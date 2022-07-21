package env

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Config to hold environment based configuration
var Config = map[string]string{}

// LoadEnvironment loads .env for production and .env.<env_name> for specific environment
func LoadEnvironment() {
	log.Println("Initializing environment.")
	//set mode to production, development, staging etc
	mode := os.Getenv("GO_ENV")
	if mode == "" {
		mode = "development"
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(dir)
	env_file := dir + "/.env"
	if mode != "production" {
		env_file = env_file + "." + mode
	}
	data, err := ioutil.ReadFile(env_file)
	if err != nil {
		log.Println(err.Error())
	}
	lines := strings.Split(fmt.Sprintf("%s", data), "\n")
	log.Println(len(lines))
	for _, line := range lines {
		// fmt.Println("Inside loop")
		// fmt.Println(line)
		kvPair := strings.Split(line, "=")
		if len(kvPair) == 2 {
			if strings.Trim(kvPair[0], " ") != "" && strings.Trim(kvPair[1], " ") != "" {
				Config[kvPair[0]] = kvPair[1]
			}
		}
	}
	log.Println("Finished loading the environment values")
	fmt.Println(Config)
}
