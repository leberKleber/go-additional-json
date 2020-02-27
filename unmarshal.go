package go_additional_json

import (
	"encoding/json"
)

var DefaultUnmarshaler = Unmarshaler{UnmarshalFunc: json.Unmarshal}

type Unmarshaler struct {
	UnmarshalFunc func([]byte, interface{}) error
}

func (um Unmarshaler) Unmarshal(data []byte, i interface{}) error {
	return nil
}
