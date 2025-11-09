/**
 * @Author: Noaghzil
 * @Date:   2025-11-02
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02
 */
package service

import (
	"llm-chat/application/dto"
	"llm-chat/domain/entity"
	"llm-chat/domain/valueobject"
	"testing"
)

// mockSessionRepository 模拟会话仓储
type mockChatSessionRepository struct {
	sessions map[string]*entity.Session
}

func newMockChatSessionRepository() *mockChatSessionRepository {
	return &mockChatSessionRepository{
		sessions: make(map[string]*entity.Session),
	}
}

func (m *mockChatSessionRepository) Save(session *entity.Session) error {
	m.sessions[session.ID.Value()] = session
	return nil
}

func (m *mockChatSessionRepository) FindByID(id valueobject.SessionID) (*entity.Session, error) {
	session, ok := m.sessions[id.Value()]
	if !ok {
		return nil, nil
	}
	return session, nil
}

func (m *mockChatSessionRepository) Delete(id valueobject.SessionID) error {
	delete(m.sessions, id.Value())
	return nil
}

func (m *mockChatSessionRepository) FindAll() ([]*entity.Session, error) {
	sessions := make([]*entity.Session, 0, len(m.sessions))
	for _, session := range m.sessions {
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (m *mockChatSessionRepository) Count() int {
	return len(m.sessions)
}

// mockMessageQueue 模拟消息队列
type mockMessageQueue struct {
	publishedRequests []struct {
		sessionID valueobject.SessionID
		content   string
	}
}

func newMockMessageQueue() *mockMessageQueue {
	return &mockMessageQueue{
		publishedRequests: make([]struct {
			sessionID valueobject.SessionID
			content   string
		}, 0),
	}
}

func (m *mockMessageQueue) PublishRequest(sessionID valueobject.SessionID, content string) error {
	m.publishedRequests = append(m.publishedRequests, struct {
		sessionID valueobject.SessionID
		content   string
	}{sessionID, content})
	return nil
}

func (m *mockMessageQueue) SubscribeResults(handler func(sessionID valueobject.SessionID, chunk string, isLast bool)) error {
	// 模拟实现，不实际订阅
	return nil
}

func (m *mockMessageQueue) StartWorker() error {
	return nil
}

func TestNewChatApplicationService(t *testing.T) {
	repo := newMockChatSessionRepository()
	queue := newMockMessageQueue()
	service := NewChatApplicationService(repo, queue)

	if service == nil {
		t.Fatal("NewChatApplicationService() returned nil")
	}

	if service.sessionRepo != repo {
		t.Error("NewChatApplicationService() sessionRepo not set correctly")
	}

	if service.messageQueue != queue {
		t.Error("NewChatApplicationService() messageQueue not set correctly")
	}
}

func TestChatApplicationService_SendMessage(t *testing.T) {
	repo := newMockChatSessionRepository()
	queue := newMockMessageQueue()
	service := NewChatApplicationService(repo, queue)

	// 创建会话
	sessionID, _ := valueobject.NewSessionID("test-session")
	session := entity.NewSession(sessionID)
	repo.Save(session)

	// 发送消息
	req := &dto.ChatRequest{
		SessionID: "test-session",
		Content:   "Hello, world!",
	}

	err := service.SendMessage(req)
	if err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}

	// 验证消息已发布
	if len(queue.publishedRequests) != 1 {
		t.Fatalf("SendMessage() published requests = %v, want 1", len(queue.publishedRequests))
	}

	published := queue.publishedRequests[0]
	if !published.sessionID.Equals(sessionID) {
		t.Errorf("SendMessage() sessionID = %v, want %v", published.sessionID, sessionID)
	}

	if published.content != "Hello, world!" {
		t.Errorf("SendMessage() content = %v, want Hello, world!", published.content)
	}
}

func TestChatApplicationService_SendMessage_InvalidSessionID(t *testing.T) {
	repo := newMockChatSessionRepository()
	queue := newMockMessageQueue()
	service := NewChatApplicationService(repo, queue)

	req := &dto.ChatRequest{
		SessionID: "",
		Content:   "Hello, world!",
	}

	err := service.SendMessage(req)
	if err == nil {
		t.Error("SendMessage() should return error for invalid session ID")
	}
}

func TestChatApplicationService_SendMessage_SessionNotFound(t *testing.T) {
	repo := newMockChatSessionRepository()
	queue := newMockMessageQueue()
	service := NewChatApplicationService(repo, queue)

	req := &dto.ChatRequest{
		SessionID: "non-existent",
		Content:   "Hello, world!",
	}

	err := service.SendMessage(req)
	if err == nil {
		t.Error("SendMessage() should return error for non-existent session")
	}

	if _, ok := err.(*SessionNotFoundError); !ok {
		t.Errorf("SendMessage() error type = %T, want *SessionNotFoundError", err)
	}
}

func TestChatApplicationService_SendMessage_SessionNotActive(t *testing.T) {
	repo := newMockChatSessionRepository()
	queue := newMockMessageQueue()
	service := NewChatApplicationService(repo, queue)

	// 创建非活跃会话
	sessionID, _ := valueobject.NewSessionID("test-session")
	session := entity.NewSession(sessionID)
	session.Deactivate()
	repo.Save(session)

	req := &dto.ChatRequest{
		SessionID: "test-session",
		Content:   "Hello, world!",
	}

	err := service.SendMessage(req)
	if err == nil {
		t.Error("SendMessage() should return error for inactive session")
	}

	if _, ok := err.(*SessionNotActiveError); !ok {
		t.Errorf("SendMessage() error type = %T, want *SessionNotActiveError", err)
	}
}

func TestSessionNotFoundError(t *testing.T) {
	err := &SessionNotFoundError{SessionID: "test-session"}
	expected := "session not found: test-session"
	if err.Error() != expected {
		t.Errorf("SessionNotFoundError.Error() = %v, want %v", err.Error(), expected)
	}
}

func TestSessionNotActiveError(t *testing.T) {
	err := &SessionNotActiveError{SessionID: "test-session"}
	expected := "session not active: test-session"
	if err.Error() != expected {
		t.Errorf("SessionNotActiveError.Error() = %v, want %v", err.Error(), expected)
	}
}

