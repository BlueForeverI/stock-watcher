package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func getAllStocks(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var stocks []Stock
		db.Find(&stocks)
		json.NewEncoder(w).Encode(stocks)
	}
}

func addStock(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(mux.Vars(r)["id"])
		var stock Stock
		json.NewDecoder(r.Body).Decode(&stock)

		var user User
		db.Find(&user, userId)
		db.Find(&stock, stock.ID)

		db.Model(&user).Association("Stocks").Append([]Stock{stock})
	}
}

func getUserStocks(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(mux.Vars(r)["id"])
		var user User
		db.Find(&user, userId)

		var stocks []Stock
		db.Model(&user).Association("Stocks").Find(&stocks)
		json.NewEncoder(w).Encode(stocks)
	}
}

func deleteStock(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := strconv.Atoi(mux.Vars(r)["id"])
		var user User
		db.Find(&user, userID)

		stockID, _ := strconv.Atoi(mux.Vars(r)["stockId"])
		var stock Stock
		db.Find(&stock, stockID)

		db.Model(&user).Association("Stocks").Delete([]Stock{stock})
	}
}
