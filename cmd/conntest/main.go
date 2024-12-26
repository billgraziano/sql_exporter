package main

import (
	"database/sql"
	"log"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	connstr := "sqlserver://D40/SQL2016"
	db, err := sql.Open("sqlserver", connstr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ping success.")
	var name string
	err = db.QueryRow("SELECT @@SERVERNAME").Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("@@servername:", name)
}
