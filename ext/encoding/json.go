package encoding

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const Name = "json"

var (
	MarshalOptions = protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}

	UnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

type codec struct{}

func (codec) Marshal(obj any) ([]byte, error) {
	var (
		jsonBytes []byte
		err       error
	)

	switch m := obj.(type) {
	case json.Marshaler:
		jsonBytes, err = m.MarshalJSON()
	case proto.Message:
		jsonBytes, err = MarshalOptions.Marshal(m)
	default:
		jsonBytes, err = json.Marshal(m)
	}
	if err != nil {
		return nil, err
	}
	return jsonBytes, err
}

func (codec) Unmarshal(data []byte, obj any) error {
	switch m := obj.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	case proto.Message:
		return UnmarshalOptions.Unmarshal(data, m)
	default:
		return json.Unmarshal(data, obj)
	}
}

func (codec) Name() string {
	return Name
}

func init() {
	RegisterCodec(codec{})
}
