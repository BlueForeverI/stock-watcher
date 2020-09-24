package main

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email  string
	Stocks []Stock `gorm:"many2many:user_stocks;"`
}
