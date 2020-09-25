package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type LoginRequest struct {
	Email string `json:"email"`
}

func login(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var existingUser User
		var request LoginRequest

		json.NewDecoder(r.Body).Decode(&request)
		searchResult := db.First(&existingUser, "Email = ?", request.Email)

		if errors.Is(searchResult.Error, gorm.ErrRecordNotFound) {
			existingUser = User{Email: request.Email}
			db.Create(&existingUser)
		}

		json.NewEncoder(w).Encode(existingUser)
	}
}
