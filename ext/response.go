package ext

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
)

type Context struct {
	*gin.Context
}

type Response interface {
	Error(c *Context, err error)
	Success(c *Context, data interface{})
}

// NewContext 包装一下
func NewContext(ctx *gin.Context) *Context {
	return &Context{ctx}
}

type Resp struct {
}

func (r Resp) Error(c *Context, err error) {
	e := errors.FromError(err)
	c.JSON(200, gin.H{
		"code": e.Reason,
		"msg":  e.Message,
	})
	return
}

func (r Resp) Success(c *Context, data interface{}) {
	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}
