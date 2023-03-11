package sql

import (
	"time"

	"github.com/gopalrohra/gofly/log"
)

type MigrateFunc = func(Database, string)
type DBMigration interface {
	Init()
	CreateDatabase()
	MigrateDB()
}

type FlyDBMigration struct {
	Migrations map[string]MigrateFunc
}

// CreateDatabase creates database by making grpc call to grpcdb service
func (migration *FlyDBMigration) CreateDatabase() {
	db := NewDatabase()
	err := db.CreateDatabase()
	if err != nil {
		log.Fatal("Error occured")
		return
	}
	log.Fatal("Database created")
}

// Init initializes migration for the application
func (migration *FlyDBMigration) Init() {
	db := NewDatabase()
	columns := []string{
		"id serial primary key",
		"version text not null",
		"creation_time timestamp",
	}
	err := db.CreateTable("migrations", columns)
	if err != nil {
		log.Fatalf("Error occured %v\n", err)
	}
}

// MigrateDB migrates the database to the latest changes
func (migration *FlyDBMigration) MigrateDB() {
	db := NewDatabase()
	var m []Migration
	err := db.Scan(&m)
	if err != nil {
		log.Fatal(err)
	}

	for version, migrateFunc := range migration.Migrations {
		if !contains(m, version) {
			//todo: handle in single transaction
			migrateFunc(db, version)
			updateMigration(db, version)
		}
	}
	log.Info("Migrations done")
}
func updateMigration(db Database, version string) {
	m := Migration{Version: version, CreationTime: time.Now()}
	err := db.Insert(&m)
	if err != nil {
		log.Fatal(err)
	}
}
func contains(elements []Migration, item string) bool {
	for _, element := range elements {
		if element.Version == item {
			return true
		}
	}
	return false
}
