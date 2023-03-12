package database

import (
	"database/sql"
	"fmt"
	"log"

	"example.com/games/environment"
	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

func Connection() {

	user := environment.ViperEnvVariable("MYSQL_USER")
	password := environment.ViperEnvVariable("MYSQL_PASSWORD")
	host := environment.ViperEnvVariable("MYSQL_HOST")
	dbname := environment.ViperEnvVariable("DB_NAME")
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
