package services

import (
	"context"
	"io"

	"github.com/google/uuid"

	"gofiber-template/domain/dto"
)

type AIService interface {
	// AI Quick Search (single response)
	AISearch(ctx context.Context, userID uuid.UUID, req *dto.AISearchRequest) (*dto.AISearchResponse, error)

	// AI Chat Sessions
	CreateChatSession(ctx context.Context, userID uuid.UUID, req *dto.CreateAIChatRequest) (*dto.AIChatSessionDetailResponse, error)
	GetChatSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) (*dto.AIChatSessionDetailResponse, error)
	GetChatSessions(ctx context.Context, userID uuid.UUID, req *dto.GetAIChatSessionsRequest) (*dto.AIChatSessionListResponse, error)
	DeleteChatSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error
	ClearAllChatSessions(ctx context.Context, userID uuid.UUID) error

	// AI Chat Messages
	SendMessage(ctx context.Context, userID uuid.UUID, req *dto.SendAIChatMessageRequest) (*dto.AIChatMessageResponse, error)
	SendMessageStream(ctx context.Context, userID uuid.UUID, req *dto.SendAIChatMessageRequest, writer io.Writer) error

	// AI Search Streaming
	AISearchStream(ctx context.Context, userID uuid.UUID, req *dto.AISearchRequest, writer io.Writer) error
}
