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
			// 容器内以 root 运行：禁用沙箱与 leakless 守护。
			//   - 沙箱在容器内无法建立，必须禁用；
			//   - leakless guard 需与主进程 TCP 握手，容器内一旦握手失败会永久阻塞
			//     Launch（`<-ll.Pid()` 无超时），导致截图 worker 卡死、截图无响应；
			//     容器退出时进程本就会被清理，无需守护。
			l = l.NoSandbox(true).Leakless(false)
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
