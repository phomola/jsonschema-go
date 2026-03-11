package jsonschema

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/stretchr/testify/require"
)

func TestSchemaInference(t *testing.T) {
	req := require.New(t)

	type s struct {
		Text   string  `json:"text"`
		Num1   int     `json:"num1"`
		Num2   float32 `json:"num2"`
		Num3   float64 `json:"num3"`
		Check  bool    `json:"check"`
		Inner1 struct {
			Name string `json:"name" jsonschema:"Some name."`
		} `json:"inner1"`
		Inner2 *struct {
			Name string `json:"name,omitempty"`
		} `json:"inner2"`
	}

	sch, err := For[s]()
	req.Nil(err)

	b, err := json.Marshal(sch)
	req.Nil(err)
	fmt.Println(string(b))

	sch2, err := jsonschema.For[s](nil)
	req.Nil(err)
	b, err = json.Marshal(sch2)
	req.Nil(err)
	fmt.Println(string(b))
}
