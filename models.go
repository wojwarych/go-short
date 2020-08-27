package main

import (
	"github.com/jinzhu/gorm"
)

type URL struct {
	gorm.Model
	ShortURL string `gorm:"not_null"`
	LongURL  string `gorm:"not_null;unique"`
}
