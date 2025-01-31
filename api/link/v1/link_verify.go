/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 *
 * 本文件包含 XiaoMain 的源代码，该项目的所有源代码均遵循MIT开源许可证协议。
 * --------------------------------------------------------------------------------
 * 许可证声明：
 *
 * 版权所有 (c) 2016-2024 筱锋。保留所有权利。
 *
 * 本软件是“按原样”提供的，没有任何形式的明示或暗示的保证，包括但不限于
 * 对适销性、特定用途的适用性和非侵权性的暗示保证。在任何情况下，
 * 作者或版权持有人均不承担因软件或软件的使用或其他交易而产生的、
 * 由此引起的或以任何方式与此软件有关的任何索赔、损害或其他责任。
 *
 * 使用本软件即表示您了解此声明并同意其条款。
 *
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 * 免责声明：
 *
 * 使用本软件的风险由用户自担。作者或版权持有人在法律允许的最大范围内，
 * 对因使用本软件内容而导致的任何直接或间接的损失不承担任何责任。
 * --------------------------------------------------------------------------------
 */

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// LinkVerifyReq
//
// # 审核站点
//
// 审核站点是否通过
//
// # 参数
//   - ID: 链接ID
//   - Status: 审核状态
type LinkVerifyReq struct {
	g.Meta          `path:"/link/verify" method:"GET" tags:"链接控制器" summary:"审核站点" dc:"审核站点是否通过"`
	Id              int64 `json:"id" v:"required#请输入链接ID"`
	DesiredLocation int64 `json:"desired_location" v:"required#请输入期望位置"`
	DesiredColor    int64 `json:"desired_color" v:"required#请输入期望颜色"`
	Status          bool  `json:"status" v:"required#请输入审核状态"`
}

// LinkVerifyRes
//
// # 审核站点响应
//
// 审核站点是否通过响应
type LinkVerifyRes struct {
	g.Meta `mime:"application/json"`
}
