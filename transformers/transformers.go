package transformers

import (
	"reflect"

	"github.com/gopalrohra/gofly/log"
	"github.com/gopalrohra/gofly/util"
)

type TransformerFunc = func(reflect.Value, string)

var Transformers = map[string]TransformerFunc{
	"int":        IntTransformer,
	"string":     StringTransformer,
	"time.Time":  TimeTransformer,
	"*time.Time": TimePointerTransformer,
}

func IntTransformer(f reflect.Value, v string) {
	log.Debugf("Value of request param:%v\n", v)
	intV := util.ToInt(v)
	if !f.OverflowInt(int64(intV)) {
		f.SetInt(int64(intV))
	}
}
func StringTransformer(f reflect.Value, v string) {
	log.Debugf("Inside string transformer with value: %v for field: %v\n", v, f.String())
	f.SetString(v)
}
func TimeTransformer(f reflect.Value, v string) {
	f.Set(reflect.ValueOf(util.ToDate(v)))
}
func TimePointerTransformer(f reflect.Value, v string) {
	if v != "" && v != "<nil>" {
		t := util.ToDate(v)
		f.Set(reflect.ValueOf(&t))
	} else {
		log.Info("Time pointer")
		f.Set(reflect.Zero(f.Type()))
	}
}
