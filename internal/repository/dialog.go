package repository

import (
	"github.com/tamper000/freybot/internal/models"
	"gorm.io/gorm"
)

type DialogRepository interface {
	AddMessage(userID int64, role, content string) error
	GetHistory(userID int64) ([]models.Message, error)
	ClearHistory(userID int64) error
	DeleteLastMessage(userID int64) error
}

type DialogRepo struct {
	db          *gorm.DB
	maxMessages int
}

func NewDialogRepository(db *gorm.DB, maxMsg int) DialogRepository {
	return &DialogRepo{db: db, maxMessages: maxMsg}
}

func (r *DialogRepo) AddMessage(userID int64, role, content string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		msg := models.Message{
			UserID:  userID,
			Role:    role,
			Content: content,
		}
		if err := tx.Create(&msg).Error; err != nil {
			return err
		}

		var messages []models.Message
		if err := tx.Where("user_id = ?", userID).Order("created_at").Find(&messages).Error; err != nil {
			return err
		}

		if len(messages) > r.maxMessages {
			toDelete := messages[:len(messages)-r.maxMessages]
			for _, msg := range toDelete {
				if err := tx.Delete(&models.Message{}, msg.ID).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *DialogRepo) GetHistory(userID int64) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Where("user_id = ?", userID).
		Order("created_at").
		Limit(r.maxMessages).
		Find(&messages).Error

	return messages, err
}

func (r *DialogRepo) ClearHistory(userID int64) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Message{}).Error
}

func (r *DialogRepo) DeleteLastMessage(userID int64) error {
	var lastMessage models.Message
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&lastMessage).Error

	if err != nil {
		return err
	}

	return r.db.Delete(&lastMessage).Error
}
