package jsonschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Type ...
type Type []string

// MarshalJSON ...
func (t Type) MarshalJSON() ([]byte, error) {
	if len(t) == 1 {
		return json.Marshal(t[0])
	}
	return json.Marshal([]string(t))
}

// UnmarshalJSON ...
func (t *Type) UnmarshalJSON(b []byte) error {
	if b[0] == '[' {
		return json.Unmarshal(b, (*[]string)(t))
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*t = []string{s}
	return nil
}

var (
	_ json.Marshaler   = Type{}
	_ json.Unmarshaler = (*Type)(nil)
)

// Schema ...
type Schema struct {
	Title       string             `json:"title,omitempty"`
	Type        Type               `json:"type"`
	Items       *Schema            `json:"items,omitempty"`
	Description string             `json:"description,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
	Required    []string           `json:"required,omitempty"`
	XOrder      []string           `json:"x-order,omitempty"`
}

// For ...
func For[T any]() (*Schema, error) {
	return ForType(reflect.TypeFor[T]())
}

// ForType ...
func ForType(typ reflect.Type) (*Schema, error) {
	properties := make(map[string]*Schema, typ.NumField())
	required := make([]string, 0, typ.NumField())
	xorder := make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		js := field.Tag.Get("json")
		if js == "-" {
			continue
		}
		schema, err := getSchema(field.Type)
		if err != nil {
			return nil, err
		}
		schema.Description = field.Tag.Get("jsonschema")
		name := field.Name
		comps := strings.Split(js, ",")
		if comps[0] != "" {
			name = comps[0]
		}
		properties[name] = schema
		if len(comps) == 1 || comps[1] != "omitempty" && comps[1] != "omitzero" {
			required = append(required, name)
		}
		xorder = append(xorder, name)
	}
	return &Schema{Title: typ.Name(), Type: []string{"object"}, Properties: properties, Required: required, XOrder: xorder}, nil
}

var (
	intType     = reflect.TypeFor[int]()
	stringType  = reflect.TypeFor[string]()
	float32Type = reflect.TypeFor[float32]()
	float64Type = reflect.TypeFor[float64]()
	boolType    = reflect.TypeFor[bool]()
)

func getSchema(typ reflect.Type) (*Schema, error) {
	switch typ {
	case intType:
		return &Schema{Type: []string{"integer"}}, nil
	case stringType:
		return &Schema{Type: []string{"string"}}, nil
	case float32Type, float64Type:
		return &Schema{Type: []string{"number"}}, nil
	case boolType:
		return &Schema{Type: []string{"boolean"}}, nil
	}
	switch {
	case typ.Kind() == reflect.Struct:
		return ForType(typ)
	case typ.Kind() == reflect.Pointer && typ.Elem().Kind() == reflect.Struct:
		sch, err := ForType(typ.Elem())
		if err != nil {
			return nil, err
		}
		sch.Type = append(sch.Type, "null")
		return sch, nil
	case typ.Kind() == reflect.Slice:
		return nil, errors.ErrUnsupported
	}
	return nil, fmt.Errorf("unknown type for schema: %s", typ)
}
