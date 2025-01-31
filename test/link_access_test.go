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

package test

import (
	"context"
	"testing"
	"xiaoMain/internal/logic/link"
	"xiaoMain/internal/service"
)

// TestBlogUrlAbleToUse 测试博客地址是否可以访问
// 用于测试博客地址是否可以访问，如果可以访问则返回 nil，否则返回错误
//
// 参数：
// t: 测试实例
func TestBlogURLAbleToUse(t *testing.T) {
	service.RegisterLink(link.New())
	// 测试其中一个地址是否可以访问
	// 正确可以返回的
	err := service.Link().CheckLinkCanAccess(context.Background(), "https://blog.x-lf.com/")
	if err != nil {
		t.Errorf("测试失败，错误原因：%s", err.Error())
	}

	// 错误的不可返回的
	err = service.Link().CheckLinkCanAccess(context.Background(), "https://blog.x-lf.com/404")
	if err != nil {
		t.Logf("测试失败，错误原因：%s", err.Error())
	}
}

// TestLogoUrlAbleToUse 测试 Logo 地址是否可以访问
// 用于测试 Logo 地址是否可以访问，如果可以访问则返回 nil，否则返回错误
// 进行测试时，需要注意 Logo 地址是否为图片，否则会返回错误
// 该函数只允许访问图片地址，如果不是图片地址则会返回错误
//
// 参数：
// t: 测试实例
func TestLogoURLAbleToUse(t *testing.T) {
	service.RegisterLink(link.New())
	// 测试其中一个地址是否可以访问
	// 正确可以返回的
	err := service.Link().CheckLogoCanAccess(
		context.Background(),
		"https://www.baidu.com/img/flexible/logo/pc/peak-result.png",
	)
	if err != nil {
		t.Errorf("测试失败，错误原因：%s", err.Error())
	}

	// 返回非图片
	err = service.Link().CheckLogoCanAccess(context.Background(), "https://blog.x-lf.com/")
	if err != nil {
		t.Logf("测试失败，错误原因：%s", err.Error())
	}

	// 错误的不可返回的
	err = service.Link().CheckLogoCanAccess(context.Background(), "https://blog.x-lf.com/404")
	if err != nil {
		t.Logf("测试失败，错误原因：%s", err.Error())
	}
}

// TestRSSURLAbleToUse 测试 RSS 地址是否可以访问
// 用于测试 RSS 地址是否可以访问，如果可以访问则返回 nil，否则返回错误
// 进行测试时，需要注意 RSS 地址是否为 XML 格式，否则会返回错误
//
// 参数：
// t: 测试实例
func TestRSSURLAbleToUse(t *testing.T) {
	service.RegisterLink(link.New())
	// 测试其中一个地址是否可以访问
	// 正确可以返回的
	err := service.Link().CheckRSSCanAccess(context.Background(), "https://blog.x-lf.com/atom.xml")
	if err != nil {
		t.Errorf("测试失败，错误原因：%s", err.Error())
	}

	// 错误的不可返回的
	err = service.Link().CheckRSSCanAccess(context.Background(), "https://blog.x-lf.com/404")
	if err != nil {
		t.Logf("测试失败，错误原因：%s", err.Error())
	}
}
