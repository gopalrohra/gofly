package gofly

import (
	"github.com/gopalrohra/gofly/rest"
	"github.com/gopalrohra/gofly/sql"
)

type FlyConfig struct {
	Migrations map[string]sql.MigrateFunc
	Routes     []rest.Route
}
