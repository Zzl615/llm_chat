package server

import (
	"llm-chat/internal/service"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer 创建HTTP服务器
func NewHTTPServer(chatService *service.ChatService, addr string) *http.Server {
	srv := http.NewServer(
		http.Address(addr),
	)
	
	// 注册 WebSocket 路由
	srv.HandleFunc("/ws", chatService.HandleWebSocket)
	
	return srv
}

