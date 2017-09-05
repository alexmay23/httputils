package httputils

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Router struct {
	router *httprouter.Router
}

func (self *Router) Get(path string, handler http.Handler) {
	self.router.GET(path, wrapHandler(handler))
}

func (self *Router) Post(path string, handler http.Handler) {
	self.router.POST(path, wrapHandler(handler))
}

func (self *Router) Put(path string, handler http.Handler) {
	self.router.PUT(path, wrapHandler(handler))
}

func (self *Router) Delete(path string, handler http.Handler) {
	self.router.DELETE(path, wrapHandler(handler))
}

func (self *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	self.router.ServeHTTP(w, req)
}

func NewRouter() *Router {
	return &Router{httprouter.New()}
}