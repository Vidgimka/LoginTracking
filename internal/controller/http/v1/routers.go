package v1

import (
	"github.com/gin-gonic/gin"
)

type handlers interface {
	GetPointsByLogin(c *gin.Context)
	GetLinesByLogin(c *gin.Context)
	GetPointsByDate(c *gin.Context)
	GetLinesByDate(c *gin.Context)
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
	router.GET("/loginytracking/v1/logins/:login/date", r.h.GetLinesByDate) // http://localhost:8080/loginytracking/v1/logins/tsb645/date?start=....&end=....&visual=...
	//2025-06-29T17:33:54.253593+03:00
	return router
}
