package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func Connection() *sql.DB {
	db, err := sql.Open("postgres", viper.GetString("POSTGRES_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	log.Println("Successfully Connected to DB")

	return db
}
