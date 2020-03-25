package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (server *Server) Home(c *gin.Context) {
	//responses.JSON(w, http.StatusOK, "Welcome To This API")
	c.JSON(http.StatusOK, "Welcome To This API")
}

