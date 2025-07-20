package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/Vidgimka/LoginTracking.git/repository"
	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Login            string    `json:"login"`
	Session_id       int       `json:"session_id"`
	Lat              float64   `json:"lat"`
	Lon              float64   `json:"lon"`
	Station_distance float64   `json:"station_distance"`
	CreatedAt        time.Time `json:"сreated_at"`
}

type handlers struct {
	repo repository.PostgresGormRepoInterfase
}

type HandlersInterface interface {
	GetAllUsers(c *gin.Context)
	GetUserByLogin(c *gin.Context)
	GetUserByLoginAndSessionId(c *gin.Context)
	GetUserByLoginAnDatetime(c *gin.Context)
}

func NewHandlers(db repository.PostgresGormRepoInterfase) HandlersInterface {
	return &handlers{
		repo: db,
	}
}

func (repo *handlers) GetUserByLoginAnDatetime(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	login := c.Param("login")
	if login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login cannot be empty"})
		return
	}
	CreatedAt := c.Param("datetime")
	if CreatedAt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SessionId cannot be empty"})
		return
	}
	if err := repo.repo.GetByDatetime(c.Request.Context(), c.Writer, login, CreatedAt); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"db error": "login and datetime data not recirved"})
	}
	c.Writer.Flush()
}

// http://localhost:8080/UsersOnline2/nje232/date/2025-06-22T21:02:30.896313+03:00

func (repo *handlers) GetAllUsers(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	if err := repo.repo.GetByAllUser(c.Request.Context(), c.Writer); err != nil {
		log.Printf("get all users:%v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"db error": "users not received"})
		return
	}
	c.Writer.Flush()
}

// http://localhost:8080/UsersOnline2

func (repo *handlers) GetUserByLogin(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	login := c.Param("login")
	if login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login cannot be empty"})
		return
	}
	if err := repo.repo.GetByLogin(c.Request.Context(), c.Writer, login); err != nil {
		log.Printf("get login:%v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"db error": "login data not recirved"})
		return
	}
	c.Writer.Flush()
}

// http://localhost:8080/UsersOnline2/aza235

func (repo *handlers) GetUserByLoginAndSessionId(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	login := c.Param("login")
	if login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login cannot be empty"})
		return
	}
	Session_id := c.Param("session_id")
	if Session_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id cannot be empty"})
		return
	}
	if err := repo.repo.GetBySessionId(c.Request.Context(), c.Writer, login, Session_id); err != nil {
		log.Printf("get by session_id: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"db error": "login and session_id data not recirved"})
	}
	c.Writer.Flush()
}

// http://localhost:8080/UsersOnline2/nje232/4968
