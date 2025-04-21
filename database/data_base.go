package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func DbInit() *sqlx.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Не удалось загрузить .env")
	}
	Str := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	DB, err = sqlx.Connect("postgres", Str)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных:%v", err)
	}

	fmt.Println("Успешное подключение к базе данных")
	return DB
}
