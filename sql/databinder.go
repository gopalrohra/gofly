package sql

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gopalrohra/flyapi/log"
	"github.com/gopalrohra/flyapi/transformers"
	grpcdb "github.com/gopalrohra/grpcdb/grpc_database"
)

type dataBinder struct {
	target interface{}
}

func (b *dataBinder) bind(sqr *grpcdb.SelectQueryResult, err error) error {
	if err != nil {
		log.Error(err)
		return err
	}
	if "success" != strings.ToLower(sqr.Status) {
		return errors.New("Something went wrong")
	}
	if reflect.TypeOf(b.target).Kind() != reflect.Ptr {
		log.Error("Invalid target, must be ptr")
		return errors.New("Target must be ptr")
	}
	if reflect.TypeOf(b.target).Elem().Kind() == reflect.Slice {
		bindRecords(b.target, sqr)
	} else {
		bindRecord(b.target, sqr)
	}
	return nil
}
func bindRecords(target interface{}, sqr *grpcdb.SelectQueryResult) {
	log.Info("Inside bindRecords")
	s := reflect.ValueOf(target).Elem()
	t := reflect.TypeOf(target).Elem().Elem()
	for _, row := range sqr.Records {
		nv := reflect.New(t).Elem()
		log.Info(nv.Kind())
		processRecord(nv, row)
		s.Set(reflect.Append(s, nv))
	}
}
func processRecord(nv reflect.Value, row *grpcdb.Row) {
	m := ToMap(row)
	for i := 0; i < nv.NumField(); i++ {
		name := nv.Type().Field(i).Name
		tag := nv.Type().Field(i).Tag
		f := nv.FieldByName(name)
		if f.CanSet() && f.Kind() != reflect.Struct && f.Kind() != reflect.Pointer {
			processField(f, tag, m)
		} else {
			if _, ok := transformers.Transformers[f.Type().String()]; ok {
				processField(f, tag, m)
			}
		}
	}
}
func processField(f reflect.Value, tag reflect.StructTag, data map[string]interface{}) {
	if data[tag.Get("dbColumnName")] != nil {
		v := data[tag.Get("dbColumnName")].(string)
		transformers.Transformers[f.Type().String()](f, v)
	}
}
func bindRecord(target interface{}, sqr *grpcdb.SelectQueryResult) {
	log.Info("Inside bindRecord")
	for _, row := range sqr.Records {
		nv := reflect.ValueOf(target).Elem()
		log.Info(nv.Kind())
		processRecord(nv, row)
	}
}
func ToMap(record *grpcdb.Row) map[string]interface{} {
	result := map[string]interface{}{}
	for _, column := range record.Columns {
		result[column.ColumnName] = column.ColumnValue
	}
	log.Info("ToMap: ", result)
	return result
}
