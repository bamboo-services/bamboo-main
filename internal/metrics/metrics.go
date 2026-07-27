// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

// Package metrics 提供 Prometheus 请求埋点中间件、/metrics 暴露处理器，
// 以及跨平台的运行时指标采集（协程数 / CPU / 内存）。
package metrics

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/v4/process"
)

// namespace 统一指标命名空间前缀
const namespace = "bamboo"

// startTime 进程启动时间，用于计算运行时长
var startTime = time.Now()

// proc 当前进程句柄，用于采集 CPU 使用率（gopsutil 跨平台）
var proc *process.Process

var (
	// registry 自定义注册表，汇聚 HTTP 指标、运行时指标与标准 Go/进程采集器
	registry = prometheus.NewRegistry()

	// httpRequestsTotal HTTP 请求计数（按方法 / 路由模板 / 状态码）
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_requests_total",
			Help:      "HTTP 请求总数（按方法、路由模板、状态码）",
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestDuration HTTP 请求耗时分布（按方法 / 路由模板）
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP 请求耗时分布（秒，按方法、路由模板）",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	// 初始化进程句柄并预热一次 CPU 采样，建立基线，
	// 使后续 Percent(0) 返回「距上次采样」的区间均值而非 0。
	if p, err := process.NewProcess(int32(os.Getpid())); err == nil {
		proc = p
		_, _ = proc.Percent(0)
	}

	registry.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		newRuntimeCollector(),
		collectors.NewGoCollector(), // go_goroutines / go_memstats_* / go_gc_* 等
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}), // process_cpu_seconds_total / process_resident_memory_bytes 等（Linux 完整）
	)
}

// Middleware Prometheus 请求埋点中间件：在请求完成后记录计数与耗时。
//
// path 使用路由模板（c.FullPath()，如 /api/v1/admin/links/:id）而非真实 URL，
// 避免 ID 等动态段造成指标基数爆炸；未匹配路由记为 unmatched。
// 跳过 /metrics 自身，避免抓取请求污染统计。
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = "unmatched"
		}
		status := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(time.Since(start).Seconds())
	}
}

// Handler 返回 /metrics 端点处理器。promhttp 直写 Prometheus 文本格式，
// 框架 ResponseMiddleware 检测到响应已写入会自动放行，不会包装为 BaseResponse。
func Handler() gin.HandlerFunc {
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// ---------------------------------------------------------------------------
// 运行时快照（供 health 接口与 runtimeCollector 复用）
// ---------------------------------------------------------------------------

// RuntimeStats 运行时快照
type RuntimeStats struct {
	Uptime      time.Duration // 运行时长
	Goroutines  int           // 当前 Goroutine 数量
	MemoryBytes uint64        // 已分配内存（字节，runtime.ReadMemStats.Alloc）
	CPUPercent  float64       // 进程 CPU 使用率（百分比，gopsutil）
}

// Snapshot 采集当前运行时快照：运行时长、协程数、内存（runtime）、CPU（gopsutil）。
func Snapshot() RuntimeStats {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	cpuPercent := 0.0
	if proc != nil {
		if pct, err := proc.Percent(0); err == nil {
			cpuPercent = pct
		}
	}

	return RuntimeStats{
		Uptime:      time.Since(startTime),
		Goroutines:  runtime.NumGoroutine(),
		MemoryBytes: memStats.Alloc,
		CPUPercent:  cpuPercent,
	}
}

// runtimeCollector 自定义运行时采集器：跨平台暴露协程数 / CPU / 内存三项指标。
// Collect 时实时读取 Snapshot，无需后台协程；保证 macOS 开发机与 Linux 生产机指标齐全。
type runtimeCollector struct {
	goroutines  *prometheus.Desc
	cpuPercent  *prometheus.Desc
	memoryBytes *prometheus.Desc
}

func newRuntimeCollector() *runtimeCollector {
	return &runtimeCollector{
		goroutines: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "runtime", "goroutines"),
			"当前 Goroutine 数量", nil, nil,
		),
		cpuPercent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "runtime", "cpu_percent"),
			"进程 CPU 使用率（百分比）", nil, nil,
		),
		memoryBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "runtime", "memory_bytes"),
			"进程已分配内存（字节）", nil, nil,
		),
	}
}

// Describe 实现 prometheus.Collector 接口
func (c *runtimeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.goroutines
	ch <- c.cpuPercent
	ch <- c.memoryBytes
}

// Collect 实现 prometheus.Collector 接口，抓取时实时采样
func (c *runtimeCollector) Collect(ch chan<- prometheus.Metric) {
	snap := Snapshot()
	ch <- prometheus.MustNewConstMetric(c.goroutines, prometheus.GaugeValue, float64(snap.Goroutines))
	ch <- prometheus.MustNewConstMetric(c.cpuPercent, prometheus.GaugeValue, snap.CPUPercent)
	ch <- prometheus.MustNewConstMetric(c.memoryBytes, prometheus.GaugeValue, float64(snap.MemoryBytes))
}
