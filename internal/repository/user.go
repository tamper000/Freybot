// internal/repository/user_repository.go
package repository

import (
	"github.com/tamper000/freybot/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	AddUser(userID int64) error
	DelUser(userID int64) error
	GetUser(userID int64) (*models.User, error)

	UpdateRole(userID int64, role string) error
	UpdateProvider(userID int64, provider string) error
	UpdateGroup(userID int64, group string) error
	UpdateTextModel(userID int64, model string) error
	UpdatePhotoModel(userID int64, model string) error
	UpdateEditModel(userID int64, model string) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) AddUser(userID int64) error {
	return r.db.Create(&models.User{ID: userID}).Error
}

func (r *UserRepo) DelUser(userID int64) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Delete(nil).Error
}

func (r *UserRepo) GetUser(userID int64) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) UpdateRole(userID int64, role string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("Role", role).Error
}

func (r *UserRepo) UpdateProvider(userID int64, provider string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("Provider", provider).Error
}

func (r *UserRepo) UpdateGroup(userID int64, group string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("Group", group).Error
}

func (r *UserRepo) UpdateTextModel(userID int64, model string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("Model", model).Error
}

func (r *UserRepo) UpdatePhotoModel(userID int64, model string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("Photo", model).Error
}

func (r *UserRepo) UpdateEditModel(userID int64, model string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("Edit", model).Error
}
