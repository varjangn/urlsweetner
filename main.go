package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/varjangn/urlsweetner/db"
	"github.com/varjangn/urlsweetner/handlers"
	"github.com/varjangn/urlsweetner/middlewares"
)

func main() {
	var err error
	db.DbRepo, err = db.NewSQLiteRepository("urlsweetner.db")
	if err != nil {
		panic(err)
	}
	if err := db.DbRepo.Migrate(); err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/register", middlewares.Chain(
		handlers.Register, middlewares.Method("POST"), middlewares.Logging()))

	router.HandleFunc("/login", middlewares.Chain(
		handlers.Login, middlewares.Method("POST"), middlewares.Logging()))

	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", router)
}
