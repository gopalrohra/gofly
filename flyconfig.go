package flyapi

import (
	"github.com/gopalrohra/flyapi/rest"
	"github.com/gopalrohra/flyapi/sql"
)

type FlyConfig struct {
	Migrations map[string]sql.MigrateFunc
	Routes     []rest.Route
}
