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
