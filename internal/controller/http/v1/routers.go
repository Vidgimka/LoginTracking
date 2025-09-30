package v1

import (
	"time"

	"github.com/gin-gonic/gin"
)

type handlers interface {
	GetPointsByLogin(c *gin.Context, login string)
	GetLinesByLogin(c *gin.Context, login string)
	GetPointsByDate(c *gin.Context, login string, start time.Time, end time.Time)
	GetLinesByDate(c *gin.Context, login string, start time.Time, end time.Time)
}

type router struct {
	h handlers
}

func NewRouter(h handlers) *router {
	return &router{
		h: h,
	}
}

func (r *router) SetRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/UsersOnline2/points/:login", r.h.GetPointsByLogin())     // http://localhost:8080/UsersOnline2/points/tsb645
	router.GET("/UsersOnline2/lines/:login", r.h.GetLinesByLogin())       // http://localhost:8080/UsersOnline2/lines/tsb645
	router.GET("/UsersOnline2/points/:login/date", r.h.GetPointsByDate()) // http://localhost:8080/UsersOnline2/points/tsb645|date?start=....&end=....
	router.GET("/UsersOnline2/lines/:login/date", r.h.GetLinesByDate())   // http://localhost:8080/UsersOnline2/lines/tsb645|date?start=....&end=....
	//2025-06-29T17:33:54.253593+03:00
	return router
}
