package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Server struct {
	DB     *sql.DB
	Router *gin.Engine
}

func (server *Server) Initialize(DBDriver, DBUrl string) {
	if DBDriver == "postgres" {
		var err error
		server.DB, err = sql.Open(DBDriver, DBUrl)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", DBDriver)
			fmt.Println()
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", DBDriver)
			fmt.Println()
		}
	}
	server.Router = gin.Default()
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to localhost" + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
