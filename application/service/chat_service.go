/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:25
 */
package service

import (
	"llm-chat/application/dto"
	"llm-chat/domain/repository"
	"llm-chat/domain/valueobject"
)

// MessageQueue 消息队列接口（定义在应用层，由基础设施层实现）
type MessageQueue interface {
	// PublishRequest 发布请求消息
	PublishRequest(sessionID valueobject.SessionID, content string) error

	// SubscribeResults 订阅结果消息
	SubscribeResults(handler func(sessionID valueobject.SessionID, chunk string, isLast bool)) error

	// StartWorker 启动工作协程
	StartWorker() error
}

// ChatApplicationService 聊天应用服务
// 负责聊天用例编排
type ChatApplicationService struct {
	sessionRepo repository.SessionRepository
	messageQueue MessageQueue
}

// NewChatApplicationService 创建聊天应用服务
func NewChatApplicationService(
	sessionRepo repository.SessionRepository,
	messageQueue MessageQueue,
) *ChatApplicationService {
	return &ChatApplicationService{
		sessionRepo: sessionRepo,
		messageQueue: messageQueue,
	}
}

// SendMessage 发送消息用例
func (s *ChatApplicationService) SendMessage(req *dto.ChatRequest) error {
	// 验证会话是否存在
	sessionID, err := valueobject.NewSessionID(req.SessionID)
	if err != nil {
		return err
	}

	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return &SessionNotFoundError{SessionID: req.SessionID}
	}
	if !session.IsValid() {
		return &SessionNotActiveError{SessionID: req.SessionID}
	}

	// 发布消息到队列
	return s.messageQueue.PublishRequest(sessionID, req.Content)
}

// SessionNotFoundError 会话未找到错误
type SessionNotFoundError struct {
	SessionID string
}

func (e *SessionNotFoundError) Error() string {
	return "session not found: " + e.SessionID
}

// SessionNotActiveError 会话未激活错误
type SessionNotActiveError struct {
	SessionID string
}

func (e *SessionNotActiveError) Error() string {
	return "session not active: " + e.SessionID
}

