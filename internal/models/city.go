package models

import "github.com/jinzhu/gorm"

type City struct {
	gorm.Model
	Name   string `gorm:"index"`
	Region string `gorm:"index"`
}
