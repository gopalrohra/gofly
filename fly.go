package flyapi

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gopalrohra/flyapi/env"
	"github.com/rs/cors"
)

func Fly(config FlyConfig) {
	var createDB = flag.Bool("createDB", false, "Creates a new database for tracking expenses")
	var initializeMigration = flag.Bool("initMigration", false, "Initializes migration feature for the application")
	var migrateDB = flag.Bool("migrateDB", false, "Runs the migration for the attendance database")
	var start = flag.Bool("start", true, "Starts the expense tracker rest service")
	flag.Parse()
	env.LoadEnvironment()
	if *createDB {
		config.Migration.CreateDatabase()
		return
	}
	if *initializeMigration {
		config.Migration.Init()
		return
	}
	if *migrateDB {
		config.Migration.MigrateDB()
		return
	}
	if *start {
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
