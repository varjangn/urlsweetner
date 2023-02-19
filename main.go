package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/varjangn/urlsweetner/db"
	"github.com/varjangn/urlsweetner/handlers"
	"github.com/varjangn/urlsweetner/middlewares"

	"github.com/joho/godotenv"
)

func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func initDB() {
	var err error
	db.DbRepo, err = db.NewSQLiteRepository(fmt.Sprintf("%s.db", os.Getenv("DB_NAME")))
	if err != nil {
		panic(err)
	}
	if err := db.DbRepo.Migrate(); err != nil {
		panic(err)
	}
}

func main() {
	loadEnvVars()
	initDB()

	router := mux.NewRouter()

	router.HandleFunc("/register", middlewares.Chain(
		handlers.Register, middlewares.Method("POST"), middlewares.Logging()))

	router.HandleFunc("/login", middlewares.Chain(
		handlers.Login, middlewares.Method("POST"), middlewares.Logging()))

	router.HandleFunc("/hello", middlewares.Chain(
		handlers.Hello, middlewares.RequireAuth()))

	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", router)
}
