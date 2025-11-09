/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:47
 */
package websocket

import (
	"llm-chat/application/dto"
	"llm-chat/application/service"
	"llm-chat/infrastructure/websocket"
	"llm-chat/domain/valueobject"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

// Handler WebSocket处理器
type Handler struct {
	chatService    *service.ChatApplicationService
	sessionService *service.SessionApplicationService
}

// NewHandler 创建WebSocket处理器
func NewHandler(
	chatService *service.ChatApplicationService,
	sessionService *service.SessionApplicationService,
) *Handler {
	return &Handler{
		chatService:    chatService,
		sessionService: sessionService,
	}
}

// HandleWebSocket 处理WebSocket连接
func (h *Handler) HandleWebSocket(c *gin.Context) {
	upgrader := gorilla.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(400, "upgrade failed: %v", err)
		return
	}

	// 创建会话
	sessionResp, err := h.sessionService.CreateSession(&dto.CreateSessionRequest{})
	if err != nil {
		log.Printf("[WebSocket] create session error: %v", err)
		conn.Close()
		return
	}

	sessionID, err := valueobject.NewSessionID(sessionResp.SessionID)
	if err != nil {
		log.Printf("[WebSocket] invalid session id: %v", err)
		conn.Close()
		return
	}

	// 创建WebSocket连接封装
	wsConn := websocket.NewConnection(conn)

	// 启动心跳协程
	go wsConn.HeartbeatPump()

	// 启动读协程，处理用户消息
	go wsConn.ReadPump(func(msg []byte) {
		// 发送消息到聊天服务
		chatReq := &dto.ChatRequest{
			SessionID: sessionID.Value(),
			Content:   string(msg),
		}
		if err := h.chatService.SendMessage(chatReq); err != nil {
			log.Printf("[WebSocket] send message error: %v", err)
		}
	})
}

