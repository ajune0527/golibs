package ext

import (
	"github.com/gin-gonic/gin"
	errors "github.com/go-kratos/kratos/v2/errors"
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
	if e.Code != 0 {
		c.JSON(200, gin.H{
			"code": e.Code,
			"msg":  e.Message,
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 400,
		"msg":  err.Error(),
	})
}

func (r Resp) Success(c *Context, data interface{}) {
	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}
