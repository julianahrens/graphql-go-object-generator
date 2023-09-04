package graphql_go_struct_schema

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"reflect"
	"slices"
	"strings"
)

type StructObjectGenerator struct {
	ScalarMapping map[reflect.Kind]*graphql.Scalar
}

func NewStructObjectGenerator() *StructObjectGenerator {
	return &StructObjectGenerator{
		ScalarMapping: map[reflect.Kind]*graphql.Scalar{
			reflect.Bool:       graphql.Boolean,
			reflect.Int:        graphql.Int,
			reflect.Int8:       graphql.Int,
			reflect.Int16:      graphql.Int,
			reflect.Int64:      graphql.Int,
			reflect.Uint:       graphql.Int,
			reflect.Uint8:      graphql.Int,
			reflect.Uint16:     graphql.Int,
			reflect.Uint32:     graphql.Int,
			reflect.Uint64:     graphql.Int,
			reflect.Uintptr:    graphql.String,
			reflect.Float32:    graphql.Float,
			reflect.Float64:    graphql.Float,
			reflect.Complex64:  graphql.Float,
			reflect.Complex128: graphql.Float,
			reflect.String:     graphql.String,
		},
	}
}

func (g *StructObjectGenerator) GenerateObject(s interface{}) (*graphql.Object, *graphql.InputObject, error) {
	return g.generateObjectsByType(reflect.TypeOf(s))
}

func (g *StructObjectGenerator) generateObjectsByType(t reflect.Type) (*graphql.Object, *graphql.InputObject, error) {
	if t.Kind() == reflect.Pointer {
		return g.generateObjectsByType(t.Elem())
	} else if t.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("input has to be a struct")
	} else {
		return g.generateOutputObject(t), g.generateInputObject(t), nil
	}
}

func (g *StructObjectGenerator) generateOutputObject(t reflect.Type) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:   t.Name(),
		Fields: g.getOutputFields(t),
	})
}

func (g *StructObjectGenerator) generateInputObject(t reflect.Type) *graphql.InputObject {
	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   fmt.Sprintf("%sInput", t.Name()),
		Fields: g.getInputFields(t),
	})
}

func (g *StructObjectGenerator) getOutputFields(t reflect.Type) graphql.Fields {
	fields := make(map[string]*graphql.Field)

	for i := 0; i < t.NumField(); i++ {
		name, omitempty := getJSONTag(t.Field(i).Tag)
		fields[name] = &graphql.Field{
			Type: g.getType(t.Field(i).Type, omitempty),
		}
	}

	return fields
}

func (g *StructObjectGenerator) getInputFields(t reflect.Type) graphql.InputObjectFieldMap {
	fields := make(map[string]*graphql.InputObjectField)

	for i := 0; i < t.NumField(); i++ {
		if readonly := getGraphqlTag(t.Field(i).Tag); readonly {
			continue
		}
		name, omitempty := getJSONTag(t.Field(i).Tag)
		fields[name] = &graphql.InputObjectField{
			Type: g.getType(t.Field(i).Type, omitempty),
		}
	}

	return fields
}

func (g *StructObjectGenerator) getType(t reflect.Type, isOmitempty bool) graphql.Output {
	if !isOmitempty {
		return graphql.NewNonNull(g.getType(t, true))
	} else if slices.Contains([]reflect.Kind{reflect.Array, reflect.Slice}, t.Kind()) {
		return graphql.NewList(g.getType(t.Elem(), isOmitempty))
	} else if t.Kind() == reflect.Pointer {
		return g.getType(t.Elem(), isOmitempty)
	} else {
		return g.ScalarMapping[t.Kind()]
	}
}

func getJSONTag(tag reflect.StructTag) (string, bool) {
	json := strings.Split(tag.Get("json"), ",")
	return json[0], slices.Contains(json, "omitempty")
}

func getGraphqlTag(tag reflect.StructTag) bool {
	values := strings.Split(tag.Get("graphql"), ",")
	return slices.Contains(values, "readonly")
}
