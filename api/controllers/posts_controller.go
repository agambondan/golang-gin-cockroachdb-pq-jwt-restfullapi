package controllers

//
import (
	"../models"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (server *Server) CreatePost(c *gin.Context) {
	userToken := extractToken(c)
	post := models.Post{}
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	if userToken.ID != post.AuthorID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized, Author Id != uid " + userToken.ID.String() + " != " + post.AuthorID.String()})
		return
	}
	err = post.Validate()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	savePost, err := post.SavePost(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": savePost})
}

func (server *Server) GetAllPost(c *gin.Context) {
	post := models.Post{}
	var posts []models.Post
	findAllPost, err := post.FindAllPost(server.DB)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}
	bytes, err := json.Marshal(findAllPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	err = json.Unmarshal([]byte(bytes), &posts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message1": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": findAllPost})
}

func (server *Server) GetPostById(c *gin.Context) {
	id := c.Params.ByName("id")
	uId := uuid.MustParse(id)
	post := models.Post{}
	postById, err := post.FindPostByID(server.DB, uId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"message": `Data By Id ` + uId.String() + ` Is Found`,
		"data":    postById,
	})
}

func (server *Server) UpdatePostById(c *gin.Context) {
	id := c.Params.ByName("id")
	pId := uuid.MustParse(id)
	userToken := extractToken(c)
	post := models.Post{}
	_, err := post.FindPostByID(server.DB, pId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	if userToken.ID != post.AuthorID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": errors.New("your user id not the same like post author id " + userToken.ID.String() + " " + post.AuthorID.String()).Error()})
		return
	}
	postById, err := post.UpdatePostById(server.DB, pId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": postById})
}

func (server *Server) DeletePostById(c *gin.Context) {
	userToken := extractToken(c)
	id := c.Params.ByName("id")
	pId := uuid.MustParse(id)
	post := models.Post{}
	user := models.User{}
	_, err := post.FindPostByID(server.DB, post.AuthorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}
	_, err = user.FindUserById(server.DB, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}
	if userToken.ID != user.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": errors.New("your user id not the same like post author id " + user.ID.String() + " " + userToken.ID.String())})
		return
	}
	deleteUserById, err := post.SoftDeletePostById(server.DB, pId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"message": `Delete Data By Id ` + id + ` Successfully`,
		"data":    deleteUserById,
	})
}
