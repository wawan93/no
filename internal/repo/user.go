package repo

import (
	"no/internal/models"

	"github.com/jinzhu/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Get(chatID int64) (*models.User, error) {
	var user models.User
	user.ChatID = chatID

	err := r.db.
		Preload("City").
		Where("chat_id=?", chatID).
		FirstOrCreate(&user).
		Error
	return &user, err
}

func (r *UserRepo) IncrementPhotos(user *models.User) error {
	user.Photos++
	return r.db.Save(&user).Error
}

func (r *UserRepo) Update(user *models.User) error {
	return r.db.Save(&user).Error
}
