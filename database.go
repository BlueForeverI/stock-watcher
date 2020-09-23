package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db, _ = gorm.Open(mysql.Open(os.Getenv("MYSQL_CONN")), &gorm.Config{})

func allProducts(w http.ResponseWriter, r *http.Request) {
	var products []Product
	db.Find(&products)

	json.NewEncoder(w).Encode(products)
}

func productByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var product Product
	db.First(&product, "code = ?", id)

	json.NewEncoder(w).Encode(product)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	json.NewDecoder(r.Body).Decode(&product)

	db.Create(&product)

	json.NewEncoder(w).Encode(product)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(commonMiddleware)
	myRouter.HandleFunc("/products", allProducts).Methods("GET")
	myRouter.HandleFunc("/products/{id}", productByID)
	myRouter.HandleFunc("/products", createProduct).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), myRouter))
}
