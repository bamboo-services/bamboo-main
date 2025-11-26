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

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "bamboo-main/docs"
	"bamboo-main/internal/router"
	"bamboo-main/pkg/startup"
	"fmt"

	xConsts "github.com/bamboo-services/bamboo-base-go/constants"
	xInit "github.com/bamboo-services/bamboo-base-go/init"
)

// @title BambooMain
// @version v1.0.0
// @description 友情链接管理系统
// @termsOfService https://www.aiawaken.top/
// @contact.name 筱锋 xiao_lfeng
// @contact.url https://www.x-lf.com/
// @contact.email gm@x-lf.cn
// @host localhost:23333
// @BasePath /api/v1
func main() {
	// 配置注册 - 两层初始化模式
	getServ := startup.Register(xInit.Register())

	// 变量赋值
	log := getServ.Serv.Logger.Sugar().Named(xConsts.LogMAIN)
	getGin := getServ.Serv.Serve

	// 初始化路由表
	router.Init(getGin, getServ.Config)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", getServ.Config.Xlf.Server.Port),
		Handler: getGin,
	}

	// 在独立的 goroutine 中启动服务器
	go func() {
		log.Infof("启动程序端口: %d", getServ.Config.Xlf.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[MAIN] 系统启动失败: %v", err)
		}
	}()

	// 监听系统信号，优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("收到关闭信号，正在优雅关闭...")

	// 停止邮件工作协程
	if mailWorker := startup.GetMailWorker(); mailWorker != nil {
		mailWorker.Stop()
	}

	// 设置 5 秒超时关闭 HTTP 服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("HTTP 服务器关闭失败: %v", err)
	}

	log.Info("服务已停止")
}
