package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"llm-chat/internal/biz"
	"llm-chat/internal/conf"
	"llm-chat/internal/server"
	"llm-chat/internal/service"
)

var (
	// Name 应用名称
	Name = "llm-chat"
	// Version 应用版本
	Version = "v1.0.0"
	// flagconf 配置文件路径
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	
	// 加载配置
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// 初始化日志
	logger := kratoslog.NewStdLogger(os.Stdout)
	log := kratoslog.NewHelper(logger)

	// 初始化业务组件
	manager := biz.NewManager()
	queue := biz.NewMockQueue()
	
	// 启动模拟的模型推理流
	queue.StartMockModelWorker()

	// 初始化服务
	chatService := service.NewChatService(manager, queue)

	// 创建 HTTP 服务器
	httpSrv := server.NewHTTPServer(chatService, bc.Server.HTTP.Addr)

	// 创建应用
	app := kratos.New(
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			httpSrv,
		),
	)

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动应用
	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("failed to start app: %v", err)
		}
	}()

	log.Infof("🚀 Chat demo server started at %s", bc.Server.HTTP.Addr)

	// 等待中断信号
	<-quit
	log.Info("Shutting down server...")

	// 优雅关闭
	if err := app.Stop(); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited")
}

