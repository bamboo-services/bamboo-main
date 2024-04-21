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

package result

// ================================================================================================================
// 常见错误代码定义
// 框架保留了内部错误代码：代码 < xxx30
// ================================================================================================================

var (
	NotLoggedIn                = ErrorBaseResponse{Output: "NotLoggedIn", Code: 40101, Message: "用户未登录"}
	RequestBodyValidationError = ErrorBaseResponse{Output: "RequestBodyValidationError", Code: 40301, Message: "请求正文验证错误"}
	QueryValidationError       = ErrorBaseResponse{Output: "QueryValidationError", Code: 40302, Message: "查询验证错误"}
	PathValidationError        = ErrorBaseResponse{Output: "PathValidationError", Code: 40303, Message: "路径参数错误"}
	VerificationFailed         = ErrorBaseResponse{Output: "VerificationFailed", Code: 40304, Message: "校验失败"}
	MailError                  = ErrorBaseResponse{Output: "MailError", Code: 40305, Message: "邮件错误"}
	AccessError                = ErrorBaseResponse{Output: "AccessError", Code: 40306, Message: "访问错误"}
	AddLinkFailed              = ErrorBaseResponse{Output: "AddLinkFailed", Code: 40307, Message: "添加链接失败"}

	ServerInternalError = ErrorBaseResponse{Output: "ServerInternalError", Code: 50000, Message: "服务器内部错误", Data: errorData{ErrorMessage: "服务器内部错误"}} //nolint:lll
	DatabaseError       = ErrorBaseResponse{Output: "DatabaseError", Code: 50001, Message: "服务器内部错误", Data: errorData{ErrorMessage: "数据库错误"}}         //nolint:lll
)
