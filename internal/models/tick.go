package models

import (
	"github.com/jinzhu/gorm"
)

type Tick struct {
	gorm.Model
	Longitude float64
	Latitude  float64
	Photo     string
	UserID    uint
	CityID    uint
	User      User `gorm:"foreignkey:UserID;association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
	City      City `gorm:"foreignkey:CityID;association_autoupdate:false;association_autocreate:false;association_save_reference:false"`
}
