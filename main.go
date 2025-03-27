package main

import (
	"JenArgo/controller"
	"JenArgo/db"
	"JenArgo/logger"
	"JenArgo/middleware/snowflake"
	"JenArgo/service"
	"JenArgo/settings"
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	// 1. 加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}

	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")

	// 3. 初始化MySQL连接
	if err := db.Init(settings.Conf.MySQLConfig); err != nil {
		zap.L().Error("init mysql failed, err:%v\n", zap.Error(err))
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			zap.L().Fatal("数据库关闭异常:", zap.Error(err))
		}
	}()

	// 4. 注册雪花算法 ID 生成器
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	// 5. 启动定时获取 ArgoCD token
	go func() {
		ticker := time.NewTicker(12 * time.Hour)
		defer ticker.Stop()

		// 首次运行
		if err := service.ArgoCD.Session(); err != nil {
			zap.L().Error("Get ArgoCD session failed", zap.Error(err))
		}

		for range ticker.C {
			if err := service.ArgoCD.Session(); err != nil {
				zap.L().Error("Get ArgoCD session failed", zap.Error(err))
			}
		}
	}()

	// 6. 注册路由
	r := controller.Router.Setup()

	// 7. gin server 启动
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s", zap.Error(err))
		}
	}()

	// 8. 优雅关闭server
	// 声明一个系统信号的channel，并监听他，如果没有信号，就一直阻塞，如果有，就继续执行
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 9. 设置ctx超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//cancel用于释放ctx
	defer cancel()

	// 10. 关闭 gin server
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Gin Server 关闭异常：", zap.Error(err))
	}
	zap.L().Info("Gin Server 退出成功")
}
