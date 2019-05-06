package controller

import (
	"encoding/json"
	"letgo/context"
)

type Controller interface {
	Init(ctx context.Context)
	ServeJSON(data interface{})
	Query(key string) string
}

type myController struct {
	ctx context.Context
}

func new() Controller {
	return &myController{}
}

func (c *myController) Init(ctx context.Context) {
	c.ctx = ctx
}

func (c *myController) ServeJSON(data interface{}) {
	json.NewEncoder(c.ctx.Response()).Encode(data)
}

func (c *myController) Query(key string) string {
	req := c.ctx.Request()
	if req.Form == nil {
		req.ParseForm()
	}
	return req.Form.Get(key)
}
