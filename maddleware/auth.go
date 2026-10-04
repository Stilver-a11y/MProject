/*package maddleware

import

var JwtSecret = []byte("change-me-to-a-long-random-secret")

type ContextKey string

const UserIDKey ContextKey = "user_id"

func Auth(next http.HandleFunc) http.HandleFunc{
	return func(w http.ResponseWriter, r *http.Request)
}