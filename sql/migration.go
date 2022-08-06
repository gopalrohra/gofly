package sql

import (
	"fmt"
	"log"
)

type MigrateFunc = func(Database, string)
type FlyDBMigration struct {
	Migrations map[string]MigrateFunc
}

// CreateDatabase creates database by making grpc call to grpcdb service
func (migration *FlyDBMigration) CreateDatabase() {
	db := NewDatabase()
	err := db.CreateDatabase()
	if err != nil {
		log.Fatalf("Error occured")
		return
	}
	fmt.Println("Database created")
}
func connectionFailure() {
	log.Fatal("Connection to grpcDB service failed")
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
		log.Fatalf("Error occured")
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
			migrateFunc(db, version)
		}
	}
	fmt.Println("Migrations done")
}

func contains(elements []Migration, item string) bool {
	for _, element := range elements {
		if element.Version == item {
			return true
		}
	}
	return false
}
