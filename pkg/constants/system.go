// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package constants

// KeySystemAdminID 系统唯一管理员的用户 ID 配置键。
//
// 初始化默认管理员后写入 bm_system 表，鉴权时以「当前用户 ID == 该值」判定管理身份，
// 取代 bm_system_user 上的角色字段。缺失或空值视为无管理员。
const KeySystemAdminID = "system.admin.id"

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
