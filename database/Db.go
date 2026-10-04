package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Db *pgxpool.Pool

func Connect() {
	connString := "postgres://postgres:root@localhost:5432/auth_db"

	var err error
	Db, err = pgxpool.New(context.Background(), connString)
	if err != nil{
		log.Fatal("Unable to connect to database:", err)
	}

	err = Db.Ping(context.Background())
	if err != nil {
		log.Fatal("Database ping failed:", err)
	}

	fmt.Println("Connected to PostgreSQL")
}