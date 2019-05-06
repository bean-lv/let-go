package context

import "net/http"

type Context interface {
	Request() *http.Request
	Response() http.ResponseWriter
	Reset(rw http.ResponseWriter, req *http.Request)
}

type myContext struct {
	request        *http.Request
	responseWriter http.ResponseWriter
}

func New() Context {
	return &myContext{}
}

func (ctx *myContext) Request() *http.Request {
	return ctx.request
}

func (ctx *myContext) Response() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *myContext) Reset(rw http.ResponseWriter, req *http.Request) {
	ctx.responseWriter = rw
	ctx.request = req
}
