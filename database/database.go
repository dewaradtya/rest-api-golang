package database

import (
    "database/sql"
    "log"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Inisialisasi koneksi ke database
func Connect() {
	dsn := "root@tcp(localhost:3306)/gocrud_app"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Error connecting to database:", err)
    }

    // Pastikan koneksi berhasil
    if err := db.Ping(); err != nil {
        log.Fatal("Cannot reach database:", err)
    }

    log.Println("Database connected successfully")
    DB = db
}
