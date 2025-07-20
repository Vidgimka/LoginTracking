package myhttp

import (
	"github.com/Vidgimka/LoginTracking.git/myhttp/handlers"
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
	router.GET("/UsersOnline2/:login/date/:datetime", r.h.GetUserByLoginAndSessionId)
	router.GET("/UsersOnline2", r.h.GetAllUsers)
	router.GET("/UsersOnline2/:login", r.h.GetUserByLogin)
	router.GET("/UsersOnline2/:login/:session_id", r.h.GetUserByLoginAndSessionId)

	return router
}
