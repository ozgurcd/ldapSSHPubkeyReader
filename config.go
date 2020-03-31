package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// LDAPServer contains LDAPAddress
type LDAPServer struct {
	URL string `json:"URL"`
}

// Config contains the LdapServer
type Config struct {
	LdapServer LDAPServer `json:"LdapServer"`
	BaseDN     string     `json:"BaseDN"`
}

func readConfig(config *Config) {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
}