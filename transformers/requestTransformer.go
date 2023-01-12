package transformers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/gopalrohra/flyapi/log"
)

type RequestTransformer struct {
	request        *http.Request
	routePath      string
	postData       map[string]interface{}
	pathParameters map[string]interface{}
}

func NewRequestTransformer(r *http.Request, path string) RequestTransformer {
	rt := RequestTransformer{request: r, routePath: path}
	rt.parseParameters()
	return rt
}
func (transformer *RequestTransformer) parseParameters() {
	transformer.postData = parsePostParams(transformer.request)
	transformer.pathParameters = parsePathParams(transformer.request, transformer.routePath)
}
func (transformer *RequestTransformer) PopulateData(dest interface{}) {
	e := reflect.ValueOf(dest).Elem()
	log.Debugf("Kind of e: %v\n", e.Kind())
	transformer.processFields(e)
}
func parsePostParams(r *http.Request) map[string]interface{} {
	m := make(map[string]interface{})
	log.Debugf("Value of m: %v", m)
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &m)
	return m
}
func parsePathParams(r *http.Request, routePath string) map[string]interface{} {
	m := make(map[string]interface{})
	if !strings.Contains(routePath, ":") {
		return m
	}
	var requestPathElements, routePathElements []string
	requestPathElements = strings.Split(r.URL.Path, "/")
	routePathElements = strings.Split(routePath, "/")
	if len(requestPathElements) != len(routePathElements) {
		return m
	}
	for i := 0; i < len(requestPathElements); i++ {
		if strings.Contains(routePathElements[i], ":") {
			key := strings.Replace(routePathElements[i], ":", "", 1)
			m[key] = requestPathElements[i]
		}
	}
	log.Debugf("Value of pathParams: %v\n", m)
	return m
}
func (transformer *RequestTransformer) processFields(e reflect.Value) {
	for i := 0; i < e.NumField(); i++ {
		name := e.Type().Field(i).Name
		tag := e.Type().Field(i).Tag
		log.Debugf("Value of name: %v and value of tag: %v\n", name, tag)
		f := e.FieldByName(name)
		log.Debugf("Kind of f: %v and CanSet for f: %v and type of v: %v \n", f.Kind(), f.CanSet(), f.Type().String())
		if f.IsValid() && f.CanSet() && f.Kind() != reflect.Struct {
			transformer.processField(f, tag)
		} else if f.IsValid() && f.CanSet() && f.Kind() == reflect.Struct {
			if _, ok := Transformers[f.Type().String()]; ok {
				log.Debugf("Transformer key: %s\n", f.Type().String())
				transformer.processField(f, tag)
			} else {
				log.Debug("Processing nested struct")
				transformer.processFields(f)
			}
		}
	}
}
func (transformer *RequestTransformer) processField(f reflect.Value, tag reflect.StructTag) {
	var value string
	if tag.Get("requestParamSource") == "query" {
		value = transformer.request.URL.Query().Get(tag.Get("requestParamName"))
	} else if tag.Get("requestParamSource") == "body" {
		if transformer.postData[tag.Get("requestParamName")] != nil {
			value = fmt.Sprint(transformer.postData[tag.Get("requestParamName")])
		}
	} else if tag.Get("requestParamSource") == "path" {
		if transformer.pathParameters[tag.Get("requestParamName")] != nil {
			value = fmt.Sprint(transformer.pathParameters[tag.Get("requestParamName")])
		}
	}
	Transformers[f.Type().String()](f, value)
}
