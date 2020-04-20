package controllers

import (
	"../middlewares"
)

func (server *Server) initializeRoutes() {

	// Home Route
	server.Router.GET("/", server.Home)

	// Login Route
	server.Router.POST("/login", server.Login)

	// Message Broker
	server.Router.POST("/booking", server.SendMessage)
	server.Router.GET("/notify", server.GetConsumeMessageBroker)

	// User Route
	server.Router.POST("/user", server.CreateUser)
	server.Router.GET("/users", server.GetAllUser)
	server.Router.GET("/user/:id", server.GetUserById)
	server.Router.PUT("/user/:id", middlewares.TokenAuthMiddleware(), server.UpdateUserById)
	server.Router.DELETE("/user/:id", middlewares.TokenAuthMiddleware(), server.DeleteUserById)
	server.Router.POST("/user/:id/upload-image", middlewares.TokenAuthMiddleware(), server.UploadFile)

	// Post Route
	server.Router.POST("/post", middlewares.TokenAuthMiddleware(), server.CreatePost)
	server.Router.GET("/posts", server.GetAllPost)
	server.Router.GET("/post/:id", server.GetPostById)
	server.Router.PUT("/post/:id", middlewares.TokenAuthMiddleware(), server.UpdatePostById)
	server.Router.DELETE("/post/:id", middlewares.TokenAuthMiddleware(), server.DeletePostById)

	// Role Route
	server.Router.POST("/role", middlewares.TokenAuthMiddleware(), server.CreateRole)
	server.Router.GET("/roles", server.GetAllRole)
	server.Router.GET("/role/:id", server.GetRoleById)
	server.Router.PUT("/role/:id", middlewares.TokenAuthMiddleware(), server.UpdateRoleById)
	server.Router.DELETE("/role/:id", middlewares.TokenAuthMiddleware(), server.DeleteRoleById)

	// Export || Import Route
	server.Router.GET("/export-user", middlewares.TokenAuthMiddleware(), server.ExportUserToExcel)
	server.Router.POST("/import-user", middlewares.TokenAuthMiddleware(), server.ImportExcelToUser)

	// SMTP Route
	server.Router.POST("/send-email", middlewares.TokenAuthMiddleware(), server.SendEmail)

}
