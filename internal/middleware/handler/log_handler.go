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

package handler

import (
	"bamboo-main/internal/dao"
	"bamboo-main/internal/model/dto/base"
	"bamboo-main/internal/model/entity"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/google/uuid"
	"io"
	"strings"
)

// LogHandler
//
//	是一个中间件，用于记录请求日志。
//
// 它会在请求处理完成后，记录请求的路径、方法、请求体、IP 地址、用户代理和引用页面等信息。
// 如果请求处理过程中没有发生错误，则记录操作成功的日志；
// 如果发生了错误，则记录错误信息。
func LogHandler(r *ghttp.Request) {
	// 记录请求日志
	r.Middleware.Next()

	// 对以下开头的请求不记录日志
	var urlPaths = []string{
		"/api/v1/auth",
	}
	for _, path := range urlPaths {
		if strings.HasPrefix(r.URL.Path, path) {
			return
		}
	}

	// 创建日志
	readAll, _ := io.ReadAll(r.Body)
	logDTO := base.LogDTO{
		Path:    r.URL.Path,
		Method:  r.Method,
		Body:    gconv.String(readAll),
		IP:      r.GetClientIp(),
		UA:      r.Header.Get("User-Agent"),
		Referer: r.Header.Get("Referer"),
	}

	// 对请求进行分类
	if r.GetError() == nil {
		// 如果没有错误，记录成功日志
		logDTO.Msg = "操作成功"
	} else {
		// 如果有错误，记录错误日
		logDTO.Msg = r.GetError().Error()
	}

	_, daoErr := dao.Logs.Ctx(r.GetCtx()).Insert(&entity.Logs{
		LogUuid:      uuid.New().String(),
		LogType:      0,
		LogContent:   gjson.MustEncodeString(logDTO),
		LogCreatedAt: gtime.Now(),
	})
	if daoErr != nil {
		r.SetError(berror.ErrorAddData(&berror.ErrInternalServer, daoErr.Error()))
	}
}
