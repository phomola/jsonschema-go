package jsonschema

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchemaInference(t *testing.T) {
	req := require.New(t)

	sch, err := For[struct {
		Text   string
		Num1   int
		Num2   float32
		Num3   float64
		Check  bool
		Inner1 struct {
			Name string `json:"name" jsonschema:"Some name."`
		}
		Inner2 *struct {
			Name string `json:"name,omitempty"`
		}
	}]()
	req.Nil(err)
	b, err := json.Marshal(sch)
	req.Nil(err)
	fmt.Println(string(b))
}
