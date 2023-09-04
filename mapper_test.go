package graphql_go_struct_schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type A struct {
	One []*string `json:"one"`
	Two int64     `json:"two,omitempty" graphql:"readonly"`
}

func TestGenerateObject(t *testing.T) {
	a := assert.New(t)

	g := NewStructObjectGenerator()
	aOutputType, aInputType, err := g.GenerateObject((*A)(nil))

	a.NotNil(err)

	a.Equal("A", aOutputType.Name())
	a.Equal("AInput", aInputType.Name())

	a.Equal(2, len(aOutputType.Fields()))
}
