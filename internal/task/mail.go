/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package task

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/smtp"
	"strings"
	"sync"
	"time"

	"bamboo-main/internal/model/base"
	dtoRedis "bamboo-main/internal/model/dto/redis"
	"bamboo-main/pkg/constants"

	"github.com/jordan-wright/email"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// MailWorker 邮件工作协程
type MailWorker struct {
	rdb         *redis.Client      // Redis 客户端
	config      *base.EmailConfig  // 邮件配置
	pool        *TLSMailPool       // TLS 邮件连接池
	logger      *zap.SugaredLogger // 日志记录器
	ctx         context.Context    // 上下文
	cancel      context.CancelFunc // 取消函数
	wg          sync.WaitGroup     // 等待组
	workerCount int                // 工作协程数
	timeout     time.Duration      // 发送超时
}

// NewMailWorker 创建邮件工作协程实例
//
// 参数说明:
//   - rdb: Redis 客户端
//   - config: 邮件配置
//   - logger: 日志记录器
//
// 返回值:
//   - 邮件工作协程实例
func NewMailWorker(rdb *redis.Client, config *base.EmailConfig, logger *zap.SugaredLogger) *MailWorker {
	// 设置默认工作协程数
	workerCount := config.WorkerCount
	if workerCount <= 0 {
		workerCount = 4
	}

	// 设置默认超时时间
	timeout := time.Duration(config.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &MailWorker{
		rdb:         rdb,
		config:      config,
		logger:      logger,
		workerCount: workerCount,
		timeout:     timeout,
	}
}

// Start 启动守护协程
func (w *MailWorker) Start() {
	// 检查邮件配置是否有效
	if w.config.SMTPHost == "" || w.config.Username == "" {
		w.logger.Warn("邮件配置不完整，邮件服务未启动")
		return
	}

	// 创建 SMTP 连接池
	addr := fmt.Sprintf("%s:%d", w.config.SMTPHost, w.config.SMTPPort)
	auth := smtp.PlainAuth("", w.config.Username, w.config.Password, w.config.SMTPHost)

	// TLS 配置
	tlsConfig := &tls.Config{
		ServerName: w.config.SMTPHost,
	}

	// 智能判断 TLS 模式
	useTLS := w.config.UseTLS
	if !useTLS && !w.config.UseStartTLS {
		// 如果两者都未配置，根据端口自动判断
		if w.config.SMTPPort == 465 {
			useTLS = true
			w.logger.Info("检测到 465 端口，自动启用 TLS 直连模式")
		} else if w.config.SMTPPort == 587 {
			useTLS = false
			w.logger.Info("检测到 587 端口，自动启用 STARTTLS 模式")
		}
	}

	// 创建自定义 TLS 连接池
	w.pool = NewTLSMailPool(addr, w.workerCount, auth, tlsConfig, useTLS)

	w.ctx, w.cancel = context.WithCancel(context.Background())

	// 启动多个工作协程
	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.processQueue(i)
	}

	// 启动重试调度器
	w.wg.Add(1)
	go w.retryScheduler()

	w.logger.Infof("邮件工作协程已启动，工作协程数: %d, TLS模式: %v", w.workerCount, useTLS)
}

// Stop 优雅停止守护协程
func (w *MailWorker) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
	w.wg.Wait()

	// 关闭连接池
	if w.pool != nil {
		w.pool.Close()
	}

	w.logger.Info("邮件工作协程已停止")
}

// processQueue 处理队列（单个工作协程）
func (w *MailWorker) processQueue(workerID int) {
	defer w.wg.Done()

	w.logger.Infof("工作协程 #%d 已启动", workerID)

	for {
		select {
		case <-w.ctx.Done():
			w.logger.Infof("工作协程 #%d 正在退出", workerID)
			return
		default:
			// 阻塞获取任务（最多等待 5 秒）
			result, err := w.rdb.BRPop(w.ctx, 5*time.Second, constants.MailQueueKey).Result()
			if err != nil {
				if err == redis.Nil || strings.Contains(err.Error(), "context canceled") {
					continue
				}
				w.logger.Errorf("工作协程 #%d 获取任务失败: %v", workerID, err)
				continue
			}

			// 解析任务
			var task dtoRedis.MailTaskDTO
			if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
				w.logger.Errorf("解析邮件任务失败: %v", err)
				continue
			}

			w.logger.Infof("工作协程 #%d 处理任务: ID=%s, To=%v", workerID, task.ID, task.To)

			// 发送邮件
			if err := w.sendEmail(&task); err != nil {
				w.handleError(&task, err)
			} else {
				w.updateStats("success")
				w.logger.Infof("邮件发送成功: ID=%s, To=%v", task.ID, task.To)
			}
		}
	}
}

// sendEmail 使用 email 库发送邮件
func (w *MailWorker) sendEmail(task *dtoRedis.MailTaskDTO) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", w.config.FromName, w.config.FromEmail)
	e.To = task.To
	e.Cc = task.Cc
	e.Subject = task.Subject
	e.HTML = []byte(task.Body)

	// 使用连接池发送
	return w.pool.Send(e, w.timeout)
}

// handleError 处理发送失败
func (w *MailWorker) handleError(task *dtoRedis.MailTaskDTO, sendErr error) {
	w.logger.Errorf("邮件发送失败: ID=%s, To=%v, Error=%v, RetryCount=%d",
		task.ID, task.To, sendErr, task.RetryCount)

	// 检查是否还能重试
	if task.RetryCount < task.MaxRetry {
		if err := w.requeueWithBackoff(task); err != nil {
			w.logger.Errorf("邮件重试入队失败: ID=%s, Error=%v", task.ID, err)
			w.saveFailedTask(task, sendErr)
		}
	} else {
		w.logger.Errorf("邮件达到最大重试次数，放弃发送: ID=%s", task.ID)
		w.saveFailedTask(task, sendErr)
		w.updateStats("failed")
	}
}

// calculateBackoff 计算指数退避时间
//
// 公式: baseDelay * 2^retryCount + jitter
// retryCount=0: 1s  + jitter (0-500ms)
// retryCount=1: 2s  + jitter
// retryCount=2: 4s  + jitter
// retryCount=3: 8s  + jitter
func (w *MailWorker) calculateBackoff(retryCount int) time.Duration {
	baseDelay := 1 * time.Second
	maxDelay := 30 * time.Second

	delay := time.Duration(math.Pow(2, float64(retryCount))) * baseDelay
	if delay > maxDelay {
		delay = maxDelay
	}

	// 添加随机抖动（0-500ms）避免惊群效应
	jitter := time.Duration(rand.Intn(500)) * time.Millisecond
	return delay + jitter
}

// requeueWithBackoff 重试入队（指数退避）
func (w *MailWorker) requeueWithBackoff(task *dtoRedis.MailTaskDTO) error {
	task.RetryCount++
	backoff := w.calculateBackoff(task.RetryCount)
	task.NextRetryAt = time.Now().Add(backoff)

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	// 使用 Sorted Set 存储，score 为下次执行时间戳
	score := float64(task.NextRetryAt.Unix())
	return w.rdb.ZAdd(w.ctx, constants.MailRetryQueueKey, redis.Z{
		Score:  score,
		Member: taskJSON,
	}).Err()
}

// retryScheduler 重试调度器
//
// 每 5 秒检查重试队列，将到期的任务移回主队列
func (w *MailWorker) retryScheduler() {
	defer w.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	w.logger.Info("重试调度器已启动")

	for {
		select {
		case <-w.ctx.Done():
			w.logger.Info("重试调度器正在退出")
			return
		case <-ticker.C:
			w.processRetryQueue()
		}
	}
}

// processRetryQueue 处理重试队列
func (w *MailWorker) processRetryQueue() {
	now := float64(time.Now().Unix())

	// 获取已到期的重试任务
	tasks, err := w.rdb.ZRangeByScore(w.ctx, constants.MailRetryQueueKey, &redis.ZRangeBy{
		Min: "0",
		Max: fmt.Sprintf("%f", now),
	}).Result()

	if err != nil {
		w.logger.Errorf("获取重试任务失败: %v", err)
		return
	}

	for _, taskJSON := range tasks {
		// 使用事务确保原子性
		pipe := w.rdb.Pipeline()
		pipe.LPush(w.ctx, constants.MailQueueKey, taskJSON)
		pipe.ZRem(w.ctx, constants.MailRetryQueueKey, taskJSON)

		if _, err := pipe.Exec(w.ctx); err != nil {
			w.logger.Errorf("移动重试任务失败: %v", err)
		} else {
			w.logger.Info("重试任务已移回主队列")
		}
	}
}

// saveFailedTask 保存失败的任务
func (w *MailWorker) saveFailedTask(task *dtoRedis.MailTaskDTO, err error) {
	failedKey := fmt.Sprintf(constants.MailFailedKey, task.ID)

	failedData := map[string]interface{}{
		"task":      task,
		"error":     err.Error(),
		"failed_at": time.Now(),
	}

	dataJSON, _ := json.Marshal(failedData)
	// 保留 7 天
	w.rdb.Set(w.ctx, failedKey, dataJSON, 7*24*time.Hour)
}

// updateStats 更新统计信息
func (w *MailWorker) updateStats(status string) {
	w.rdb.HIncrBy(w.ctx, constants.MailStatsKey, status, 1)
	w.rdb.HSet(w.ctx, constants.MailStatsKey, "last_updated", time.Now().Format(time.RFC3339))
}
