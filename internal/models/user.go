package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	ChatID    int64 `gorm:"unique_index"`
	Photos    uint
	CityID    uint
	City      City `gorm:"foreignkey:CityID;association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	Materials uint
}
