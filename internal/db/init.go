package db

import (
	"log"
	"no/internal/models"

	"github.com/jinzhu/gorm"
)

var (
	Conn *gorm.DB
)

func Connect(dialect, host, port, user, pass, database string) {
	log.Println("Сonnect to database")
	conn, err := gorm.Open(
		dialect,
		user+
			":"+pass+
			"@tcp("+host+
			":"+port+")"+
			"/"+database+
			"?charset=utf8mb4&parseTime=true",
	)
	if err != nil {
		log.Panic(err)
	}
	Conn = conn
}

func Mock() {
	Conn = &gorm.DB{}
}

func Migrate() {
	Conn.AutoMigrate(
		new(models.User),
	)
}
