package additionaljson

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var DefaultUnmarshaler = Unmarshaler{UnmarshalFunc: json.Unmarshal}

var UnexpectedTypeErr = errors.New("failed to unmarshal json: interface to unmarshal into must be a pointer to a struct")

type Unmarshaler struct {
	UnmarshalFunc func([]byte, interface{}) error
}

func (um Unmarshaler) Unmarshal(data []byte, i interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct { //pointer to struct
		return UnexpectedTypeErr
	}

	err := um.UnmarshalFunc(data, i)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	declaredJSONFields := listDeclaredJSONFieldNames(v.Elem())

	for i := 0; i < v.Elem().Type().NumField(); i++ {
		f := v.Elem().Type().Field(i)
		ajTag, found := f.Tag.Lookup("aj")
		if !found {
			continue
		}

		if ajTag == "all" {
			err := json.Unmarshal(data, v.Elem().Field(i).Addr().Interface())
			if err != nil {
				return fmt.Errorf("failed to unmarshal json field with tag %q: %w", ajTag, err)
			}

		} else if ajTag == "other" {
			var tmp map[string]json.RawMessage
			err := json.Unmarshal(data, &tmp)
			if err != nil {
				return fmt.Errorf("failed to unmarshal generic json field-map for tag %q: %w", ajTag, err)
			}
			//drop declared fields
			for _, f := range declaredJSONFields {
				delete(tmp, f)
			}

			//return err can be ignored cause:
			// - json.RawMessage{...}.MarshalJSON() does never respond with an error
			// - map keys (string) cannot be nil cause (syntax error at unmarshal step)
			b, _ := json.Marshal(tmp)

			err = json.Unmarshal(b, v.Elem().Field(i).Addr().Interface())
			if err != nil {
				return fmt.Errorf("failed to unmarshal json field with tag %q: %w", ajTag, err)
			}
		}
	}

	return nil
}

func listDeclaredJSONFieldNames(v reflect.Value) []string {
	var declaredJSONFields []string
	for i := 0; i < v.NumField(); i++ {
		jsonTag, found := v.Type().Field(i).Tag.Lookup("json")
		if !found {
			continue
		}

		splitTag := strings.Split(jsonTag, ",")
		if len(splitTag) == 0 || splitTag[0] == "-" {
			continue
		}

		declaredJSONFields = append(declaredJSONFields, jsonTag)
	}

	return declaredJSONFields
}
