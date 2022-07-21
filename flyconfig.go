package flyapi

type FlyConfig struct {
	Migration DBMigration
	Routes    []Route
}
type DBMigration interface {
	Init()
	CreateDatabase()
	MigrateDB()
}
