// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package constants

import "github.com/bamboo-services/bamboo-base-go/foundation/v2/apperr"

// 业务专属错误码（遵循 foundation/v2/apperr 分组内 50..99 号段约定）。
var (
	// ErrCodeLoginFailed 登录凭证校验失败（HTTP 401）
	ErrCodeLoginFailed = &apperr.ErrorCode{
		Code:    40150,
		Output:  "BAMBOO_LOGIN_FAILED",
		Message: "用户名或密码错误",
	}
)
