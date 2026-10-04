package main

import (
	"fmt"
	"log"
	database "mproject/database"
	"mproject/handlers"
	"net/http"
)

func main() {
	database.Connect()
	database.Migrate()
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintln(w, "Server is running locally")
	})
	
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)

	log.Println("server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}