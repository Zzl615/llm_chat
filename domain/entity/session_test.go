/**
 * @Author: Noaghzil
 * @Date:   2025-11-02
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02
 */
package entity

import (
	"llm-chat/domain/valueobject"
	"testing"
	"time"
)

func TestNewSession(t *testing.T) {
	sessionID, _ := valueobject.NewSessionID("test-session")
	session := NewSession(sessionID)

	if session == nil {
		t.Fatal("NewSession() returned nil")
	}

	if !session.ID.Equals(sessionID) {
		t.Errorf("NewSession() ID = %v, want %v", session.ID, sessionID)
	}

	if !session.IsActive {
		t.Error("NewSession() IsActive = false, want true")
	}

	if session.CreatedAt.IsZero() {
		t.Error("NewSession() CreatedAt is zero")
	}

	if session.UpdatedAt.IsZero() {
		t.Error("NewSession() UpdatedAt is zero")
	}
}

func TestSession_Activate(t *testing.T) {
	sessionID, _ := valueobject.NewSessionID("test-session")
	session := NewSession(sessionID)
	session.Deactivate()

	beforeUpdate := session.UpdatedAt
	time.Sleep(10 * time.Millisecond) // 确保时间不同

	session.Activate()

	if !session.IsActive {
		t.Error("Session.Activate() IsActive = false, want true")
	}

	if !session.UpdatedAt.After(beforeUpdate) {
		t.Error("Session.Activate() UpdatedAt should be updated")
	}
}

func TestSession_Deactivate(t *testing.T) {
	sessionID, _ := valueobject.NewSessionID("test-session")
	session := NewSession(sessionID)

	beforeUpdate := session.UpdatedAt
	time.Sleep(10 * time.Millisecond) // 确保时间不同

	session.Deactivate()

	if session.IsActive {
		t.Error("Session.Deactivate() IsActive = true, want false")
	}

	if !session.UpdatedAt.After(beforeUpdate) {
		t.Error("Session.Deactivate() UpdatedAt should be updated")
	}
}

func TestSession_IsValid(t *testing.T) {
	sessionID, _ := valueobject.NewSessionID("test-session")
	session := NewSession(sessionID)

	if !session.IsValid() {
		t.Error("Session.IsValid() = false, want true for active session")
	}

	session.Deactivate()

	if session.IsValid() {
		t.Error("Session.IsValid() = true, want false for inactive session")
	}
}

