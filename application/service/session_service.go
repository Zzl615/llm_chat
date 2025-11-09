/**
 * @Author: Noaghzil
 * @Date:   2025-11-02 11:08:44
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2025-11-02 11:21:25
 */
package service

import (
	"llm-chat/application/dto"
	"llm-chat/domain/entity"
	"llm-chat/domain/repository"
	"llm-chat/domain/service"
	"llm-chat/domain/valueobject"
	"time"
)

// SessionApplicationService 会话应用服务
// 负责用例编排和DTO转换
type SessionApplicationService struct {
	domainService *service.SessionService
	repo          repository.SessionRepository
}

// NewSessionApplicationService 创建会话应用服务
func NewSessionApplicationService(
	domainService *service.SessionService,
	repo repository.SessionRepository,
) *SessionApplicationService {
	return &SessionApplicationService{
		domainService: domainService,
		repo:          repo,
	}
}

// CreateSession 创建会话用例
func (s *SessionApplicationService) CreateSession(req *dto.CreateSessionRequest) (*dto.SessionResponse, error) {
	var sessionID valueobject.SessionID
	var err error

	if req.SessionID != "" {
		sessionID, err = valueobject.NewSessionID(req.SessionID)
		if err != nil {
			return nil, err
		}
	} else {
		// 自动生成会话ID
		count := s.repo.Count()
		sessionID = valueobject.GenerateSessionID(count + 1)
	}

	session, err := s.domainService.CreateSession(sessionID)
	if err != nil {
		return nil, err
	}

	return s.toSessionResponse(session), nil
}

// GetSession 获取会话用例
func (s *SessionApplicationService) GetSession(sessionIDStr string) (*dto.SessionResponse, error) {
	sessionID, err := valueobject.NewSessionID(sessionIDStr)
	if err != nil {
		return nil, err
	}

	session, err := s.domainService.GetSession(sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}

	return s.toSessionResponse(session), nil
}

// CloseSession 关闭会话用例
func (s *SessionApplicationService) CloseSession(sessionIDStr string) error {
	sessionID, err := valueobject.NewSessionID(sessionIDStr)
	if err != nil {
		return err
	}

	return s.domainService.CloseSession(sessionID)
}

// toSessionResponse 转换为会话响应DTO
func (s *SessionApplicationService) toSessionResponse(session *entity.Session) *dto.SessionResponse {
	return &dto.SessionResponse{
		SessionID: session.ID.Value(),
		IsActive:  session.IsActive,
		CreatedAt: session.CreatedAt.Format(time.RFC3339),
	}
}

