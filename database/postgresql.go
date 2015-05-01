package database

import (
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	pq "github.com/lib/pq"
)

func NewPostgresql() gorm.DB {
	var database gorm.DB
	var err error

	databaseUrl, _ := pq.ParseURL(os.Getenv("DATABASE_URL"))
	database, err = gorm.Open("postgres", databaseUrl)

	if err != nil {
		log.Fatalln(err.Error())
	}

	err = database.DB().Ping()
	if err != nil {
		log.Fatalln(err.Error())
	}

	database.LogMode(os.Getenv("DEBUG") == "true")

	Migrations(database)

	return database
}

func Migrations(database gorm.DB) {
	database.AutoMigrate(&Image{})

	database.Model(&Image{}).AddIndex("idx_id", "id")
	database.Model(&Image{}).AddIndex("idx_by_url", "url")
}

type Image struct {
	Id        int64
	Url       string `sql:"unique"`
	Data      []byte `sql:"type:bytea"`
	UpdatedAt time.Time
}
