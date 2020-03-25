package controllers

import (
	//"../auth"
)

func (server *Server) initializeRoutes() {

	// Home Route
	server.Router.GET("/", server.Home)

	// Login Route
	server.Router.POST("/login", server.Login)

	// User Route
	server.Router.POST("/user", server.CreateUser)
	server.Router.GET("/users", server.GetAllUser)
	server.Router.GET("/user/:id", server.GetUserById)
	server.Router.PUT("/user/:id", server.UpdateUserById)
	server.Router.DELETE("/user/:id", server.DeleteUserById)

	// Post Router
	server.Router.POST("/post", server.CreatePost)
	server.Router.GET("/posts", server.GetAllPost)
	server.Router.GET("/post/:id", server.GetPostById)
	server.Router.PUT("/post/:id", server.UpdatePostById)
	server.Router.DELETE("/post/:id", server.DeletePostById)
}
