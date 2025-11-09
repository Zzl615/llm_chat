/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:25
 */
package valueobject

import (
	"fmt"
	"strings"
)

// SessionID 会话ID值对象
type SessionID struct {
	value string
}

// NewSessionID 创建新的会话ID
func NewSessionID(value string) (SessionID, error) {
	if strings.TrimSpace(value) == "" {
		return SessionID{}, fmt.Errorf("session id cannot be empty")
	}
	return SessionID{value: value}, nil
}

// GenerateSessionID 生成新的会话ID
func GenerateSessionID(sequence int) SessionID {
	return SessionID{value: fmt.Sprintf("sess-%d", sequence)}
}

// Value 获取会话ID的值
func (s SessionID) Value() string {
	return s.value
}

// String 实现Stringer接口
func (s SessionID) String() string {
	return s.value
}

// Equals 比较两个会话ID是否相等
func (s SessionID) Equals(other SessionID) bool {
	return s.value == other.value
}

