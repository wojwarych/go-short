package main

import (
	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	ShortURL string `gorm:"notNull"`
	LongURL  string `gorm:"notNull;unique"`
	Visits   int
}
