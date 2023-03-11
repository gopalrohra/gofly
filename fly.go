package gofly

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gopalrohra/gofly/env"
	"github.com/gopalrohra/gofly/log"
	"github.com/gopalrohra/gofly/rest"
	"github.com/gopalrohra/gofly/sql"
	"github.com/rs/cors"
)

type cliOptions struct {
	createDB            bool
	initializeMigration bool
	migrateDB           bool
	start               bool
}

func parseFlags() cliOptions {
	var options cliOptions
	flag.BoolVar(&options.createDB, "createDB", false, "Creates a new database for the database")
	flag.BoolVar(&options.initializeMigration, "initMigration", false, "Initializes migration feature for the application")
	flag.BoolVar(&options.migrateDB, "migrateDB", false, "Runs the database migration for the application")
	flag.BoolVar(&options.start, "start", true, "Starts the http server")
	flag.Parse()
	return options
}
func Fly(config FlyConfig) {
	env.LoadEnvironment()
	dbMigration := sql.FlyDBMigration{Migrations: config.Migrations}
	doWork(&dbMigration, parseFlags(), config)
}
func doWork(dbMigration sql.DBMigration, opts cliOptions, config FlyConfig) {
	if opts.createDB {
		dbMigration.CreateDatabase()
		return
	}
	if opts.initializeMigration {
		dbMigration.Init()
		return
	}
	if opts.migrateDB {
		dbMigration.MigrateDB()
		return
	}
	if opts.start {
		log.Info("Initializing cors...")
		c := initializeCors()
		log.Info("Registering routes.")
		router := registerRoutes(c, &config)
		log.Infof("Starting the server on port %s\n", env.Config["SERVER_PORT"])
		log.Info(env.Config)
		handler := http.HandlerFunc(router.HandleRouting)
		serverHost := fmt.Sprintf("%s:%s", env.Config["SERVER_HOST"], env.Config["SERVER_PORT"])
		log.Info(serverHost)
		log.Fatal(http.ListenAndServe(serverHost, handler))
	}
}
func initializeCors() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{env.Config["ALLOWED_ORIGINS"]},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"HEAD", "GET", "POST", "PUT"},
	})
}
func registerRoutes(c *cors.Cors, config *FlyConfig) *rest.Router {
	return &rest.Router{Routes: config.Routes, Cors: c}
}
