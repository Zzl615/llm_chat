/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:25
 */
package service

import (
	"llm-chat/domain/entity"
	"llm-chat/domain/repository"
	"llm-chat/domain/valueobject"
)

// SessionService 会话领域服务
// 包含跨实体的业务逻辑
type SessionService struct {
	repo repository.SessionRepository
}

// NewSessionService 创建会话领域服务
func NewSessionService(repo repository.SessionRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}

// CreateSession 创建新会话
func (s *SessionService) CreateSession(id valueobject.SessionID) (*entity.Session, error) {
	// 检查会话是否已存在
	existing, err := s.repo.FindByID(id)
	if err == nil && existing != nil {
		// 如果已存在且有效，则激活它
		if existing.IsValid() {
			return existing, nil
		}
		// 如果已存在但无效，则重新激活
		existing.Activate()
		return existing, s.repo.Save(existing)
	}

	// 创建新会话
	session := entity.NewSession(id)
	return session, s.repo.Save(session)
}

// GetSession 获取会话
func (s *SessionService) GetSession(id valueobject.SessionID) (*entity.Session, error) {
	session, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}
	if !session.IsValid() {
		return nil, nil
	}
	return session, nil
}

// CloseSession 关闭会话
func (s *SessionService) CloseSession(id valueobject.SessionID) error {
	session, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if session == nil {
		return nil
	}
	session.Deactivate()
	return s.repo.Save(session)
}

// DeleteSession 删除会话
func (s *SessionService) DeleteSession(id valueobject.SessionID) error {
	return s.repo.Delete(id)
}

