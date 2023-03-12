package main

import (
	"net/http"
	"os"

	"example.com/games/controllers"
	"example.com/games/database"
	"example.com/games/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	database.Connection()
	r := mux.NewRouter()
	r.HandleFunc("/api/signup", controllers.Signup).Methods("POST")
	r.HandleFunc("/api/login", controllers.Login).Methods("POST")
	r.HandleFunc("/test", middleware.VerifyToken(controllers.Test)).Methods("POST")
	PORT := os.Getenv("PORT")
	http.ListenAndServe(PORT, r)
}
