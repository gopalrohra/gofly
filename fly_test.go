package gofly

import "testing"

func TestDoWorkDBInit(t *testing.T) {
	testTable := []struct {
		m       *MockDBMigration
		name    string
		options cliOptions
	}{{
		m:       &MockDBMigration{expectedMethodCall: "Init"},
		name:    "InitTest",
		options: cliOptions{createDB: false, initializeMigration: true, migrateDB: false, start: false},
	}, {
		m:       &MockDBMigration{expectedMethodCall: "CreateDatabase"},
		name:    "CreateDB",
		options: cliOptions{createDB: true, initializeMigration: false, migrateDB: false, start: false},
	}, {
		m:       &MockDBMigration{expectedMethodCall: "MigrateDB"},
		name:    "MigrateDB",
		options: cliOptions{createDB: false, initializeMigration: false, migrateDB: true, start: false},
	},
	}
	for _, table := range testTable {
		doWork(table.m, table.options, buildTestConfig(table.m))
		if table.m.expectedMethodCall != table.m.methodCalled {
			t.Errorf("Test: %s -- Expected %s to be called and %s method was called", table.name, table.m.expectedMethodCall, table.m.methodCalled)
		}
	}
}
