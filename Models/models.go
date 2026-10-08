package Models

import "time"

type User struct {
	ID	uint	`json:"id"`
	Avatar	string	`json:"avatar"`
	Username	string	`json:"username"`
	Email	string	`json:"email"`
	PasswordHash	string	`json:"-"`
	CreatedAt	time.Time	`json:"created_at"`
}

type Post struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Media       string    `json:"media"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}