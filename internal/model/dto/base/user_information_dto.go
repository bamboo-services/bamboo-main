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

package base

// UserSimpleDTO
//
// 用户简单信息数据传输对象；
// 包含用户名、电子邮箱和手机号码等基本信息。
type UserSimpleDTO struct {
	Username string `json:"username" dc:"用户名"`
	Email    string `json:"email" dc:"电子邮箱"`
	Phone    string `json:"phone" dc:"手机号码"`
}

// UserDetailDTO
//
// 用户详细信息数据传输对象；
// 包含用户简单信息、昵称、头像类型、头像Base64编码和头像URL等详细信息。
type UserDetailDTO struct {
	UserSimpleDTO
	Nickname     string `json:"nickname" dc:"昵称"`
	AvatarType   string `json:"avatar_type" dc:"头像类型(local,url)"`
	AvatarBase64 string `json:"avatar_base64" dc:"头像Base64"`
	AvatarURL    string `json:"avatar_url" dc:"头像URL"`
}
