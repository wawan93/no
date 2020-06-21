package repo

import (
	"fmt"
	"no/internal/models"

	"github.com/jinzhu/gorm"
)

type TickRepo struct {
	db *gorm.DB
}

func NewTickRepo(db *gorm.DB) *TickRepo {
	return &TickRepo{
		db: db,
	}
}

func (r *TickRepo) Get(chatID int64) (*models.Tick, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *TickRepo) Save(tick *models.Tick) error {
	return r.db.Save(tick).Error
}
