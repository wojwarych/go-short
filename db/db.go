package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = "host=db port=5432 user=postgres password=password sslmode=disable"

func Open() (*gorm.DB, error) {
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return DB, nil
}
