/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:47
 */
package sse

import (
	"encoding/json"
	"llm-chat/application/service"
	"llm-chat/domain/repository"
	"llm-chat/domain/valueobject"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler SSE处理器
type Handler struct {
	sessionRepo repository.SessionRepository
	messageQueue service.MessageQueue
}

// NewHandler 创建SSE处理器
func NewHandler(
	sessionRepo repository.SessionRepository,
	messageQueue service.MessageQueue,
) *Handler {
	return &Handler{
		sessionRepo: sessionRepo,
		messageQueue: messageQueue,
	}
}

// HandleSSE 处理SSE连接
func (h *Handler) HandleSSE(c *gin.Context) {
	sessionIDStr := c.Param("sessionId")
	sessionID, err := valueobject.NewSessionID(sessionIDStr)
	if err != nil {
		c.String(400, "invalid session id: %s", sessionIDStr)
		return
	}

	// 验证会话是否存在
	session, err := h.sessionRepo.FindByID(sessionID)
	if err != nil {
		c.String(500, "internal error")
		return
	}
	if session == nil {
		c.String(404, "session not found: %s", sessionIDStr)
		return
	}

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // 禁用nginx缓冲

	// 发送初始连接消息
	c.SSEvent("connected", map[string]string{"sessionId": sessionIDStr})
	c.Writer.Flush()

	// 创建SSE消息通道
	messageCh := make(chan []byte, 128)
	closeCh := make(chan struct{})

	// 订阅消息队列结果
	err = h.messageQueue.SubscribeResults(func(id valueobject.SessionID, chunk string, isLast bool) {
		// 只处理当前会话的消息
		if !id.Equals(sessionID) {
			return
		}

		// 构造JSON格式的消息
		msg := map[string]interface{}{
			"chunk":  chunk,
			"isLast": isLast,
		}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			log.Printf("[SSE] session %s marshal error: %v", sessionIDStr, err)
			return
		}

		select {
		case messageCh <- msgBytes:
		default:
			log.Printf("[SSE] session %s message buffer full, drop chunk", sessionIDStr)
		}
	})
	if err != nil {
		log.Printf("[SSE] subscribe error: %v", err)
		return
	}

	// 从消息通道读取并发送数据
	ticker := time.NewTicker(30 * time.Second) // 保持连接的心跳
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-messageCh:
			if !ok {
				// channel已关闭，发送结束事件
				c.SSEvent("end", map[string]string{"message": "stream ended"})
				c.Writer.Flush()
				return
			}

			// 发送数据事件
			var data map[string]interface{}
			if err := json.Unmarshal(msg, &data); err != nil {
				// 如果不是JSON，直接作为文本发送
				c.SSEvent("message", string(msg))
			} else {
				c.SSEvent("message", data)
			}
			c.Writer.Flush()

		case <-ticker.C:
			// 发送心跳保持连接
			c.SSEvent("ping", map[string]string{"time": time.Now().Format(time.RFC3339)})
			c.Writer.Flush()

		case <-closeCh:
			log.Printf("[SSE] session %s SSE closed", sessionIDStr)
			return

		case <-c.Request.Context().Done():
			log.Printf("[SSE] client disconnected: %s", sessionIDStr)
			close(closeCh)
			return
		}
	}
}

