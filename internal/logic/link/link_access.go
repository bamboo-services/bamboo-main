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

package link

import (
	"context"
	"encoding/xml"
	"errors"
	"github.com/bamboo-services/bamboo-utils/berror"
	"github.com/gogf/gf/v2/frame/g"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	xerror "xiaoMain/internal/config/error"
	"xiaoMain/internal/constants"
	"xiaoMain/internal/lutil"
)

// CheckLinkCanAccess
//
// # 检查链接是否可以访问
//
// 用于检查用户添加的链接是否可以访问，如果可以访问则返回 nil，否则返回错误。
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - siteURL: 网站地址(string)
//
// # 返回:
//   - err: 如果链接可以访问，返回 nil。否则返回错误信息。
func (s *sLink) CheckLinkCanAccess(ctx context.Context, siteURL string) (err error) {
	g.Log().Notice(ctx, "[LOGIC] Link:CheckLinkCanAccess | 检查链接是否可以访问")
	getNowTimestamp := time.Now().UnixMilli()
	// 检查链接是否可以访问
	getResp, err := lutil.LinkAccess(siteURL)
	defer func() {
		if getResp != nil {
			_ = getResp.Body.Close()
		}
	}()
	// 错误处理
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		var netErr *net.OpError
		if errors.As(urlErr.Err, &netErr) {
			if netErr.Timeout() {
				g.Log().Debug(ctx, "[LOGIC] Link 链接超时")
				return berror.NewError(xerror.LinkAccessTimeout)
			} else {
				g.Log().Debug(ctx, "[LOGIC] Link 链接错误")
				return berror.NewError(xerror.LinkAccessTimeout)
			}
		} else {
			g.Log().Debug(ctx, "[LOGIC] Link 链接错误")
			return berror.NewError(xerror.LinkAccessTimeout)
		}
	}
	// 检查链接是否可以访问
	if getResp == nil {
		g.Log().Debug(ctx, "[LOGIC] 网站不可达")
		return berror.NewError(xerror.WebsiteIsUnreachable)
	}
	if getResp.StatusCode != http.StatusOK {
		g.Log().Debug(ctx, "[LOGIC] Link 访问状态不正确")
		return berror.NewError(xerror.WebIncorrectStatus)
	}
	_ = getResp.Body.Close()
	g.Log().Debugf(ctx, "[LOGIC] Link 访问成功，耗时：%dms", time.Now().UnixMilli()-getNowTimestamp)
	return nil
}

// CheckLogoCanAccess
//
// # 检查 Logo 是否可以访问
//
// 用于检查用户添加的 Logo 是否可以访问，如果可以访问则返回 nil，否则返回错误。
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - siteLogo: 网站 Logo 地址(string)
//
// # 返回:
//   - err: 如果 Logo 可以访问，返回 nil。否则返回错误信息。
func (s *sLink) CheckLogoCanAccess(ctx context.Context, siteLogo string) (err error) {
	g.Log().Notice(ctx, "[LOGIC] 执行 Link:CheckLogoCanAccess 服务层")
	getNowTimestamp := time.Now().UnixMilli()
	// 检查链接是否可以访问
	getResp, err := lutil.LinkAccess(siteLogo)
	defer func() {
		if getResp != nil {
			_ = getResp.Body.Close()
		}
	}()
	// 错误处理
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		var netErr *net.OpError
		if errors.As(urlErr.Err, &netErr) {
			if netErr.Timeout() {
				g.Log().Debug(ctx, "[LOGIC] Logo 链接超时")
				return berror.NewError(xerror.LinkAccessTimeout)
			} else {
				g.Log().Debug(ctx, "[LOGIC] Logo 链接错误")
				return berror.NewError(xerror.LinkAccessTimeout)
			}
		} else {
			g.Log().Debug(ctx, "[LOGIC] Logo 链接错误")
			return berror.NewError(xerror.LinkAccessTimeout)
		}
	}
	// 检查链接是否可以访问
	if getResp == nil {
		g.Log().Debug(ctx, "[LOGIC] 网站不可达")
		return berror.NewError(xerror.WebsiteIsUnreachable)
	}
	if getResp.StatusCode != http.StatusOK {
		g.Log().Debug(ctx, "[LOGIC] Logo 访问状态不正确")
		return berror.NewError(xerror.WebIncorrectStatus)
	} else {
		_ = getResp.Body.Close()
		// 检查获取的状态是否是图片
		for _, imageContentType := range constants.ImageContentTypes {
			// 拆分 Content-Type 获取图片类型
			for _, getUserContentType := range strings.Split(getResp.Header.Get("Content-Type"), ";") {
				if imageContentType == getUserContentType {
					g.Log().Debugf(ctx, "[LOGIC] Logo 访问成功，耗时：%dms", time.Now().UnixMilli()-getNowTimestamp)
					return nil
				}
			}
		}
		g.Log().Noticef(ctx, "[LOGIC] Logo 不是图片，获取请求头为 [%s]", getResp.Header.Get("Content-Type"))
		return errors.New("该 URL 不是图片")
	}
}

// CheckRSSCanAccess
//
// # 检查 RSS 是否可以访问
//
// 用于检查用户添加的 RSS 是否可以访问，如果可以访问则返回 nil，否则返回错误。
//
// # 参数:
//   - ctx: 上下文对象，用于传递和控制请求的生命周期。
//   - siteRSS: 网站 RSS 地址(string)
//
// # 返回:
//   - err: 如果 RSS 可以访问，返回 nil。否则返回错误信息。
func (s *sLink) CheckRSSCanAccess(ctx context.Context, siteRSS string) (err error) {
	g.Log().Notice(ctx, "[LOGIC] 执行 Link:CheckRSSCanAccess 服务层")
	getNowTimestamp := time.Now().UnixMilli()
	// 检查链接是否可以访问
	getResp, err := lutil.LinkAccess(siteRSS)
	defer func() {
		if getResp != nil {
			_ = getResp.Body.Close()
		}
	}()
	// 错误处理
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		var netErr *net.OpError
		if errors.As(urlErr.Err, &netErr) {
			if netErr.Timeout() {
				g.Log().Debug(ctx, "[LOGIC] RSS 链接超时")
				return berror.NewError(xerror.LinkAccessTimeout)
			} else {
				g.Log().Debug(ctx, "[LOGIC] RSS 链接错误")
				return berror.NewError(xerror.LinkAccessTimeout)
			}
		} else {
			g.Log().Debug(ctx, "[LOGIC] RSS 链接错误")
			return berror.NewError(xerror.LinkAccessTimeout)
		}
	}
	// 检查链接是否可以访问
	if getResp == nil {
		g.Log().Debug(ctx, "[LOGIC] RSS 网站不可达")
		return berror.NewError(xerror.WebsiteIsUnreachable)
	}
	if getResp.StatusCode != http.StatusOK {
		g.Log().Debug(ctx, "[LOGIC] RSS 访问状态不正确")
		return berror.NewError(xerror.WebIncorrectStatus)
	} else {
		// 检查是否为 XML 格式
		getBody, err := io.ReadAll(getResp.Body)
		if err != nil {
			g.Log().Debug(ctx, "[LOGIC] 读取 RSS 错误")
			return berror.NewError(xerror.ReadRSSFailed)
		}
		err = xml.Unmarshal(getBody, new(interface{}))
		if err == nil {
			g.Log().Debugf(ctx, "[LOGIC] RSS 访问成功，耗时：%dms", time.Now().UnixMilli()-getNowTimestamp)
			return nil
		} else {
			g.Log().Errorf(ctx, "[LOGIC] RSS 不是 XML 格式")
			return berror.NewError(xerror.RSSIsNotXML)
		}
	}
}
