package http_server

import (
	"github.com/gin-gonic/gin"
)

type APIHandlerRouter interface {
	AddRoutes(r *gin.RouterGroup)
	GetVersion() string
}

type Router struct {
	router *gin.Engine
}

func NewRouter() Router {
	return Router{router: gin.Default()}
}

func (r *Router) WithHandle(method string, path string, handler gin.HandlerFunc) *Router {
	r.router.Handle(method, path, handler)
	return r
}

func (r *Router) WithHandleGET(path string, handler gin.HandlerFunc) *Router {
	r.router.GET(path, handler)
	return r
}

func (r *Router) GetRouter() *gin.Engine {
	return r.router
}
