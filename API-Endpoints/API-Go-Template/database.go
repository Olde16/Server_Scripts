package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDatabase() {
	dsn := "olympics_db_user:EFI24B2026@unix(/var/run/mysqld/mysqld.sock)/olympics?parseTime=true"

	var err error

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Could not open database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	log.Println("Connected to MySQL")
}
