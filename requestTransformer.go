package flyapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/gopalrohra/flyapi/util"
)

type RequestTransformer struct {
	request  *http.Request
	postData map[string]interface{}
}

func (transformer RequestTransformer) populateData(dest interface{}) {
	transformer.postData = parsePostParams(transformer.request)
	e := reflect.ValueOf(dest).Elem()
	fmt.Printf("Kind of e: %v\n", e.Kind())
	transformer.processFields(e)
}
func parsePostParams(r *http.Request) map[string]interface{} {
	m := make(map[string]interface{})
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &m)
	return m
}
func (transformer RequestTransformer) processFields(e reflect.Value) {
	for i := 0; i < e.NumField(); i++ {
		name := e.Type().Field(i).Name
		tag := e.Type().Field(i).Tag
		fmt.Printf("Value of name: %v and value of tag: %v\n", name, tag)
		f := e.FieldByName(name)
		fmt.Printf("Kind of f: %v and CanSet for f: %v and type of v: %v \n", f.Kind(), f.CanSet(), f.Type().String())
		if f.IsValid() && f.CanSet() && f.Kind() != reflect.Struct {
			transformer.processField(f, tag)
		} else if f.IsValid() && f.CanSet() && f.Kind() == reflect.Struct {
			if _, ok := transformers[f.Type().String()]; ok {
				fmt.Printf("Transformer key: %s\n", f.Type().String())
				transformer.processField(f, tag)
			} else {
				fmt.Println("Processing nested struct")
				transformer.processFields(f)
			}
		}
	}
}
func (transformer RequestTransformer) processField(f reflect.Value, tag reflect.StructTag) {
	var value string
	if tag.Get("requestParamSource") == "query" {
		value = transformer.request.URL.Query().Get(tag.Get("requestParamName"))
	} else if tag.Get("requestParamSource") == "body" {
		if transformer.postData[tag.Get("requestParamName")] != nil {
			value = fmt.Sprint(transformer.postData[tag.Get("requestParamName")])
		}
	}
	transformers[f.Type().String()](f, value)
}

type TransformerFunc = func(reflect.Value, string)

func intTransformer(f reflect.Value, v string) {
	fmt.Printf("Value of request param:%v\n", v)
	intV := util.ToInt(v)
	if !f.OverflowInt(int64(intV)) {
		f.SetInt(int64(intV))
	}
}
func stringTransformer(f reflect.Value, v string) {
	fmt.Println(v)
	f.SetString(v)
}

var transformers = map[string]TransformerFunc{
	"int":       intTransformer,
	"string":    stringTransformer,
	"time.Time": timeTransformer,
}

func timeTransformer(f reflect.Value, v string) {
	f.Set(reflect.ValueOf(util.ToDate(v)))
}
