package sql

import (
	"reflect"
	"strings"
	"time"

	"github.com/gopalrohra/flyapi/log"
	"github.com/gopalrohra/flyapi/util"

	grpcdb "github.com/gopalrohra/grpcdb/grpc_database"
)

type queryBuilder struct{}

func (qb *queryBuilder) updateQuery(dbInfo *grpcdb.DatabaseInfo, i interface{}, clauses []string) *grpcdb.UpdateQuery {
	uq := &grpcdb.UpdateQuery{
		Info:      dbInfo,
		TableName: getTableName(i),
		Columns:   getUpdateColumns(i),
		Clauses:   clauses,
	}
	return uq
}
func (qb *queryBuilder) insertQuery(dbInfo *grpcdb.DatabaseInfo, i interface{}) *grpcdb.InsertQueryRequest {
	insertQuery := &grpcdb.InsertQueryRequest{
		Info:      dbInfo,
		TableName: getTableName(i),
	}
	insertQuery.Columns, insertQuery.ColumnValues = getInsertColumnNameValues(i)
	insertQuery.ReturningIdColumnName = "id"
	return insertQuery
}
func (qb *queryBuilder) selectQuery(dbInfo *grpcdb.DatabaseInfo, i interface{}, queryClauses ...[]string) *grpcdb.SelectQuery {
	var t reflect.Type
	if reflect.TypeOf(i).Elem().Kind() == reflect.Slice {
		t = reflect.TypeOf(i).Elem().Elem()
	} else {
		t = reflect.TypeOf(i).Elem()
	}
	log.Debug(t)
	fields := extractFields(t)
	sq := &grpcdb.SelectQuery{Info: dbInfo, Fields: fields, TableName: getTableName(i)}
	if len(queryClauses) > 0 {
		sq.Clauses = queryClauses[0]
	} else if len(queryClauses) > 1 {
		sq.Groupby = queryClauses[1]
	} else if len(queryClauses) > 2 {
		sq.Orderby = queryClauses[2]
	}
	log.Debug(sq)
	return sq
}
func getTableName(i interface{}) string {
	var t reflect.Type
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		log.Debug("Inside ptr condition")
		if reflect.TypeOf(i).Elem().Kind() == reflect.Slice {
			log.Debug("Found slice in getTableName")
			t = reflect.TypeOf(i).Elem().Elem()
		} else {
			log.Debug("Inside struct condition")
			t = reflect.TypeOf(i).Elem()
		}
	}

	nv := reflect.New(t)
	rv := nv.MethodByName("GetTableName").Call(nil)
	log.Debug(rv)
	return rv[len(rv)-1].String()
}

func extractFields(t reflect.Type) []string {
	//t should be struct
	var fields []string
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		if tag.Get("dbColumnName") != "" {
			fields = append(fields, tag.Get("dbColumnName"))
		}
	}
	return fields
}
func getInsertColumnNameValues(entity interface{}) ([]string, []string) {
	var columnNames, columnValues []string
	e := reflect.ValueOf(entity).Elem()
	for i := 0; i < e.NumField(); i++ {
		name := e.Type().Field(i).Name
		tag := e.Type().Field(i).Tag
		f := e.FieldByName(name)
		if tag.Get("dbColumnName") != "" && operationAllowed(tag.Get("dbOperations"), "insert") {
			log.Debugf("Entity field name: %v\n", name)
			columnNames = append(columnNames, tag.Get("dbColumnName"))
			if f.Kind() == reflect.String {
				columnValues = append(columnValues, util.SQ(f.String()))
			} else if f.Kind() == reflect.Int {
				columnValues = append(columnValues, util.ToString(int(f.Int())))
			} else if f.Type().String() == "time.Time" {
				columnValues = append(columnValues, util.SQ(f.Interface().(time.Time).Format(time.RFC3339))+"::timestamp")
			}
		}
	}
	log.Debugf("Columns: %v and ColumnValues: %v\n", columnNames, columnValues)
	return columnNames, columnValues
}
func operationAllowed(operations string, operation string) bool {
	operationList := strings.Split(operations, ",")
	for _, val := range operationList {
		if val == operation {
			return true
		}
	}
	return false
}

func getUpdateColumns(entity interface{}) []*grpcdb.Column {
	var columns []*grpcdb.Column
	e := reflect.ValueOf(entity).Elem()
	for i := 0; i < e.NumField(); i++ {
		name := e.Type().Field(i).Name
		tag := e.Type().Field(i).Tag
		f := e.FieldByName(name)
		if tag.Get("dbColumnName") != "" && operationAllowed(tag.Get("dbOperations"), "update") {
			log.Debugf("Entity field name: %v\n", name)
			column := grpcdb.Column{ColumnName: tag.Get("dbColumnName")}
			if f.Kind() == reflect.String {
				column.ColumnValue = util.SQ(f.String())
			} else if f.Kind() == reflect.Int {
				column.ColumnValue = util.ToString(int(f.Int()))
			} else if f.Type().String() == "time.Time" {
				column.ColumnValue = util.SQ(f.Interface().(time.Time).Format(time.RFC3339)) + "::timestamp"
			}
			columns = append(columns, &column)
		}
	}
	log.Debugf("Columns: %v\n", columns)
	return columns
}
