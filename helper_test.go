package flyapi

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

func buildTestConfig(m DBMigration) FlyConfig {
	return FlyConfig{
		Migration: m,
	}
}
