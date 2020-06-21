package repo

import (
	"github.com/jinzhu/gorm"

	"no/internal/models"
)

type CityRepo struct {
	db *gorm.DB
}

func NewCityRepo(db *gorm.DB) *CityRepo {
	return &CityRepo{
		db: db,
	}
}

func (r *CityRepo) Find(city, region string) (*models.City, error) {
	var out models.City
	err := r.db.Where("name=? AND region=?", city, region).Find(&out).Error
	return &out, err
}

func (r *CityRepo) Regions() (regions []*models.City, err error) {
	err = r.db.Select("DISTINCT region").Table("cities").Find(&regions).Error
	return
}

func (r *CityRepo) Cities(region string) (cities []*models.City, err error) {
	err = r.db.Where("region=?", region).Table("cities").Find(&cities).Error
	return
}
