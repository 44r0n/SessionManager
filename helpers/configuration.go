package helpers

import (
  "github.com/tkanos/gonfig"
)

// Configuration type to read configuration file
type Configuration struct {
  Port int
  ConnString string
}

// GetConnString function. Gets the configuration string of MySQL from a given json formated file
func GetConnString(configFileName string) string {
  // get Configuration
  configuration := Configuration{}
  err := gonfig.GetConf(configFileName,&configuration)
  if err != nil {
    panic(err)
  }

  return configuration.ConnString
}
