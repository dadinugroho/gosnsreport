package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
)

func connectToDB() (*sql.DB, error) {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbSocket := os.Getenv("DB_SOCKET")

    dsn := fmt.Sprintf("%s:%s@unix(%s)/%s", dbUser, dbPassword, dbSocket, dbName)

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %v", err)
    }

    return db, nil
}
