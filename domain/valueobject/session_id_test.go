/**
 * @Author: Noaghzil
 * @Date:   2025-11-02
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02
 */
package valueobject

import (
	"testing"
)

func TestNewSessionID(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid session id",
			value:   "sess-123",
			wantErr: false,
		},
		{
			name:    "empty session id",
			value:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			value:   "   ",
			wantErr: true,
		},
		{
			name:    "valid with spaces",
			value:   "  sess-123  ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSessionID(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSessionID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Value() != tt.value {
					t.Errorf("NewSessionID() value = %v, want %v", got.Value(), tt.value)
				}
			}
		})
	}
}

func TestGenerateSessionID(t *testing.T) {
	tests := []struct {
		name     string
		sequence int
		want     string
	}{
		{
			name:     "sequence 1",
			sequence: 1,
			want:     "sess-1",
		},
		{
			name:     "sequence 100",
			sequence: 100,
			want:     "sess-100",
		},
		{
			name:     "sequence 0",
			sequence: 0,
			want:     "sess-0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSessionID(tt.sequence)
			if got.Value() != tt.want {
				t.Errorf("GenerateSessionID() = %v, want %v", got.Value(), tt.want)
			}
		})
	}
}

func TestSessionID_Value(t *testing.T) {
	sid, _ := NewSessionID("test-session")
	if sid.Value() != "test-session" {
		t.Errorf("SessionID.Value() = %v, want test-session", sid.Value())
	}
}

func TestSessionID_String(t *testing.T) {
	sid, _ := NewSessionID("test-session")
	if sid.String() != "test-session" {
		t.Errorf("SessionID.String() = %v, want test-session", sid.String())
	}
}

func TestSessionID_Equals(t *testing.T) {
	sid1, _ := NewSessionID("test-session")
	sid2, _ := NewSessionID("test-session")
	sid3, _ := NewSessionID("other-session")

	if !sid1.Equals(sid2) {
		t.Error("SessionID.Equals() should return true for same values")
	}

	if sid1.Equals(sid3) {
		t.Error("SessionID.Equals() should return false for different values")
	}
}

