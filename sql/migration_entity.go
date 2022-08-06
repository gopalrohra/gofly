package sql

import "time"

type Migration struct {
	ID           int       `dbColumnName:"id"`
	Version      string    `dbColumnName:"version" dbOperations:"insert"`
	CreationTime time.Time `dbColumnName:"creation_time" dbOperations:"insert"`
}

func (m *Migration) GetTableName() string {
	return "migrations"
}
