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

package utility

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gregex"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"xiaoMain/internal/consts"
)

// TokenLeftBearer 是一个处理带有 "Bearer " 前缀的 Token 的函数。
// 它接收一个 string 类型的参数，表示带有 "Bearer " 前缀的 Token，并返回去除 "Bearer " 前缀后的 Token。
// 如果在处理过程中出现错误，函数将返回一个错误。
//
// 参数:
// token: string 类型，表示带有 "Bearer " 前缀的 Token。
//
// 返回:
// *string 类型，表示去除 "Bearer " 前缀后的 Token。
// error 类型，如果在处理过程中出现错误，返回一个错误；否则，返回 nil。
func TokenLeftBearer(token string) (*string, error) {
	replace, err := gregex.Replace("Bearer ", []byte(""), []byte(token))
	if err != nil {
		return nil, err
	}
	getReplace := string(replace)
	return &getReplace, nil
}

// GetUUIDFromHeader 是一个从请求头中获取 X-User-Uid 字段的函数。
// 它接收一个 *ghttp.Request 类型的参数，表示 HTTP 请求。
// 如果请求头中存在 X-User-Uid 字段，函数将返回该字段的值。
// 如果请求头中不存在 X-User-Uid 字段，函数将返回一个错误。
//
// 参数:
// getRequest: *ghttp.Request 类型，表示 HTTP 请求。
//
// 返回:
// *string 类型，表示从请求头中获取的 X-User-Uid 字段的值。
// error 类型，如果无法从请求头中获取 X-User-Uid 字段，返回一个错误；否则，返回 nil。
func GetUUIDFromHeader(getRequest *ghttp.Request) (*string, error) {
	getUUID := getRequest.Header.Get("X-User-Uid")
	if getUUID != "" {
		return &getUUID, nil
	} else {
		return nil, errors.New("无法从请求头获取 UUID")
	}
}

// GetAuthorizationFromHeader 是一个从请求头中获取 Authorization 字段的函数。
// 它接收一个 *ghttp.Request 类型的参数，表示 HTTP 请求。
// 如果请求头中存在 Authorization 字段，函数将调用 TokenLeftBearer 函数去除 "Bearer " 前缀，并返回剩余的 Token 内容。
// 如果请求头中不存在 Authorization 字段，函数将返回一个错误。
//
// 参数:
// getRequest: *ghttp.Request 类型，表示 HTTP 请求。
//
// 返回:
// *string 类型，表示从请求头中获取的 Authorization 字段的值（去除 "Bearer " 前缀后的值）。
// error 类型，如果无法从请求头中获取 Authorization 字段，返回一个错误；否则，返回 nil。
func GetAuthorizationFromHeader(getRequest *ghttp.Request) (*string, error) {
	getAuthorization := getRequest.Header.Get("Authorization")
	if getAuthorization != "" {
		return TokenLeftBearer(getAuthorization)
	} else {
		return nil, errors.New("无法从请求头获取 Authorization")
	}
}

// PasswordEncode 是一个对密码进行加密的函数。
// 它接收一个 string 类型的参数，表示原始密码，并返回该密码的加密形式。
// 首先，它将原始密码转换为 Base64 编码，然后使用 bcrypt.GenerateFromPassword 函数对 Base64 编码的密码进行加密。
// 加密的成本因子为 bcrypt.DefaultCost。
//
// 参数:
// password: string 类型，表示原始密码。
//
// 返回: string 类型，表示加密后的密码。
func PasswordEncode(password string) string {
	getEncodePassword, _ := bcrypt.GenerateFromPassword(StringToBase64(password), bcrypt.DefaultCost)
	return string(getEncodePassword)
}

// PasswordVerify 是一个验证密码的函数。
// 它接收两个参数，一个是哈希密码，另一个是原始密码。
// 首先，它将原始密码转换为 Base64 编码，然后使用 bcrypt.CompareHashAndPassword 函数比较哈希密码和 Base64 编码的密码。
// 如果两者匹配，函数返回 true；否则，返回 false。
//
// 参数:
// hashPassword: string 类型，表示哈希密码。
// password: string 类型，表示原始密码。
//
// 返回: bool 类型，如果哈希密码和 Base64 编码的密码匹配，返回 true；否则，返回 false。
func PasswordVerify(hashPassword string, password string) bool {
	// 对密码进行 Base64 加密
	basePassword := StringToBase64(password)
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), basePassword)
	return err == nil
}

// StringToBase64Hex 是一个将字符串转换为 Base64 编码的函数。
// 它接收一个 string 类型的参数，并返回该字符串的 Base64 编码形式。
//
// 参数:
// stringData: string 类型，表示要转换为 Base64 编码的原始字符串。
//
// 返回: string 类型，表示转换后的 Base64 编码字符串。
func StringToBase64Hex(stringData string) string {
	return base64.StdEncoding.EncodeToString([]byte(stringData))
}

// StringToBase64 是一个将字符串转换为 Base64 编码的函数。
// 它接收一个 string 类型的参数，并返回该字符串的 Base64 编码形式。
//
// 参数:
// stringData: string 类型，表示要转换为 Base64 编码的原始字符串。
//
// 返回: []byte 类型，表示转换后的 Base64 编码字符串。
func StringToBase64(stringData string) []byte {
	return []byte(base64.StdEncoding.EncodeToString([]byte(stringData)))
}

// CheckScenesInScope 是一个检查给定场景是否在预定义范围内的函数。
// 它遍历 consts.Scenes 中的所有场景，并检查给定的场景是否在这些场景中。
// 如果给定的场景在预定义的场景中，函数将返回 true；否则，返回 false。
//
// 参数:
// scene: string 类型，表示要检查的场景名称。
//
// 返回: bool 类型，如果给定的场景在预定义的场景中，返回 true；否则，返回 false。
func CheckScenesInScope(scene consts.Scene) bool {
	for _, getScene := range consts.Scenes {
		if string(scene) == getScene {
			return true
		}
	}
	return false
}

// GetMailTemplateByScene 是一个根据给定的场景返回对应的邮件模板的函数。
// 它遍历 consts.MailTemplate 中的所有模板，并返回与给定场景名称匹配的模板的数据。
// 如果没有找到匹配的模板，它将返回一个空字符串。
//
// 参数:
// scene: consts.Scene 类型，表示邮件模板的场景名称。
//
// 返回: string 类型，表示找到的邮件模板的数据。如果没有找到匹配的模板，返回空字符串。
func GetMailTemplateByScene(scene consts.Scene) string {
	// 从 Scene 获取模板
	for _, template := range consts.MailTemplate {
		if template.Name == string(scene) {
			return template.Data
		}
	}
	return ""
}

// GetRandomString 是一个生成随机字符串的函数。
// 它接收一个 int 类型的参数，表示要生成的随机字符串的长度。
// 函数将生成一个指定长度的随机字符串，并返回该字符串。
// 随机字符串由数字、大写字母和小写字母组成。
//
// 参数:
// length: int 类型，表示要生成的随机字符串的长度。
//
// 返回: string 类型，表示生成的随机字符串。
func GetRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		randIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[randIndex.Int64()]
	}
	return string(result)
}

// GetVisitorProtocol 是一个获取访问者使用的协议类型的函数。
// 它接收一个 context.Context 类型的参数，表示上下文。
// 如果访问者使用的是 HTTPS 协议（即请求中的 TLS 字段不为 nil），函数将返回 "SSL"。
// 如果访问者使用的是 HTTP 协议（即请求中的 TLS 字段为 nil），函数将返回 "TLS"。
//
// 参数:
// ctx: context.Context 类型，表示上下文。
//
// 返回:
// string 类型，表示访问者使用的协议类型。如果访问者使用的是 HTTPS 协议，返回 "SSL"；如果访问者使用的是 HTTP 协议，返回 "TLS"。
func GetVisitorProtocol(ctx context.Context) string {
	getRequest := ghttp.RequestFromCtx(ctx)
	if getRequest.TLS != nil {
		return "SSL"
	} else {
		return "TLS"
	}
}

// GetMailSendPort 是一个获取邮件发送端口的函数。
// 它接收一个 context.Context 类型的参数，表示上下文。
// 如果访问者使用的是 HTTPS 协议（即请求中的 TLS 字段不为 nil），函数将返回 consts.SMTPPortSSL。
// 如果访问者使用的是 HTTP 协议（即请求中的 TLS 字段为 nil），函数将返回 consts.SMTPPortTLS。
//
// 参数:
// ctx: context.Context 类型，表示上下文。
//
// 返回:
// int 类型，表示邮件发送端口。如果访问者使用的是 HTTPS 协议，返回 consts.SMTPPortSSL；如果访问者使用的是 HTTP 协议，返回 consts.SMTPPortTLS。
func GetMailSendPort(ctx context.Context) int {
	if GetVisitorProtocol(ctx) == "SSL" {
		return consts.SMTPPortSSL
	} else {
		return consts.SMTPPortTLS
	}
}
