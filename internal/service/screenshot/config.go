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
	"time"

	xEnv "github.com/bamboo-services/bamboo-base-go/defined/env"
)

// 截图服务环境变量键名（env-first 配置约定）
const (
	// EnvScreenshotEnabled 截图服务总开关
	EnvScreenshotEnabled xEnv.EnvKey = "SCREENSHOT_ENABLED"
	// EnvScreenshotDir 截图存储目录（相对进程工作目录）
	EnvScreenshotDir xEnv.EnvKey = "SCREENSHOT_DIR"
	// EnvScreenshotChromePath Chrome/Chromium 可执行文件路径（留空自动探测）
	EnvScreenshotChromePath xEnv.EnvKey = "SCREENSHOT_CHROME_PATH"
	// EnvScreenshotCDPURL 外部 CDP 服务地址（留空则内置 spawn 浏览器子进程）
	EnvScreenshotCDPURL xEnv.EnvKey = "SCREENSHOT_CDP_URL"
	// EnvScreenshotTimeout 单次截图超时秒数
	EnvScreenshotTimeout xEnv.EnvKey = "SCREENSHOT_TIMEOUT"
)

// Config 友链站点截图服务配置
type Config struct {
	Enabled    bool          // 服务总开关
	Dir        string        // 截图存储目录
	ChromePath string        // Chrome/Chromium 可执行文件路径
	CDPURL     string        // 外部 CDP 服务地址
	Timeout    time.Duration // 单次截图超时
}

// LoadConfig 从环境变量加载截图服务配置
func LoadConfig() Config {
	return Config{
		Enabled:    xEnv.GetEnvBool(EnvScreenshotEnabled, true),
		Dir:        xEnv.GetEnvString(EnvScreenshotDir, "data/screenshots"),
		ChromePath: xEnv.GetEnvString(EnvScreenshotChromePath, ""),
		CDPURL:     xEnv.GetEnvString(EnvScreenshotCDPURL, ""),
		Timeout:    time.Duration(xEnv.GetEnvInt64(EnvScreenshotTimeout, 30)) * time.Second,
	}
}
