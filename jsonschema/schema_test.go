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

	type (
		named struct {
			Name string `json:"name" jsonschema:"Some name."`
		}

		s struct {
			Text   string   `json:"text"`
			Num1   int      `json:"num1"`
			Num2   float32  `json:"num2"`
			Num3   float64  `json:"num3"`
			Check  bool     `json:"check"`
			Inner1 named    `json:"inner1"`
			Inner2 *named   `json:"inner2"`
			List1  []string `json:"list1"`
			List2  []named  `json:"list2"`
		}
	)

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
