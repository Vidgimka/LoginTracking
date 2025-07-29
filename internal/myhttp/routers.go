package myhttp

import (
	"github.com/Vidgimka/LoginTracking.git/internal/myhttp/handlers"
	"github.com/gin-gonic/gin"
)

type router struct {
	h handlers.HandlersInterface
}

func NewRouter(h handlers.HandlersInterface) *router {
	return &router{
		h: h,
	}
}

func (r *router) SetRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/UsersOnline2/:login/date/:datetime", r.h.GetUserByLoginAnDatetime) // http://localhost:8080/UsersOnline2/tsb645/date/2025-06-29T17:33:54.253593+03:00
	router.GET("/UsersOnline2", r.h.GetAllUsers)                                    // http://localhost:8080/UsersOnline2
	router.GET("/UsersOnline2/:login", r.h.GetUserByLogin)                          // http://localhost:8080/UsersOnline2/tsb645
	router.GET("/UsersOnline2/:login/date", r.h.GetlineCollection)                  // http://localhost:8080/UsersOnline2/tsb645|date?start=....&end=....
	router.GET("/UsersOnline2/:login/:session_id", r.h.GetUserByLoginAndSessionId)  // http://localhost:8080/UsersOnline2/tsb645/5984

	return router
}
