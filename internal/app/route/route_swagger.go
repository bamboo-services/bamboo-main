/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package route

import (
	"github.com/bamboo-services/bamboo-base-go/foundation/v2/env"
	"github.com/bamboo-services/bamboo-main/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// swaggerRegister 注册 Swagger 文档接口到指定的路由组。
func swaggerRegister(r gin.IRouter) {
	host := "localhost"
	if val, ok := env.Lookup(env.Host); ok && !val.IsEmpty() {
		host = val.String()
	}
	port := "5566"
	if val, ok := env.Lookup(env.Port); ok && !val.IsEmpty() {
		port = val.String()
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Title = "竹叶"
	docs.SwaggerInfo.Description = "筱锋的个人主页程序"
	docs.SwaggerInfo.Version = "v1.0.0"
	docs.SwaggerInfo.Host = host + ":" + port
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	swaggerGroup := r.Group("/swagger")
	swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
