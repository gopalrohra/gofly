package flyapi

import (
	"fmt"
	"reflect"

	"github.com/gopalrohra/flyapi/transformers"
)

type MapTransformer struct {
	Data map[string]interface{}
}

func (transformer MapTransformer) PopulateData(dest interface{}) {
	e := reflect.ValueOf(dest).Elem()
	fmt.Printf("Kind of e: %v\n", e.Kind())
	transformer.processFields(e)
}
func (transformer MapTransformer) processFields(e reflect.Value) {
	for i := 0; i < e.NumField(); i++ {
		name := e.Type().Field(i).Name
		tag := e.Type().Field(i).Tag
		fmt.Printf("Value of name: %v and value of tag: %v\n", name, tag)
		f := e.FieldByName(name)
		fmt.Printf("Kind of f: %v and CanSet for f: %v and type of v: %v \n", f.Kind(), f.CanSet(), f.Type().String())
		if f.IsValid() && f.CanSet() && f.Kind() != reflect.Struct {
			transformer.processField(f, tag)
		} else if f.IsValid() && f.CanSet() && f.Kind() == reflect.Struct {
			if _, ok := transformers.Transformers[f.Type().String()]; ok {
				fmt.Printf("Transformer key: %s\n", f.Type().String())
				transformer.processField(f, tag)
			} else {
				fmt.Println("Processing nested struct")
				transformer.processFields(f)
			}
		}
	}
}
func (transformer MapTransformer) processField(f reflect.Value, tag reflect.StructTag) {
	fmt.Printf("%v: %v\n", tag.Get("dbColumnName"), transformer.Data[tag.Get("dbColumnName")])
	if transformer.Data[tag.Get("dbColumnName")] != nil {
		value := fmt.Sprint(transformer.Data[tag.Get("dbColumnName")])
		transformers.Transformers[f.Type().String()](f, value)
	}
}
