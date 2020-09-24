package main

import (
	"gorm.io/gorm"
)

type Stock struct {
	gorm.Model
	Name   string
	Symbol string
}
