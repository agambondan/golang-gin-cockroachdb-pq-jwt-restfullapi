package controllers

import (
	"../controllers"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func Connect(server *controllers.Server, DBDriver, DBUrl string) {
	var err error
	server.DB, err = sql.Open(DBDriver, DBUrl)
	if err != nil {
		fmt.Printf("\nCannot connect to %s database", DBDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("\nWe are connected to the %s database", DBDriver)
	}
}
