package encoding

import (
	"net/http"

	"github.com/go-kratos/kratos/v2/encoding/json"
)

type JSON struct {
	Data any
}

var jsonContentType = []string{"application/json; charset=utf-8"}

func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func (r JSON) WriteJSON(w http.ResponseWriter, obj any) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := GetCodec(json.Name).Marshal(obj)

	if err != nil {
		return err
	}

	_, err = w.Write(jsonBytes)
	return err
}

func (r JSON) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)

	if err = r.WriteJSON(w, r.Data); err != nil {
		panic(err)
	}

	return
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
