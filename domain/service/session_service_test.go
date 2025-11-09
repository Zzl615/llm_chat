/**
 * @Author: Noaghzil
 * @Date:   2025-11-02
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02
 */
package service

import (
	"llm-chat/domain/entity"
	"llm-chat/domain/valueobject"
	"testing"
)

// mockSessionRepository 模拟会话仓储
type mockSessionRepository struct {
	sessions map[string]*entity.Session
}

func newMockSessionRepository() *mockSessionRepository {
	return &mockSessionRepository{
		sessions: make(map[string]*entity.Session),
	}
}

func (m *mockSessionRepository) Save(session *entity.Session) error {
	m.sessions[session.ID.Value()] = session
	return nil
}

func (m *mockSessionRepository) FindByID(id valueobject.SessionID) (*entity.Session, error) {
	session, ok := m.sessions[id.Value()]
	if !ok {
		return nil, nil
	}
	return session, nil
}

func (m *mockSessionRepository) Delete(id valueobject.SessionID) error {
	delete(m.sessions, id.Value())
	return nil
}

func (m *mockSessionRepository) FindAll() ([]*entity.Session, error) {
	sessions := make([]*entity.Session, 0, len(m.sessions))
	for _, session := range m.sessions {
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (m *mockSessionRepository) Count() int {
	return len(m.sessions)
}

func TestNewSessionService(t *testing.T) {
	repo := newMockSessionRepository()
	service := NewSessionService(repo)

	if service == nil {
		t.Fatal("NewSessionService() returned nil")
	}

	if service.repo != repo {
		t.Error("NewSessionService() repo not set correctly")
	}
}

func TestSessionService_CreateSession(t *testing.T) {
	repo := newMockSessionRepository()
	service := NewSessionService(repo)

	sessionID, _ := valueobject.NewSessionID("test-session")

	// 创建新会话
	session, err := service.CreateSession(sessionID)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	if session == nil {
		t.Fatal("CreateSession() returned nil")
	}

	if !session.ID.Equals(sessionID) {
		t.Errorf("CreateSession() ID = %v, want %v", session.ID, sessionID)
	}

	if !session.IsActive {
		t.Error("CreateSession() IsActive = false, want true")
	}

	// 验证已保存
	saved, _ := repo.FindByID(sessionID)
	if saved == nil {
		t.Error("CreateSession() session not saved")
	}
}

func TestSessionService_CreateSession_ExistingActive(t *testing.T) {
	repo := newMockSessionRepository()
	service := NewSessionService(repo)

	sessionID, _ := valueobject.NewSessionID("test-session")
	existingSession := entity.NewSession(sessionID)
	repo.Save(existingSession)

	// 创建已存在的活跃会话
	session, err := service.CreateSession(sessionID)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	if session != existingSession {
		t.Error("CreateSession() should return existing session")
	}
}

func TestSessionService_CreateSession_ExistingInactive(t *testing.T) {
	repo := newMockSessionRepository()
	service := NewSessionService(repo)

	sessionID, _ := valueobject.NewSessionID("test-session")
	existingSession := entity.NewSession(sessionID)
	existingSession.Deactivate()
	repo.Save(existingSession)

	// 创建已存在但非活跃的会话
	session, err := service.CreateSession(sessionID)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	if !session.IsActive {
		t.Error("CreateSession() should activate existing inactive session")
	}
}

func TestSessionService_GetSession(t *testing.T) {
	repo := newMockSessionRepository()
	service := NewSessionService(repo)

	sessionID, _ := valueobject.NewSessionID("test-session")
	session := entity.NewSession(sessionID)
	repo.Save(session)

	// 获取会话
	got, err := service.GetSession(sessionID)
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}

	if got == nil {
		t.Fatal("GetSession() returned nil")
	}

	if !got.ID.Equals(sessionID) {
		t.Errorf("GetSession() ID = %v, want %v", got.ID, sessionID)
	}
}

func TestSessionService_GetSession_NotFound(t *testing.T) {
	repo := newMockSessionRepository()
	service := NewSessionService(repo)

	sessionID, _ := valueobject.NewSessionID("non-existent")

	got, err := service.GetSession(sessionID)
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}

	if got != nil {
		t.Error("GetSession() should return nil for non-existent session")
	}
}

func TestSessionService_GetSession_Inactive(t *testing.T) {
	repo := newMockSessionRepository()
	service := NewSessionService(repo)

	sessionID, _ := valueobject.NewSessionID("test-session")
	session := entity.NewSession(sessionID)
	session.Deactivate()
	repo.Save(session)

	got, err := service.GetSession(sessionID)
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}

	if got != nil {
		t.Error("GetSession() should return nil for inactive session")
	}
}

func TestSessionService_CloseSession(t *testing.T) {
	repo := newMockSessionRepository()
	service := NewSessionService(repo)

	sessionID, _ := valueobject.NewSessionID("test-session")
	session := entity.NewSession(sessionID)
	repo.Save(session)

	// 关闭会话
	err := service.CloseSession(sessionID)
	if err != nil {
		t.Fatalf("CloseSession() error = %v", err)
	}

	// 验证会话已停用
	saved, _ := repo.FindByID(sessionID)
	if saved.IsActive {
		t.Error("CloseSession() should deactivate session")
	}
}

func TestSessionService_CloseSession_NotFound(t *testing.T) {
	repo := newMockSessionRepository()
	service := NewSessionService(repo)

	sessionID, _ := valueobject.NewSessionID("non-existent")

	// 关闭不存在的会话应该不报错
	err := service.CloseSession(sessionID)
	if err != nil {
		t.Fatalf("CloseSession() error = %v", err)
	}
}

func TestSessionService_DeleteSession(t *testing.T) {
	repo := newMockSessionRepository()
	service := NewSessionService(repo)

	sessionID, _ := valueobject.NewSessionID("test-session")
	session := entity.NewSession(sessionID)
	repo.Save(session)

	// 删除会话
	err := service.DeleteSession(sessionID)
	if err != nil {
		t.Fatalf("DeleteSession() error = %v", err)
	}

	// 验证会话已删除
	saved, _ := repo.FindByID(sessionID)
	if saved != nil {
		t.Error("DeleteSession() should delete session")
	}
}

