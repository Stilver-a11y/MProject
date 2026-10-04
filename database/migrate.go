package database

import (
	"context"
	"fmt"
	"log"
)

func Migrate() {
    // SQL for creating users table
    createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        avatar VARCHAR(255),
        username VARCHAR(50) UNIQUE NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT NOW()
    );`

    // SQL for creating posts table
    createPostsTable := `
    CREATE TABLE IF NOT EXISTS posts (
        id SERIAL PRIMARY KEY,
        title VARCHAR(200) NOT NULL,
        description TEXT,
        media VARCHAR(255),
        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        created_at TIMESTAMP DEFAULT NOW()
    );`

    // Create users table
    _, err := Db.Exec(context.Background(), createUsersTable)
    if err != nil {
        log.Fatal("Failed to create users table:", err)
    }

    // Create posts table
    _, err = Db.Exec(context.Background(), createPostsTable)
    if err != nil {
        log.Fatal("Failed to create posts table:", err)
    }

    fmt.Println("Tables created successfully")
}