package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	ChatID int64 `gorm:"unique_index"`
	Photos uint
}
