// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package constants

// 内置「已失效」分组在 bm_system 表中的配置键。
//
// 键值由管理员经管理端接口热修改，分组列表/友链查询在注入虚拟分组时读取，
// 缺失或空值回退内置默认名「已失效」。
const (
	// KeyBuiltinInvalidGroupName 内置「已失效」分组名称配置键。
	KeyBuiltinInvalidGroupName = "group.builtin.invalid.name"
	// KeyBuiltinInvalidGroupDesc 内置「已失效」分组描述配置键。
	KeyBuiltinInvalidGroupDesc = "group.builtin.invalid.description"
)

// 站点级高级配色开关（color.mode）在 bm_system 表中的配置键。
//
// 决定颜色选择器可见范围与渲染级别：
//   - ColorModeNormal：仅展示普通配色与炫彩（高级色隐藏，不可选）
//   - ColorModePremium：额外展示高级配色（主/副/悬停三色渐变渲染）
const (
	// KeyColorMode 高级配色模式配置键。
	KeyColorMode = "color.mode"
	// ColorModeNormal 普通模式：仅普通色 + 炫彩。
	ColorModeNormal = "normal"
	// ColorModePremium 高级模式：普通色 + 高级色 + 炫彩。
	ColorModePremium = "premium"
)
