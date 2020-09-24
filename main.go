package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email string `json:"email"`
}

var db, _ = gorm.Open(mysql.Open(os.Getenv("MYSQL_CONN")), &gorm.Config{})

func getAllStocks(w http.ResponseWriter, r *http.Request) {
	var stocks []Stock
	db.Find(&stocks)
	json.NewEncoder(w).Encode(stocks)
}

func login(w http.ResponseWriter, r *http.Request) {
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

func addStock(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.Atoi(mux.Vars(r)["id"])
	var stock Stock
	json.NewDecoder(r.Body).Decode(&stock)

	var user User
	db.Find(&user, userId)
	db.Find(&stock, stock.ID)

	db.Model(&user).Association("Stocks").Append([]Stock{stock})
}

func getUserStocks(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.Atoi(mux.Vars(r)["id"])
	var user User
	db.Find(&user, userId)

	var stocks []Stock
	db.Model(&user).Association("Stocks").Find(&stocks)
	json.NewEncoder(w).Encode(stocks)
}

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(commonMiddleware)

	myRouter.HandleFunc("/stocks", getAllStocks)
	myRouter.HandleFunc("/login", login).Methods("POST")
	myRouter.HandleFunc("/users/{id}/stocks", addStock).Methods("POST")
	myRouter.HandleFunc("/users/{id}/stocks", getUserStocks).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), myRouter))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
