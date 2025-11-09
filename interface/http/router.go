/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:47
 */
package http

import (
	"llm-chat/interface/sse"
	"llm-chat/interface/websocket"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func RegisterRoutes(
	r *gin.Engine,
	wsHandler *websocket.Handler,
	sseHandler *sse.Handler,
) {
	// WebSocket端点 - 只用于会话管理和接收用户消息
	r.GET("/ws", wsHandler.HandleWebSocket)

	// SSE端点 - 用于发送大模型返回结果
	r.GET("/sse/:sessionId", sseHandler.HandleSSE)
}

