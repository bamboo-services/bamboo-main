// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package screenshot

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// 截图视口尺寸（16:9 横屏）
const (
	screenshotWidth  = 1280
	screenshotHeight = 720
)

// networkIdleDuration 连续无新请求即视为「加载完毕」的判定时长
const networkIdleDuration = 500 * time.Millisecond

// 浏览器实例生命周期上限。
//
// Chromium 长期复用会积累内存与渲染子进程，尤其容器内无 GPU、走软件渲染路径时，
// 长时间运行后可能僵死（CDP 无响应）。到期主动关闭重建，避免「运行一段时间后
// 截图失败/卡死」。
const (
	maxBrowserAge  = 30 * time.Minute // 单实例最长存活时长
	maxBrowserUses = 50               // 单实例最多复用次数
)

// CaptureFunc 截图核心抽象，便于测试注入替身
type CaptureFunc func(ctx context.Context, url string) ([]byte, error)

// rodCapturer 基于 rod 的无头浏览器截图器。
//
// 复用单个浏览器实例（启动一次 Chrome，每个目标新建 Page 截完即关），
// 浏览器进程由 rod 作为子进程托管（或连接外部 CDP），无需独立无头服务。
// 浏览器实例有生命周期上限（maxBrowserAge / maxBrowserUses），到期主动重建；
// 任一截图失败路径都会丢弃当前实例，保证下次任务基于全新浏览器。
type rodCapturer struct {
	mu       sync.Mutex
	cfg      Config
	browser  *rod.Browser
	launcher *launcher.Launcher

	createdAt time.Time // 浏览器实例创建时间，用于到期重建
	useCount  int       // 已复用次数，用于到期重建
}

// NewRodCapture 创建基于 rod 的截图函数（默认实现）
func NewRodCapture(cfg Config) CaptureFunc {
	c := &rodCapturer{cfg: cfg}
	return c.Capture
}

// Capture 截取指定 URL 当前视口（1280×720）的 PNG 截图
func (c *rodCapturer) Capture(ctx context.Context, url string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.ensureBrowser(ctx); err != nil {
		return nil, err
	}

	// browser.Page 内部经 browser.ctx 发 CDP 调用，该 ctx 为 worker 长生命周期 ctx，
	// 无 deadline。一旦 Chromium 僵死（内存耗尽/渲染线程挂死），此调用会永久阻塞，
	// 卡死整个截图 worker。故以单次截图超时为界包装 Page 创建，超时即判定实例失效。
	page, err := c.newPage(ctx)
	if err != nil {
		c.reset()
		return nil, err
	}
	// 关闭页面走带超时的克隆 ctx：浏览器僵死时关闭页面也需兜底，避免阻塞 worker
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_ = page.Context(closeCtx).Close()
		cancel()
	}()

	// 单页操作统一受单次截图超时约束，不影响复用浏览器实例的后续使用
	page = page.Timeout(c.cfg.Timeout)

	if err := page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:             screenshotWidth,
		Height:            screenshotHeight,
		DeviceScaleFactor: 1,
		Mobile:            false,
	}); err != nil {
		c.reset()
		return nil, err
	}
	if err := page.Navigate(url); err != nil {
		c.reset()
		return nil, err
	}
	if err := page.WaitLoad(); err != nil {
		c.reset()
		return nil, err
	}
	// 等待网络空闲确认「加载完毕」；超时不阻塞（长轮询/实时站点直接截图）
	_ = rod.Try(func() {
		page.WaitRequestIdle(networkIdleDuration, nil, nil, nil)()
	})

	data, err := page.Screenshot(false, nil)
	if err != nil {
		c.reset()
		return nil, err
	}
	c.useCount++
	return data, nil
}

// newPage 以单次截图超时为界创建新页面，避免僵死浏览器阻塞 worker。
//
// browser.Page 内部经 browser.ctx（长生命周期、无 deadline）发起 CDP 调用，
// 故通过临时克隆一个带超时 ctx 的 Browser 实例来创建页面：超时后调用返回
// context deadline exceeded，上层据此丢弃实例。
func (c *rodCapturer) newPage(ctx context.Context) (*rod.Page, error) {
	pageCtx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()
	return c.browser.Context(pageCtx).Page(proto.TargetCreateTarget{})
}

// ensureBrowser 获取可用的浏览器实例，优先复用；不存在或超期/超次时按配置重建
func (c *rodCapturer) ensureBrowser(ctx context.Context) error {
	// 浏览器实例达到生命周期上限：主动关闭重建，防止长期运行后的内存膨胀/僵死
	if c.browser != nil && (c.useCount >= maxBrowserUses || time.Since(c.createdAt) > maxBrowserAge) {
		c.reset()
	}
	if c.browser != nil {
		return nil
	}

	var (
		browser *rod.Browser
		cdpURL  string
		err     error
	)
	if c.cfg.CDPURL != "" {
		// 外部 CDP 模式：连接独立浏览器服务（如 chromium 容器），无需本机浏览器
		cdpURL = c.cfg.CDPURL
	} else {
		// 内置模式：spawn 本机 Chrome/Chromium 子进程，由 rod 托管生命周期
		l := launcher.New().Headless(true)
		if c.cfg.ChromePath != "" {
			// 显式指定路径：先校验存在性，路径无效直接报错，避免 rod 的晦涩错误
			if _, err := os.Stat(c.cfg.ChromePath); err != nil {
				return fmt.Errorf("SCREENSHOT_CHROME_PATH 指向的浏览器不存在：%s", c.cfg.ChromePath)
			}
			l = l.Bin(c.cfg.ChromePath)
		} else if bin, has := launcher.LookPath(); has {
			// 未配置路径时探测系统浏览器，避免触发 rod 自动下载
			l = l.Bin(bin)
		} else {
			// 未配置路径且探测不到浏览器：直接报错，绝不触发 rod 自动下载。
			// 自动下载的浏览器为 glibc 构建，与 alpine/musl 运行时镜像不兼容，
			// 会导致截图失败且下载进度日志刷屏。
			return fmt.Errorf("未找到 Chrome/Chromium 可执行文件，请设置 SCREENSHOT_CHROME_PATH")
		}
		if os.Geteuid() == 0 {
			// 容器内以 root 运行：禁用沙箱、leakless 守护与 GPU 加速。
			//   - 沙箱在容器内无法建立，必须禁用；
			//   - leakless guard 需与主进程 TCP 握手，容器内一旦握手失败会永久阻塞
			//     Launch（`<-ll.Pid()` 无超时），导致截图 worker 卡死、截图无响应；
			//   - 容器无 GPU，chromium 默认自动启用 swiftshader 软件渲染会导致渲染进程
			//     卡死（CPU 100%、截图永不返回），显式 --disable-gpu 走纯 CPU 渲染路径；
			//     容器退出时进程本就会被清理，无需守护。
			l = l.NoSandbox(true).Leakless(false).Set("disable-gpu")
		}
		// 浏览器启动受单次截图超时约束，避免容器内启动异常（挂起/缓慢）时 Launch 永久阻塞
		launchCtx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
		defer cancel()
		l = l.Context(launchCtx)
		if cdpURL, err = l.Launch(); err != nil {
			return err
		}
		c.launcher = l
	}

	browser = rod.New().ControlURL(cdpURL).Context(ctx)
	if err = browser.Connect(); err != nil {
		c.reset()
		return err
	}
	c.browser = browser
	c.createdAt = time.Now()
	c.useCount = 0
	return nil
}

// reset 关闭并丢弃当前浏览器实例，下次 Capture 时重建。
//
// 顺序关键：先 SIGKILL 进程组（不依赖 CDP 响应），再关浏览器连接（带超时克隆 ctx），
// 最后 Cleanup 等进程退出并清理用户数据目录。全程不依赖僵死浏览器的响应，
// 避免「浏览器已僵死 → 关闭也阻塞 → worker 永久卡死」的连锁故障。
func (c *rodCapturer) reset() {
	if c.launcher != nil {
		// 强制杀掉浏览器进程组：CDP 无响应时也能立即终止僵死进程
		c.launcher.Kill()
	}
	if c.browser != nil {
		// 关闭浏览器连接；浏览器已死则快速返回，僵死则由超时兜底不阻塞
		closeCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_ = c.browser.Context(closeCtx).Close()
		cancel()
	}
	if c.launcher != nil {
		// 等进程退出并清理 user-data-dir；进程已被 Kill，此步不会阻塞
		c.launcher.Cleanup()
	}
	c.browser = nil
	c.launcher = nil
	c.createdAt = time.Time{}
	c.useCount = 0
}
