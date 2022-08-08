package flyapi

import "github.com/gopalrohra/flyapi/sql"

type MockDBMigration struct {
	methodCalled       string
	expectedMethodCall string
}

func (m *MockDBMigration) Init() {
	m.methodCalled = "Init"
}
func (m *MockDBMigration) CreateDatabase() {
	m.methodCalled = "CreateDatabase"
}
func (m *MockDBMigration) MigrateDB() {
	m.methodCalled = "MigrateDB"
}

func buildTestConfig(m sql.DBMigration) FlyConfig {
	return FlyConfig{Migrations: make(map[string]func(sql.Database, string))}
}
