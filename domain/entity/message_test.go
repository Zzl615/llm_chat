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
)

func TestNewRequestMessage(t *testing.T) {
	sessionID, _ := valueobject.NewSessionID("test-session")
	content := "Hello, world!"
	msg := NewRequestMessage(sessionID, content)

	if msg == nil {
		t.Fatal("NewRequestMessage() returned nil")
	}

	if msg.ID == "" {
		t.Error("NewRequestMessage() ID is empty")
	}

	if !msg.SessionID.Equals(sessionID) {
		t.Errorf("NewRequestMessage() SessionID = %v, want %v", msg.SessionID, sessionID)
	}

	if msg.Content != content {
		t.Errorf("NewRequestMessage() Content = %v, want %v", msg.Content, content)
	}

	if msg.Type != MessageTypeRequest {
		t.Errorf("NewRequestMessage() Type = %v, want %v", msg.Type, MessageTypeRequest)
	}

	if msg.IsLast {
		t.Error("NewRequestMessage() IsLast = true, want false")
	}

	if !msg.IsRequest() {
		t.Error("NewRequestMessage() IsRequest() = false, want true")
	}

	if msg.IsResult() {
		t.Error("NewRequestMessage() IsResult() = true, want false")
	}
}

func TestNewResultMessage(t *testing.T) {
	sessionID, _ := valueobject.NewSessionID("test-session")
	chunk := "chunk 1"
	isLast := false
	msg := NewResultMessage(sessionID, chunk, isLast)

	if msg == nil {
		t.Fatal("NewResultMessage() returned nil")
	}

	if msg.ID == "" {
		t.Error("NewResultMessage() ID is empty")
	}

	if !msg.SessionID.Equals(sessionID) {
		t.Errorf("NewResultMessage() SessionID = %v, want %v", msg.SessionID, sessionID)
	}

	if msg.Content != chunk {
		t.Errorf("NewResultMessage() Content = %v, want %v", msg.Content, chunk)
	}

	if msg.Type != MessageTypeResult {
		t.Errorf("NewResultMessage() Type = %v, want %v", msg.Type, MessageTypeResult)
	}

	if msg.IsLast != isLast {
		t.Errorf("NewResultMessage() IsLast = %v, want %v", msg.IsLast, isLast)
	}

	if !msg.IsResult() {
		t.Error("NewResultMessage() IsResult() = false, want true")
	}

	if msg.IsRequest() {
		t.Error("NewResultMessage() IsRequest() = true, want false")
	}
}

func TestMessage_IsRequest(t *testing.T) {
	sessionID, _ := valueobject.NewSessionID("test-session")
	requestMsg := NewRequestMessage(sessionID, "test")
	resultMsg := NewResultMessage(sessionID, "test", false)

	if !requestMsg.IsRequest() {
		t.Error("Message.IsRequest() = false for request message, want true")
	}

	if resultMsg.IsRequest() {
		t.Error("Message.IsRequest() = true for result message, want false")
	}
}

func TestMessage_IsResult(t *testing.T) {
	sessionID, _ := valueobject.NewSessionID("test-session")
	requestMsg := NewRequestMessage(sessionID, "test")
	resultMsg := NewResultMessage(sessionID, "test", false)

	if !resultMsg.IsResult() {
		t.Error("Message.IsResult() = false for result message, want true")
	}

	if requestMsg.IsResult() {
		t.Error("Message.IsResult() = true for request message, want false")
	}
}

