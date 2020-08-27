package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

)

var DB *gorm.DB

func Open() (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open("postgres", "host=db port=5432 user=postgres password=password sslmode=disable")
	if err != nil {
		return nil, err
	}
	return DB, nil
}

func Close() error {
	return DB.Close()
}
