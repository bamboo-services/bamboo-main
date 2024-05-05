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
	"github.com/gogf/gf/v2/os/glog"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	"xiaoMain/internal/consts"
	"xiaoMain/internal/lutil"
)

// CheckLinkCanAccess 检查链接是否可以访问
// 用于检查用户添加的链接是否可以访问，如果可以访问则返回 nil，否则返回错误
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// siteURL: 用户尝试添加的链接地址。
//
// 返回：
// err: 如果链接可以访问，返回 nil；否则返回错误。
func (s *sLinkLogic) CheckLinkCanAccess(ctx context.Context, siteURL string) (err error) {
	glog.Notice(ctx, "[LOGIC] 执行 LinkLogic:CheckLinkCanAccess 服务层")
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
				glog.Debug(ctx, "[LOGIC] Link 链接超时")
				return errors.New("链接超时")
			} else {
				glog.Debug(ctx, "[LOGIC] Link 链接错误")
				return errors.New("链接错误")
			}
		} else {
			glog.Debug(ctx, "[LOGIC] Link 链接错误")
			return errors.New("链接错误")
		}
	}
	// 检查链接是否可以访问
	if getResp == nil {
		glog.Debug(ctx, "[LOGIC] 链接错误")
		return errors.New("网站不可达")
	}
	if getResp.StatusCode != http.StatusOK {
		glog.Debug(ctx, "[LOGIC] Link 访问状态不正确")
		return errors.New("访问状态不正确")
	}
	_ = getResp.Body.Close()
	glog.Debugf(ctx, "[LOGIC] Link 访问成功，耗时：%dms", time.Now().UnixMilli()-getNowTimestamp)
	return nil
}

// CheckLogoCanAccess 检查 Logo 是否可以访问
// 用于检查用户添加的 Logo 是否可以访问，如果可以访问则返回 nil，否则返回错误
// 并且检查获取的状态是否是图片，如果不是图片则返回错误，否则返回 nil
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// siteLogo: 用户尝试添加的 Logo 地址。
//
// 返回：
// err: 如果 Logo 可以访问，返回 nil；否则返回错误。
func (s *sLinkLogic) CheckLogoCanAccess(ctx context.Context, siteLogo string) (err error) {
	glog.Notice(ctx, "[LOGIC] 执行 LinkLogic:CheckLogoCanAccess 服务层")
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
				glog.Debug(ctx, "[LOGIC] Logo 链接超时")
				return errors.New("链接超时")
			} else {
				glog.Debug(ctx, "[LOGIC] Logo 链接错误")
				return errors.New("链接错误")
			}
		} else {
			glog.Debug(ctx, "[LOGIC] Logo 链接错误")
			return errors.New("链接错误")
		}
	}
	// 检查链接是否可以访问
	if getResp == nil {
		glog.Debug(ctx, "[LOGIC] 链接错误")
		return errors.New("网站不可达")
	}
	if getResp.StatusCode != http.StatusOK {
		glog.Debug(ctx, "[LOGIC] Logo 访问状态不正确")
		return errors.New("访问状态不正确")
	} else {
		_ = getResp.Body.Close()
		// 检查获取的状态是否是图片
		for _, imageContentType := range consts.ImageContentTypes {
			// 拆分 Content-Type 获取图片类型
			for _, getUserContentType := range strings.Split(getResp.Header.Get("Content-Type"), ";") {
				if imageContentType == getUserContentType {
					glog.Debugf(ctx, "[LOGIC] Logo 访问成功，耗时：%dms", time.Now().UnixMilli()-getNowTimestamp)
					return nil
				}
			}
		}
		glog.Noticef(ctx, "[LOGIC] Logo 不是图片，获取请求头为 [%s]", getResp.Header.Get("Content-Type"))
		return errors.New("该 URL 不是图片")
	}
}

// CheckRSSCanAccess 检查链接是否可以访问
// 用于检查用户添加的链接是否可以访问，如果可以访问则返回 nil，否则返回错误
// 还需要检查是否为 RSS URL 是否为 XML 格式
//
// 参数：
// ctx: 请求的上下文，用于管理超时和取消信号。
// siteRSS: 用户尝试添加的 RSS 地址。
//
// 返回：
// err: 如果链接可以访问，返回 nil；否则返回错误。
func (s *sLinkLogic) CheckRSSCanAccess(ctx context.Context, siteRSS string) (err error) {
	glog.Notice(ctx, "[LOGIC] 执行 LinkLogic:CheckRSSCanAccess 服务层")
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
				glog.Debug(ctx, "[LOGIC] RSS 链接超时")
				return errors.New("链接超时")
			} else {
				glog.Debug(ctx, "[LOGIC] RSS 链接错误")
				return errors.New("链接错误")
			}
		} else {
			glog.Debug(ctx, "[LOGIC] RSS 链接错误")
			return errors.New("链接错误")
		}
	}
	// 检查链接是否可以访问
	if getResp == nil {
		glog.Debug(ctx, "[LOGIC] RSS 链接错误")
		return errors.New("网站不可达")
	}
	if getResp.StatusCode != http.StatusOK {
		glog.Debug(ctx, "[LOGIC] RSS 访问状态不正确")
		return errors.New("访问状态不正确")
	} else {
		// 检查是否为 XML 格式
		getBody, err := io.ReadAll(getResp.Body)
		if err != nil {
			glog.Debug(ctx, "[LOGIC] 读取 RSS 错误")
			return errors.New("读取 RSS 错误")
		}
		err = xml.Unmarshal(getBody, new(interface{}))
		if err == nil {
			glog.Debugf(ctx, "[LOGIC] RSS 访问成功，耗时：%dms", time.Now().UnixMilli()-getNowTimestamp)
			return nil
		} else {
			glog.Errorf(ctx, "[LOGIC] RSS 不是 XML 格式")
			return errors.New("该 RSS 不是 XML 格式")
		}
	}
}
