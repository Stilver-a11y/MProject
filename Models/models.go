package Models

import "time"

type User struct {
	ID	uint	`json:"id"`
	Avatar	string	`json:"avatar"`
	Username	string	`json:"username"`
	Email	string	`json:"email"`
	PasswordHash	string	`json:"-"`
	СreatedAt	time.Time	`json:"created_at"`
}