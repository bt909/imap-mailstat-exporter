// Package configread is for handling the configuration, reading and unmarshaling
package configread

import (
	"os"

	"github.com/BurntSushi/toml"
)

// a generic reusable error handling, needs to be enhanced
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

type AccountConfig struct {
	Name          string   `toml:"name"`
	Mailaddress   string   `toml:"mailaddress"`
	Username      string   `toml:"username"`
	Password      string   `toml:"password"`
	Serveraddress string   `toml:"serveraddress"`
	Serverport    int      `toml:"serverport"`
	Starttls      bool     `toml:"starttls"`
	Folders       []string `toml:"additionalfolders"`
}

type MyConfig struct {
	Accounts []*AccountConfig
}

// the main function for reading and unmarshaling the configfile
func GetConfig(configfile string) MyConfig {
	file, err := os.ReadFile(configfile)
	checkError(err)
	var config MyConfig
	err = toml.Unmarshal(file, &config)
	return config
}
