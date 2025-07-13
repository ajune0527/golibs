package ext

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/ajune0527/golibs/ext/encoding"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-kratos/kratos/v2/encoding/json"
	kerr "github.com/go-kratos/kratos/v2/errors"

	"github.com/google/uuid"
)

type Context struct {
	*gin.Context
	Body      *bytes.Buffer
	StartTime time.Time
}

type Response interface {
	Error(c *Context, err error)
	Success(c *Context, data interface{})
}

// NewContext 包装一下
func NewContext(ctx *gin.Context) *Context {
	ctx.Writer = &ResponseWriter{
		ResponseWriter: ctx.Writer,
		Buffer:         &bytes.Buffer{},
	}
	return &Context{
		Context:   ctx,
		StartTime: time.Now(),
	}
}

var (
	EnableMsgShortCode = false // 是否开启错误消息附加短码
)

type Resp struct {
}

func (r Resp) Error(c *Context, err error) {

	err2 := kerr.FromError(err)
	c.Status(int(err2.GetCode()))

	code := err2.GetReason()
	if code == "" {
		code = "Unknown"
	}

	msg := err2.GetMessage()

	if EnableMsgShortCode {
		msg = fmt.Sprintf("%s[%d]", err2.GetMessage(), time.Now().UnixMilli()%100000)
	}

	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
	})
	return
}

func (r Resp) Success(c *Context, data interface{}) {
	c.Render(200, encoding.JSON{Data: data})
}

func (c *Context) ShouldBindJSON(obj any) error {
	if c.Request == nil || c.Request.Body == nil {
		return fmt.Errorf("invalid request")
	}
	// 如果是post请求， 但是请求的Content-Type 类型是 application/x-www-form-urlencoded
	if c.Request.Method == "POST" && c.Request.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		return decodeUrlencodedJSON(c.Request.Body, obj)
	}

	return decodeJSON(c.Request.Body, obj)
}

func (c *Context) GetBody() ([]byte, error) {
	body, err := c.Context.GetRawData()
	if err != nil {
		return nil, err
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}

func (c *Context) ResponseWriter() *ResponseWriter {
	if resp, ok := c.Writer.(*ResponseWriter); ok {
		return resp
	}
	return nil
}

func (c *Context) ResponseBody() []byte {
	if resp := c.ResponseWriter(); resp != nil {
		return resp.Body()
	}
	return nil
}

func (c *Context) RequestId() string {
	requestId := c.Request.Header.Get("Request-ID")
	if requestId == "" {
		requestId = uuid.New().String()
		c.Request.Header.Set("Request-ID", requestId)
	}

	return requestId
}

func decodeJSON(r io.ReadCloser, obj any) error {
	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if err = encoding.GetCodec(json.Name).Unmarshal(body, obj); err != nil {
		return err
	}

	return validate(obj)
}

// decodeUrlencodedJSON 将 URL-encoded 格式数据（如 a=1&b=2）解析到目标结构体
func decodeUrlencodedJSON(r io.ReadCloser, obj any) error {
	// 1. 读取请求体（URL-encoded 格式，如 "name=foo&age=20"）
	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	defer r.Close()
	m := make(map[string]any)
	if d, err := url.ParseQuery(string(body)); err == nil {
		for k := range d {
			m[k] = d.Get(k)
		}
	}

	return validate(obj)
}

func validate(obj any) error {
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}
