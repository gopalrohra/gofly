package sql

import (
	"github.com/gopalrohra/flyapi/log"
	grpcdb "github.com/gopalrohra/grpcdb/grpc_database"
)

type Query struct {
	mode          string
	fields        []string
	tableName     string
	whereClauses  []string
	orderBy       []string
	groupBy       []string
	limit         int
	offset        int
	info          *grpcdb.DatabaseInfo
	executeSelect func(sq *grpcdb.SelectQuery) (*grpcdb.SelectQueryResult, error)
}

func (q *Query) Select() *Query {
	q.mode = "select"
	return q
}
func (q *Query) AddSelect(field string) *Query {
	q.fields = append(q.fields, field)
	return q
}
func (q *Query) From(table string) *Query {
	q.tableName = table
	return q
}
func (q *Query) Where(clause string) *Query {
	q.whereClauses = append(q.whereClauses, clause)
	return q
}
func (q *Query) OrderBy(field string) *Query {
	q.orderBy = append(q.orderBy, field)
	return q
}
func (q *Query) GroupBy(field string) *Query {
	q.groupBy = append(q.groupBy, field)
	return q
}
func (q *Query) Limit(value int) *Query {
	q.limit = value
	return q
}
func (q *Query) Offset(value int) *Query {
	q.offset = value
	return q
}
func (q *Query) build() *grpcdb.SelectQuery {
	sq := &grpcdb.SelectQuery{Info: q.info, Fields: q.fields, TableName: q.tableName}
	if len(q.whereClauses) > 0 {
		sq.Clauses = q.whereClauses
	}
	if len(q.groupBy) > 0 {
		sq.Groupby = q.groupBy
	}
	if len(q.orderBy) > 0 {
		sq.Orderby = q.orderBy
	}
	log.Debugf("Select query: %v", sq)
	return sq
}
func (q *Query) ExecuteSelect(i interface{}) error {
	binder := dataBinder{target: i}
	return binder.bind(q.executeSelect(q.build()))
}
