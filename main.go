package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "bamboo-main/docs"
	"bamboo-main/internal/router"
	"bamboo-main/pkg/startup"

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
	getGin := getServ.Serv.Serve
	router.Init(getGin, getServ)

	// 创建上下文和取消函数
	ctx, cancel := context.WithCancel(getServ.Serv.Context)
	defer cancel()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:    ":23333", // 使用配置文件中指定的端口
		Handler: getGin,
	}

	// 获取日志记录器
	log := getServ.Serv.Logger.Sugar().Named(xConsts.LogMAIN)

	// 启动服务器
	go func() {
		<-sigChan
		cancel()
		log.Warn("正在关闭 HTTP 服务器...")
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("HTTP 服务器关闭失败: %v", err)
		}
	}()

	// Gin 服务启动
	log.Info("正在启动 HTTP 服务器，监听端口: 23333")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP 服务器启动失败: %v", err)
	}
}
