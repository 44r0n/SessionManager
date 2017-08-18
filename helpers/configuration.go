package helpers

import (
	"log"

	"github.com/tkanos/gonfig"
)

// Configuration type to read configuration file
type Configuration struct {
	IP         string
	Port       int
	ConnString string
}

var configuration Configuration
var initialized = false

// GetConnString function. Gets the configuration string of MySQL from a given json formated file
func GetConnString(configFileName string) string {
	loadConfig(configFileName)
	return configuration.ConnString
}

func loadConfig(configFileName string) {
	// get Configuration
	if !initialized {
		err := gonfig.GetConf(configFileName, &configuration)
		if err != nil {
			log.Fatalf("Error loading file %s: %s", configFileName, err)
		}
		initialized = true
	}
}

// GetIP function gets the ip of the configurated MySQL database from a given json formated file
func GetIP(configFileName string) string {
	loadConfig(configFileName)
	return configuration.IP
}
