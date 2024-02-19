// Package configread is for handling the configuration, reading and unmarshaling
package configread

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func Test_readConfig(t *testing.T) {

	tests := []struct {
		name           string
		file           io.Reader
		wantConfigFile MyConfig
		wantErr        bool
	}{
		{
			name: "valid config",
			file: strings.NewReader(`
				[[Accounts]]
				name = "Jane Doe Mailbox"
				mailaddress = "jane.doe@example.com"
				username = ""
				password = "password"
				serveraddress = "mail.example.com"
				serverport = 143
				starttls = true
				additionalfolders = ["Trash", "Spam"]
			 `),
			wantConfigFile: MyConfig{[]*AccountConfig{
				{
					Name:          "Jane Doe Mailbox",
					Mailaddress:   "jane.doe@example.com",
					Username:      "jane.doe@example.com",
					Password:      "password",
					Serveraddress: "mail.example.com",
					Serverport:    143,
					Starttls:      true,
					Folders:       []string{"Trash", "Spam"},
				},
			}},
			wantErr: false,
		},
		{
			name: "valid multi account config",
			file: strings.NewReader(`
				[[Accounts]]
				name = "Jane Doe Mailbox"
				mailaddress = "jane.doe@example.com"
				username = ""
				password = "password"
				serveraddress = "mail.example.com"
				serverport = 143
				starttls = true
				additionalfolders = ["Trash", "Spam"]
				[[Accounts]]
				name = "Jane Mailbox"
				mailaddress = "jane@example.com"
				username = ""
				password = "password"
				serveraddress = "mail.example.com"
				serverport = 143
				starttls = true
				additionalfolders = ["Trash", "Spam"]
			 `),
			wantConfigFile: MyConfig{[]*AccountConfig{
				{
					Name:          "Jane Doe Mailbox",
					Mailaddress:   "jane.doe@example.com",
					Username:      "jane.doe@example.com",
					Password:      "password",
					Serveraddress: "mail.example.com",
					Serverport:    143,
					Starttls:      true,
					Folders:       []string{"Trash", "Spam"},
				},
				{
					Name:          "Jane Mailbox",
					Mailaddress:   "jane@example.com",
					Username:      "jane@example.com",
					Password:      "password",
					Serveraddress: "mail.example.com",
					Serverport:    143,
					Starttls:      true,
					Folders:       []string{"Trash", "Spam"},
				},
			}},
			wantErr: false,
		},
		{
			name: "serverport no int, not valid config, error check",
			file: strings.NewReader(`
				[[Accounts]]
				name = "Jane Doe Mailbox"
				mailaddress = "jane.doe@example.com"
				username = ""
				password = "password"
				serveraddress = "mail.example.com"
				serverport = "143"
				starttls = true
				additionalfolders = ["Trash", "Spam"]
			 `),
			wantConfigFile: MyConfig{[]*AccountConfig{
				{
					Name:          "Jane Doe Mailbox",
					Mailaddress:   "jane.doe@example.com",
					Username:      "jane.doe@example.com",
					Password:      "password",
					Serveraddress: "mail.example.com",
					Serverport:    143,
					Starttls:      true,
					Folders:       []string{"Trash", "Spam"},
				},
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfigFile, err := readConfig(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("readConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false {
				if !reflect.DeepEqual(gotConfigFile, tt.wantConfigFile) {
					t.Errorf("readConfig() = %v, want %v", gotConfigFile, tt.wantConfigFile)
				}
			}
		})
	}
}
