package configread

import (
	"os"

	"github.com/BurntSushi/toml"
)

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

func GetConfig() MyConfig {
	file, err := os.ReadFile("./config/config.toml")
	checkError(err)
	var config MyConfig
	err = toml.Unmarshal(file, &config)
	return config
}
