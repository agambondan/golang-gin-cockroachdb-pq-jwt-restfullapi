package controllers

import (
	"../middlewares"
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
	var err error
	server.DB, err = sql.Open(DBDriver, DBUrl)
	if err != nil {
		fmt.Printf("\nCannot connect to %s database", DBDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("\nWe are connected to the %s database", DBDriver)
	}
	server.Router = gin.Default()
	server.Router.Use(middlewares.CORSMiddleware())
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("\nListening to localhost" + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
