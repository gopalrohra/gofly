package flyapi

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gopalrohra/flyapi/env"
	"github.com/gopalrohra/flyapi/sql"
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
	doWork(parseFlags(), config)
}
func doWork(opts cliOptions, config FlyConfig) {
	dbMigration := sql.FlyDBMigration{Migrations: config.Migrations}
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
		log.Println("Initializing cors...")
		c := initializeCors()
		log.Println("Registering routes.")
		router := registerRoutes(c, &config)
		log.Printf("Starting the server on port %s\n", env.Config["SERVER_PORT"])
		log.Println(env.Config)
		handler := http.HandlerFunc(router.HandleRouting)
		serverHost := fmt.Sprintf("%s:%s", env.Config["SERVER_HOST"], env.Config["SERVER_PORT"])
		fmt.Println(serverHost)
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
func registerRoutes(c *cors.Cors, config *FlyConfig) *Router {
	return &Router{Routes: config.Routes, Cors: c}
}
