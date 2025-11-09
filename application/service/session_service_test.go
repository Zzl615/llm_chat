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
	"llm-chat/domain/service"
	"llm-chat/domain/valueobject"
	"testing"
)

// mockSessionRepositoryForApp 模拟会话仓储
type mockSessionRepositoryForApp struct {
	sessions map[string]*entity.Session
}

func newMockSessionRepositoryForApp() *mockSessionRepositoryForApp {
	return &mockSessionRepositoryForApp{
		sessions: make(map[string]*entity.Session),
	}
}

func (m *mockSessionRepositoryForApp) Save(session *entity.Session) error {
	m.sessions[session.ID.Value()] = session
	return nil
}

func (m *mockSessionRepositoryForApp) FindByID(id valueobject.SessionID) (*entity.Session, error) {
	session, ok := m.sessions[id.Value()]
	if !ok {
		return nil, nil
	}
	return session, nil
}

func (m *mockSessionRepositoryForApp) Delete(id valueobject.SessionID) error {
	delete(m.sessions, id.Value())
	return nil
}

func (m *mockSessionRepositoryForApp) FindAll() ([]*entity.Session, error) {
	sessions := make([]*entity.Session, 0, len(m.sessions))
	for _, session := range m.sessions {
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (m *mockSessionRepositoryForApp) Count() int {
	return len(m.sessions)
}

func TestNewSessionApplicationService(t *testing.T) {
	repo := newMockSessionRepositoryForApp()
	domainService := service.NewSessionService(repo)
	appService := NewSessionApplicationService(domainService, repo)

	if appService == nil {
		t.Fatal("NewSessionApplicationService() returned nil")
	}

	if appService.domainService != domainService {
		t.Error("NewSessionApplicationService() domainService not set correctly")
	}

	if appService.repo != repo {
		t.Error("NewSessionApplicationService() repo not set correctly")
	}
}

func TestSessionApplicationService_CreateSession_WithID(t *testing.T) {
	repo := newMockSessionRepositoryForApp()
	domainService := service.NewSessionService(repo)
	appService := NewSessionApplicationService(domainService, repo)

	req := &dto.CreateSessionRequest{
		SessionID: "test-session",
	}

	resp, err := appService.CreateSession(req)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	if resp == nil {
		t.Fatal("CreateSession() returned nil")
	}

	if resp.SessionID != "test-session" {
		t.Errorf("CreateSession() SessionID = %v, want test-session", resp.SessionID)
	}

	if !resp.IsActive {
		t.Error("CreateSession() IsActive = false, want true")
	}

	if resp.CreatedAt == "" {
		t.Error("CreateSession() CreatedAt is empty")
	}
}

func TestSessionApplicationService_CreateSession_AutoGenerateID(t *testing.T) {
	repo := newMockSessionRepositoryForApp()
	domainService := service.NewSessionService(repo)
	appService := NewSessionApplicationService(domainService, repo)

	req := &dto.CreateSessionRequest{}

	resp, err := appService.CreateSession(req)
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	if resp == nil {
		t.Fatal("CreateSession() returned nil")
	}

	if resp.SessionID == "" {
		t.Error("CreateSession() SessionID is empty")
	}

	// 验证ID格式
	if resp.SessionID[:5] != "sess-" {
		t.Errorf("CreateSession() SessionID = %v, should start with sess-", resp.SessionID)
	}
}

func TestSessionApplicationService_CreateSession_InvalidID(t *testing.T) {
	repo := newMockSessionRepositoryForApp()
	domainService := service.NewSessionService(repo)
	appService := NewSessionApplicationService(domainService, repo)

	req := &dto.CreateSessionRequest{
		SessionID: "   ", // 只有空格
	}

	_, err := appService.CreateSession(req)
	if err == nil {
		t.Error("CreateSession() should return error for invalid session ID")
	}
}

func TestSessionApplicationService_GetSession(t *testing.T) {
	repo := newMockSessionRepositoryForApp()
	domainService := service.NewSessionService(repo)
	appService := NewSessionApplicationService(domainService, repo)

	// 创建会话
	sessionID, _ := valueobject.NewSessionID("test-session")
	session := entity.NewSession(sessionID)
	repo.Save(session)

	// 获取会话
	resp, err := appService.GetSession("test-session")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}

	if resp == nil {
		t.Fatal("GetSession() returned nil")
	}

	if resp.SessionID != "test-session" {
		t.Errorf("GetSession() SessionID = %v, want test-session", resp.SessionID)
	}
}

func TestSessionApplicationService_GetSession_NotFound(t *testing.T) {
	repo := newMockSessionRepositoryForApp()
	domainService := service.NewSessionService(repo)
	appService := NewSessionApplicationService(domainService, repo)

	resp, err := appService.GetSession("non-existent")
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}

	if resp != nil {
		t.Error("GetSession() should return nil for non-existent session")
	}
}

func TestSessionApplicationService_GetSession_InvalidID(t *testing.T) {
	repo := newMockSessionRepositoryForApp()
	domainService := service.NewSessionService(repo)
	appService := NewSessionApplicationService(domainService, repo)

	_, err := appService.GetSession("")
	if err == nil {
		t.Error("GetSession() should return error for invalid session ID")
	}
}

func TestSessionApplicationService_CloseSession(t *testing.T) {
	repo := newMockSessionRepositoryForApp()
	domainService := service.NewSessionService(repo)
	appService := NewSessionApplicationService(domainService, repo)

	// 创建会话
	sessionID, _ := valueobject.NewSessionID("test-session")
	session := entity.NewSession(sessionID)
	repo.Save(session)

	// 关闭会话
	err := appService.CloseSession("test-session")
	if err != nil {
		t.Fatalf("CloseSession() error = %v", err)
	}

	// 验证会话已停用
	saved, _ := repo.FindByID(sessionID)
	if saved.IsActive {
		t.Error("CloseSession() should deactivate session")
	}
}

func TestSessionApplicationService_CloseSession_InvalidID(t *testing.T) {
	repo := newMockSessionRepositoryForApp()
	domainService := service.NewSessionService(repo)
	appService := NewSessionApplicationService(domainService, repo)

	err := appService.CloseSession("")
	if err == nil {
		t.Error("CloseSession() should return error for invalid session ID")
	}
}

