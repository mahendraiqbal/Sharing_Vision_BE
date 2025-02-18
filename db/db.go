package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB menginisialisasi koneksi database
func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Gagal melakukan ping ke database:", err)
	}

	log.Println("Berhasil terhubung ke database!")
}

// CloseDB menutup koneksi database
func CloseDB() {
	DB.Close()
}
