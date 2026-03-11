package jsonschema

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Schema ...
type Schema struct {
	Title                string             `json:"title,omitempty"`
	Type                 string             `json:"type"`
	Items                *Schema            `json:"items,omitempty"`
	Description          string             `json:"description,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	Required             []string           `json:"required,omitempty"`
	XOrder               []string           `json:"x-order,omitempty"`
	AdditionalProperties bool               `json:"additionalProperties"`
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
	title := typ.Name()
	if title == "" {
		return nil, errors.New("anonymous structures not allowed in schemas")
	}
	return &Schema{
		Title:      title,
		Type:       "object",
		Properties: properties,
		Required:   required,
		XOrder:     xorder,
	}, nil
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
		return &Schema{Type: "integer"}, nil
	case stringType:
		return &Schema{Type: "string"}, nil
	case float32Type, float64Type:
		return &Schema{Type: "number"}, nil
	case boolType:
		return &Schema{Type: "boolean"}, nil
	}
	switch {
	case typ.Kind() == reflect.Struct:
		return ForType(typ)
	case typ.Kind() == reflect.Pointer && typ.Elem().Kind() == reflect.Struct:
		sch, err := ForType(typ.Elem())
		if err != nil {
			return nil, err
		}
		return sch, nil
	case typ.Kind() == reflect.Slice:
		sch, err := getSchema(typ.Elem())
		if err != nil {
			return nil, err
		}
		return &Schema{
			Type:  "array",
			Items: sch,
		}, nil
	}
	return nil, fmt.Errorf("unknown type for schema: %s", typ)
}
