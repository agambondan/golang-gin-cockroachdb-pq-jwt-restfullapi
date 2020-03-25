package middlewares

import (
	"../auth"
	"../responses"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("UnAuthorized"))
			return
		}
		next(w, r)
	}
}

func SetMiddlewareAuthentication1(c *gin.Context) *gin.Context {
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
	}
	return c
}

func SetMiddlewareAuthentication2(c *gin.Context) *gin.Context {
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
	}
	return c
}
