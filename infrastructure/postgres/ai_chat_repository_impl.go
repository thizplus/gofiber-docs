package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

// AIChatSessionRepositoryImpl
type AIChatSessionRepositoryImpl struct {
	db *gorm.DB
}

func NewAIChatSessionRepository(db *gorm.DB) repositories.AIChatSessionRepository {
	return &AIChatSessionRepositoryImpl{db: db}
}

func (r *AIChatSessionRepositoryImpl) Create(ctx context.Context, session *models.AIChatSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *AIChatSessionRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.AIChatSession, error) {
	var session models.AIChatSession
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AIChatSessionRepositoryImpl) GetByIDWithMessages(ctx context.Context, id uuid.UUID) (*models.AIChatSession, error) {
	var session models.AIChatSession
	err := r.db.WithContext(ctx).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("id = ?", id).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AIChatSessionRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.AIChatSession, error) {
	var sessions []*models.AIChatSession
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&sessions).Error
	return sessions, err
}

func (r *AIChatSessionRepositoryImpl) Update(ctx context.Context, id uuid.UUID, session *models.AIChatSession) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(session).Error
}

func (r *AIChatSessionRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete messages first
		if err := tx.Where("session_id = ?", id).Delete(&models.AIChatMessage{}).Error; err != nil {
			return err
		}
		// Delete session
		return tx.Where("id = ?", id).Delete(&models.AIChatSession{}).Error
	})
}

func (r *AIChatSessionRepositoryImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get all session IDs
		var sessionIDs []uuid.UUID
		if err := tx.Model(&models.AIChatSession{}).
			Where("user_id = ?", userID).
			Pluck("id", &sessionIDs).Error; err != nil {
			return err
		}

		if len(sessionIDs) == 0 {
			return nil
		}

		// Delete all messages
		if err := tx.Where("session_id IN ?", sessionIDs).Delete(&models.AIChatMessage{}).Error; err != nil {
			return err
		}

		// Delete all sessions
		return tx.Where("user_id = ?", userID).Delete(&models.AIChatSession{}).Error
	})
}

func (r *AIChatSessionRepositoryImpl) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AIChatSession{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// AIChatMessageRepositoryImpl
type AIChatMessageRepositoryImpl struct {
	db *gorm.DB
}

func NewAIChatMessageRepository(db *gorm.DB) repositories.AIChatMessageRepository {
	return &AIChatMessageRepositoryImpl{db: db}
}

func (r *AIChatMessageRepositoryImpl) Create(ctx context.Context, message *models.AIChatMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *AIChatMessageRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.AIChatMessage, error) {
	var message models.AIChatMessage
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *AIChatMessageRepositoryImpl) GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*models.AIChatMessage, error) {
	var messages []*models.AIChatMessage
	err := r.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

func (r *AIChatMessageRepositoryImpl) GetRecentBySessionID(ctx context.Context, sessionID uuid.UUID, limit int) ([]*models.AIChatMessage, error) {
	var messages []*models.AIChatMessage
	err := r.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, err
}

func (r *AIChatMessageRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.AIChatMessage{}).Error
}

func (r *AIChatMessageRepositoryImpl) DeleteBySessionID(ctx context.Context, sessionID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("session_id = ?", sessionID).Delete(&models.AIChatMessage{}).Error
}

func (r *AIChatMessageRepositoryImpl) CountBySessionID(ctx context.Context, sessionID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AIChatMessage{}).Where("session_id = ?", sessionID).Count(&count).Error
	return count, err
}
