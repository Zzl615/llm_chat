/**
 * @Author: Noaghzil
 * @Date:   2025-11-02
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02
 */
package repository

import (
	"llm-chat/domain/entity"
	"llm-chat/domain/valueobject"
	"testing"
)

func TestNewSessionRepositoryImpl(t *testing.T) {
	repo := NewSessionRepositoryImpl()
	if repo == nil {
		t.Fatal("NewSessionRepositoryImpl() returned nil")
	}
}

func TestSessionRepositoryImpl_Save(t *testing.T) {
	repo := NewSessionRepositoryImpl().(*SessionRepositoryImpl)
	sessionID, _ := valueobject.NewSessionID("test-session")
	session := entity.NewSession(sessionID)

	err := repo.Save(session)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// 验证已保存
	saved, _ := repo.FindByID(sessionID)
	if saved == nil {
		t.Error("Save() session not saved")
	}

	if !saved.ID.Equals(sessionID) {
		t.Errorf("Save() ID = %v, want %v", saved.ID, sessionID)
	}
}

func TestSessionRepositoryImpl_FindByID(t *testing.T) {
	repo := NewSessionRepositoryImpl().(*SessionRepositoryImpl)
	sessionID, _ := valueobject.NewSessionID("test-session")
	session := entity.NewSession(sessionID)
	repo.Save(session)

	// 查找会话
	found, err := repo.FindByID(sessionID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if found == nil {
		t.Fatal("FindByID() returned nil")
	}

	if !found.ID.Equals(sessionID) {
		t.Errorf("FindByID() ID = %v, want %v", found.ID, sessionID)
	}
}

func TestSessionRepositoryImpl_FindByID_NotFound(t *testing.T) {
	repo := NewSessionRepositoryImpl().(*SessionRepositoryImpl)
	sessionID, _ := valueobject.NewSessionID("non-existent")

	found, err := repo.FindByID(sessionID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if found != nil {
		t.Error("FindByID() should return nil for non-existent session")
	}
}

func TestSessionRepositoryImpl_Delete(t *testing.T) {
	repo := NewSessionRepositoryImpl().(*SessionRepositoryImpl)
	sessionID, _ := valueobject.NewSessionID("test-session")
	session := entity.NewSession(sessionID)
	repo.Save(session)

	// 删除会话
	err := repo.Delete(sessionID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// 验证已删除
	found, _ := repo.FindByID(sessionID)
	if found != nil {
		t.Error("Delete() session not deleted")
	}
}

func TestSessionRepositoryImpl_FindAll(t *testing.T) {
	repo := NewSessionRepositoryImpl().(*SessionRepositoryImpl)

	// 创建多个会话
	for i := 1; i <= 3; i++ {
		sessionID, _ := valueobject.NewSessionID(valueobject.GenerateSessionID(i).Value())
		session := entity.NewSession(sessionID)
		repo.Save(session)
	}

	// 查找所有会话
	all, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}

	if len(all) != 3 {
		t.Errorf("FindAll() count = %v, want 3", len(all))
	}
}

func TestSessionRepositoryImpl_Count(t *testing.T) {
	repo := NewSessionRepositoryImpl().(*SessionRepositoryImpl)

	if repo.Count() != 0 {
		t.Errorf("Count() = %v, want 0", repo.Count())
	}

	// 创建会话
	for i := 1; i <= 5; i++ {
		sessionID, _ := valueobject.NewSessionID(valueobject.GenerateSessionID(i).Value())
		session := entity.NewSession(sessionID)
		repo.Save(session)
	}

	if repo.Count() != 5 {
		t.Errorf("Count() = %v, want 5", repo.Count())
	}
}

func TestSessionRepositoryImpl_Concurrent(t *testing.T) {
	repo := NewSessionRepositoryImpl().(*SessionRepositoryImpl)

	// 并发保存
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			sessionID, _ := valueobject.NewSessionID(valueobject.GenerateSessionID(id).Value())
			session := entity.NewSession(sessionID)
			repo.Save(session)
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	if repo.Count() != 10 {
		t.Errorf("Count() = %v, want 10", repo.Count())
	}
}

