package httputils

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type router struct {
	router *httprouter.Router
}

func (self *router) Get(path string, handler http.Handler) {
	self.router.GET(path, wrapHandler(handler))
}

func (self *router) Post(path string, handler http.Handler) {
	self.router.POST(path, wrapHandler(handler))
}

func (self *router) Put(path string, handler http.Handler) {
	self.router.PUT(path, wrapHandler(handler))
}

func (self *router) Delete(path string, handler http.Handler) {
	self.router.DELETE(path, wrapHandler(handler))
}

func (self *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	self.router.ServeHTTP(w, req)
}

func NewRouter() *router {
	return &router{httprouter.New()}
}