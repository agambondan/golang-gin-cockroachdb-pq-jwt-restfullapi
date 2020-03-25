package controllers

//
import (
	"../auth"
	"../models"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (server *Server) CreatePost(c *gin.Context) {
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized"})
		return
	}
	post := models.Post{}
	err = c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	uId, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized, uid : " + uId + err.Error()})
		return
	}
	if uId != post.AuthorID.String() {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized, Author Id != uid " + uId + " != " + post.AuthorID.String()})
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
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized, " + err.Error()})
		return
	}
	post := models.Post{}
	// Check if the post exist
	err = server.DB.QueryRow("SELECT id, created_at, updated_at, title, content, author_id FROM post WHERE id=$1", pId).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.Title, &post.Content, &post.AuthorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	//CHeck if the auth token is valid and get the user id from it
	uId, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	// If a user attempt to update a post not belonging to him
	if uId != post.AuthorID.String() {
		c.JSON(http.StatusUnauthorized, gin.H{"message": errors.New("your user id not the same like post author id " + uId + " " + post.AuthorID.String()).Error()})
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
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized, " + err.Error()})
		return
	}
	id := c.Params.ByName("id")
	pId := uuid.MustParse(id)
	post := models.Post{}
	user := models.User{}
	err = server.DB.QueryRow("SELECT author_id FROM post WHERE id=$1", pId).Scan(&post.AuthorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}
	err = server.DB.QueryRow("SELECT id FROM users WHERE id=$1", post.AuthorID).Scan(&user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}
	userId, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	if userId != user.ID.String() {
		c.JSON(http.StatusUnauthorized, gin.H{"message": errors.New("your user id not the same like post author id " + user.ID.String() + " " + userId)})
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
