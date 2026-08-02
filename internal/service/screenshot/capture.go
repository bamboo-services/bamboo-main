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

// CaptureFunc 截图核心抽象，便于测试注入替身
type CaptureFunc func(ctx context.Context, url string) ([]byte, error)

// rodCapturer 基于 rod 的无头浏览器截图器。
//
// 复用单个浏览器实例（启动一次 Chrome，每个目标新建 Page 截完即关），
// 浏览器进程由 rod 作为子进程托管（或连接外部 CDP），无需独立无头服务。
type rodCapturer struct {
	mu       sync.Mutex
	cfg      Config
	browser  *rod.Browser
	launcher *launcher.Launcher
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

	page, err := c.browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		// 页面创建失败视为浏览器实例失效，重建后下次重试
		c.reset()
		return nil, err
	}
	defer func() { _ = page.Close() }()

	// 单页操作统一受单次截图超时约束，不影响复用浏览器实例的后续使用
	page = page.Timeout(c.cfg.Timeout)

	if err := page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:             screenshotWidth,
		Height:            screenshotHeight,
		DeviceScaleFactor: 1,
		Mobile:            false,
	}); err != nil {
		return nil, err
	}
	if err := page.Navigate(url); err != nil {
		return nil, err
	}
	if err := page.WaitLoad(); err != nil {
		return nil, err
	}
	// 等待网络空闲确认「加载完毕」；超时不阻塞（长轮询/实时站点直接截图）
	_ = rod.Try(func() {
		page.WaitRequestIdle(networkIdleDuration, nil, nil, nil)()
	})

	return page.Screenshot(false, nil)
}

// ensureBrowser 获取可用的浏览器实例，优先复用；不存在时按配置构建
func (c *rodCapturer) ensureBrowser(ctx context.Context) error {
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
			l = l.Bin(c.cfg.ChromePath)
		}
		if os.Geteuid() == 0 {
			// 容器内以 root 运行时必须禁用沙箱，否则浏览器无法启动
			l = l.NoSandbox(true)
		}
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
	return nil
}

// reset 关闭并丢弃当前浏览器实例，下次 Capture 时重建
func (c *rodCapturer) reset() {
	if c.browser != nil {
		_ = c.browser.Close()
	}
	if c.launcher != nil {
		c.launcher.Cleanup()
	}
	c.browser = nil
	c.launcher = nil
}
