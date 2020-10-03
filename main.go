package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	myRouter := mux.NewRouter().StrictSlash(true).PathPrefix("/api").Subrouter()
	myRouter.Use(commonMiddleware)

	db, _ := gorm.Open(mysql.Open(os.Getenv("MYSQL_CONN")), &gorm.Config{})
	myRouter.HandleFunc("/login", login(db)).Methods("POST")

	myRouter.HandleFunc("/stocks", getAllStocks(db))
	myRouter.HandleFunc("/stocks/trending", getTrendingStocks)

	myRouter.HandleFunc("/users/{id}/stocks", getUserStocks(db)).Methods("GET")
	myRouter.HandleFunc("/users/{id}/stocks", addStock(db)).Methods("POST")
	myRouter.HandleFunc("/users/{id}/stocks/{stockId}", deleteStock(db)).Methods("DELETE")
	myRouter.HandleFunc("/users/{id}/watchlist", getWatchlist(db))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), myRouter))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
