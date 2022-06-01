package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string `gorm:"unique"`
	Processed bool
}
type Response struct {
	Status string `json:"status"`
}

func main() {
	db, e := gorm.Open(postgres.Open("host=127.0.0.1 port=5432 user=accounts password=accounts dbname=accounts sslmode=disable"), &gorm.Config{})
	if e != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.HandleFunc("/users", func(rw http.ResponseWriter, r *http.Request) {

		rw.Header().Add("Content-Type", "application/json")
		email := r.FormValue("email")
		err := db.Create(&User{Email: email}).Error
		if err != nil {
			json.NewEncoder(rw).Encode(Response{Status: err.Error()})
			return
		}
		json.NewEncoder(rw).Encode(Response{Status: "Ok"})

	}).Methods("POST")

	fmt.Println("started")
	server := http.Server{Addr: "127.0.0.1:1000", Handler: r}
	server.ListenAndServe()
}
