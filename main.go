package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mproject/database"
	"mproject/handlers"
	"net/http"

	"mproject/middleware"
)

func main() {
	database.Connect()
	database.Migrate()
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintln(w, "Server is running locally")
	})
	
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/me", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.UserIDKey).(uint)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id": userID,
			"message": "You are authenticated",
		})
	}))
	http.HandleFunc("/posts", middleware.Auth(func(w http.ResponseWriter, r *http.Request){
		switch r.Method {
		case http.MethodPost:
			handlers.CreatePost(w, r)
		case http.MethodGet:
			handlers.GetPost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/posts/", middleware.Auth(func(w http.ResponseWriter, r *http.Request){
		switch r.Method {
		case http.MethodPut:
			handlers.UpdatePost(w, r)
		case http.MethodDelete:
			handlers.DeletePost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	log.Println("server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}