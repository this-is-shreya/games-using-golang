package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"example.com/games/environment"
	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

func Connection() {

	user := os.Getenv("MYSQL_USER")
	if user == "" {
		user = environment.ViperEnvVariable("MYSQL_USER")
	}
	password := os.Getenv("MYSQL_PASSWORD")
	if password == "" {
		password = environment.ViperEnvVariable("MYSQL_PASSWORD")
	}
	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = environment.ViperEnvVariable("MYSQL_HOST")
	}
	dbname := os.Getenv("MYSQL_DBNAME")
	if dbname == "" {
		dbname = environment.ViperEnvVariable("DB_NAME")
	}

	conn := user + ":" + password + "@tcp(" + host + ":3306)/" + dbname

	Db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatal(err)
	}
	error := Db.Ping()
	if error != nil {
		log.Fatal(error)
	} else {
		fmt.Println("database connected")
		Database = Db
	}
}
