package main

import (
	"log"
	"net/http"
	"os"

	"./src"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	myRouter := mux.NewRouter().StrictSlash(true).PathPrefix("/api").Subrouter()
	myRouter.Use(commonMiddleware)

	db, _ := gorm.Open(mysql.Open(os.Getenv("MYSQL_CONN")), &gorm.Config{})
	myRouter.HandleFunc("/login", src.Login(db)).Methods("POST", "OPTIONS")

	myRouter.HandleFunc("/stocks", src.GetAllStocks(db))
	myRouter.HandleFunc("/stocks/trending", src.GetTrendingStocks)

	myRouter.HandleFunc("/users/{id}/stocks", src.GetUserStocks(db)).Methods("GET", "OPTIONS")
	myRouter.HandleFunc("/users/{id}/stocks", src.AddStock(db)).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/users/{id}/stocks/{stockId}", src.DeleteStock(db)).Methods("DELETE", "OPTIONS")
	myRouter.HandleFunc("/users/{id}/watchlist", src.GetWatchlist(db))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), myRouter))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Type")
		next.ServeHTTP(w, r)
	})
}
