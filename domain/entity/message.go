/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:25
 */
package entity

import (
	"llm-chat/domain/valueobject"
	"time"
)

// MessageType 消息类型
type MessageType string

const (
	MessageTypeRequest MessageType = "request" // 用户请求
	MessageTypeResult  MessageType = "result"  // 模型响应
)

// Message 消息领域实体
type Message struct {
	ID        string
	SessionID valueobject.SessionID
	Type      MessageType
	Content   string
	IsLast    bool
	CreatedAt time.Time
}

// NewRequestMessage 创建请求消息
func NewRequestMessage(sessionID valueobject.SessionID, content string) *Message {
	return &Message{
		ID:        generateMessageID(),
		SessionID: sessionID,
		Type:      MessageTypeRequest,
		Content:   content,
		IsLast:    false,
		CreatedAt: time.Now(),
	}
}

// NewResultMessage 创建结果消息
func NewResultMessage(sessionID valueobject.SessionID, chunk string, isLast bool) *Message {
	return &Message{
		ID:        generateMessageID(),
		SessionID: sessionID,
		Type:      MessageTypeResult,
		Content:   chunk,
		IsLast:    isLast,
		CreatedAt: time.Now(),
	}
}

// IsRequest 判断是否为请求消息
func (m *Message) IsRequest() bool {
	return m.Type == MessageTypeRequest
}

// IsResult 判断是否为结果消息
func (m *Message) IsResult() bool {
	return m.Type == MessageTypeResult
}

// generateMessageID 生成消息ID
func generateMessageID() string {
	return time.Now().Format("20060102150405.000000")
}

