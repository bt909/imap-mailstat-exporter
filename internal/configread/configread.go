// Package configread is for handling the configuration, reading and unmarshaling
package configread

import (
	"bytes"
	"io"
	"os"
	"sync"

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

// the main function for reading
func GetConfig(configfile string) (MyConfig, error) {
	file, err := os.Open(configfile)
	checkError(err)
	config, err := readConfig(file)
	return config, err
}

// the function for unmarshaling the configfile
func readConfig(file io.Reader) (config MyConfig, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	err = toml.Unmarshal(buf.Bytes(), &config)
	if err != nil {
		return config, err
	}
	sliceLength := len(config.Accounts)
	var wg sync.WaitGroup
	wg.Add(sliceLength)
	for account := range config.Accounts {
		go func(account int) {
			defer wg.Done()
			if config.Accounts[account].Username == "" {
				config.Accounts[account].Username = config.Accounts[account].Mailaddress
			}
		}(account)
	}
	wg.Wait()
	return config, err
}
