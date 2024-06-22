package ext

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

type ResponseWriter struct {
	gin.ResponseWriter
	Buffer *bytes.Buffer
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.Buffer.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *ResponseWriter) Body() []byte {
	return w.Buffer.Bytes()
}

func (w *ResponseWriter) WriteString(s string) (int, error) {
	w.Buffer.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
