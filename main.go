/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:32
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:17:51
 */
package main

import (
	applicationService "llm-chat/application/service"
	domainService "llm-chat/domain/service"
	"llm-chat/infrastructure/queue"
	infraRepo "llm-chat/infrastructure/repository"
	httpRouter "llm-chat/interface/http"
	"llm-chat/interface/sse"
	"llm-chat/interface/websocket"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化基础设施层
	sessionRepo := infraRepo.NewSessionRepositoryImpl()
	messageQueue := queue.NewMockQueue()

	// 启动消息队列工作协程
	if err := messageQueue.StartWorker(); err != nil {
		log.Fatalf("Failed to start message queue worker: %v", err)
	}

	// 初始化领域层
	domainSessionService := domainService.NewSessionService(sessionRepo)

	// 初始化应用层
	sessionAppService := applicationService.NewSessionApplicationService(domainSessionService, sessionRepo)
	chatAppService := applicationService.NewChatApplicationService(sessionRepo, messageQueue)

	// 注意：消息队列的结果订阅在SSE处理器中完成
	// 每个SSE连接都会独立订阅消息队列，只接收对应会话的消息

	// 初始化接口层
	wsHandler := websocket.NewHandler(chatAppService, sessionAppService)
	sseHandler := sse.NewHandler(sessionRepo, messageQueue)

	// 注册路由
	r := gin.Default()
	httpRouter.RegisterRoutes(r, wsHandler, sseHandler)

	log.Println("🚀 Chat demo server started at :8080")
	r.Run(":8080")
}
