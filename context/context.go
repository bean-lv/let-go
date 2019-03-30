package context

import "net/http"

type Context interface {
	Request() *http.Request
	SetRequest(*http.Request)
	Response() http.ResponseWriter
	SetResponse(http.ResponseWriter)
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

func (ctx *myContext) SetRequest(request *http.Request) {
	ctx.request = request
}

func (ctx *myContext) Response() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *myContext) SetResponse(response http.ResponseWriter) {
	ctx.responseWriter = response
}
