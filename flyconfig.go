package flyapi

import "github.com/gopalrohra/flyapi/sql"

type FlyConfig struct {
	Migrations map[string]sql.MigrateFunc
	Routes     []Route
}
type DBMigration interface {
	Init()
	CreateDatabase()
	MigrateDB()
}
