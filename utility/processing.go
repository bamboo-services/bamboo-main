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
	"encoding/base64"
	"errors"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gregex"
	"golang.org/x/crypto/bcrypt"
	"xiaoMain/internal/consts"
)

// TokenLeftBearer
// 获取一个带有 Bearer 的 Token，通过此工具后，将会把 Bearer 的部分去掉，返回剩余的 Token 内容
func TokenLeftBearer(token string) (*string, error) {
	replace, err := gregex.Replace("Bearer ", []byte(""), []byte(token))
	if err != nil {
		return nil, err
	}
	getReplace := string(replace)
	return &getReplace, nil
}

// GetUUIDFromHeader
// 获取请求头中的 X-User-Uid 字段，并返回其值；若能够正确获取值将会返回值，若无法返回，将会返回 err 错误内容
func GetUUIDFromHeader(getRequest *ghttp.Request) (*string, error) {
	getUUID := getRequest.Header.Get("X-User-Uid")
	if getUUID != "" {
		return &getUUID, nil
	} else {
		return nil, errors.New("无法从请求头获取 UUID")
	}
}

// GetAuthorizationFromHeader
// 获取请求头中的 Authorization 字段，并返回其值；若能够正确获取值将会返回值，若无法返回，将会返回 err 错误内容
func GetAuthorizationFromHeader(getRequest *ghttp.Request) (*string, error) {
	getAuthorization := getRequest.Header.Get("Authorization")
	if getAuthorization != "" {
		return TokenLeftBearer(getAuthorization)
	} else {
		return nil, errors.New("无法从请求头获取 Authorization")
	}
}

// PasswordEncode
// 对输入的原始密码进行加密，若密码没有出现问题将会进行解密，若出现无法加密的情况将会返回 error 内容
func PasswordEncode(password string) string {
	getEncodePassword, _ := bcrypt.GenerateFromPassword(StringToBase64(password), bcrypt.DefaultCost)
	return string(getEncodePassword)
}

// PasswordVerify
// 对输入的密码进行检查，检查需要对原密码进行 Base64 转为 Hex 后对 Hash 进行比较；
// 期中 Hash 密码与 Base64 加密后的密码之间关系为 Bcrypt
func PasswordVerify(hashPassword string, password string) bool {
	// 对密码进行 Base64 加密
	basePassword := StringToBase64(password)
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), basePassword)
	return err == nil
}

// StringToBase64Hex
// 将字符串转为 Base64 的 Hex 形式
func StringToBase64Hex(stringData string) string {
	return base64.StdEncoding.EncodeToString([]byte(stringData))
}

// StringToBase64
// 将字符串转为 Base64 形式
func StringToBase64(stringData string) []byte {
	return []byte(base64.StdEncoding.EncodeToString([]byte(stringData)))
}

// CheckScenesInScope
// 检查场景是否在范围内，若在范围内将会返回 true，若不在范围内将会返回 false
func CheckScenesInScope(scene string) bool {
	for _, getScene := range consts.Scenes {
		if scene == getScene {
			return true
		}
	}
	return false
}
