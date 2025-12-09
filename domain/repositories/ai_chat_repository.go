package repositories

import (
	"context"

	"github.com/google/uuid"

	"gofiber-template/domain/models"
)

type AIChatSessionRepository interface {
	Create(ctx context.Context, session *models.AIChatSession) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.AIChatSession, error)
	GetByIDWithMessages(ctx context.Context, id uuid.UUID) (*models.AIChatSession, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.AIChatSession, error)
	Update(ctx context.Context, id uuid.UUID, session *models.AIChatSession) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}

type AIChatMessageRepository interface {
	Create(ctx context.Context, message *models.AIChatMessage) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.AIChatMessage, error)
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*models.AIChatMessage, error)
	GetRecentBySessionID(ctx context.Context, sessionID uuid.UUID, limit int) ([]*models.AIChatMessage, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteBySessionID(ctx context.Context, sessionID uuid.UUID) error
	CountBySessionID(ctx context.Context, sessionID uuid.UUID) (int64, error)
}
