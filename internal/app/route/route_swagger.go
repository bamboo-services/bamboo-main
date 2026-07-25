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
	xEnv "github.com/bamboo-services/bamboo-base-go/defined/env"
	"github.com/bamboo-services/bamboo-main/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// swaggerRegister 注册 Swagger 文档接口到指定的路由组。
//
// 该函数通过 `SwaggerInfo` 配置 API 文档的基本信息，包括版本、地址等，
// 并绑定 `/swagger/*any` 路由展示 Swagger UI 页面。
//
// 注意:
//   - `SwaggerInfo` 的字段可通过环境变量动态设置以适配不同环境。
//   - 此函数建议仅在开发调试模式下调用以避免生产环境暴露内部文档。
//
// 参数说明:
//   - r: Gin 路由组实例，用于注册 Swagger 路由。
func swaggerRegister(r gin.IRouter) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Title = "竹叶"
	docs.SwaggerInfo.Description = "筱锋的个人主页程序"
	docs.SwaggerInfo.Version = "v1.0.0"
	docs.SwaggerInfo.Host = xEnv.GetEnvString(xEnv.Host, "localhost") + ":" + xEnv.GetEnvString(xEnv.Port, "5566")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	swaggerGroup := r.Group("/swagger")
	swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
