package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// LDAPServer contains LDAPAddress
type LDAPServer struct {
	URL string `json:"URL"`
}

// Config contains the LdapServer
type Config struct {
	LdapServer                 LDAPServer `json:"LdapServer"`
	BaseDN                     string     `json:"BaseDN"`
	PublicKeyAttribute         string     `json:"PublicKeyAttribute"`
	IgnoreInsecureCertificates bool       `json:"IgnoreInsecureCertificates"`
}

func readConfig(config *Config) error {
	viper.SetConfigName("ldapPubKeyReader")
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/ssh")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	LDAPSERVER := viper.GetString("LdapServer.URL")
	if LDAPSERVER == "" {
		return fmt.Errorf("LdapServer.URL variable is not defined in config file")
	} else {
		config.LdapServer.URL = LDAPSERVER
	}
	BASEDN := viper.GetString("BaseDN")
	if BASEDN == "" {
		return fmt.Errorf("BaseDN variable is not defined in config file")
	} else {
		config.BaseDN = BASEDN
	}
	PUBLICKEYATTRIBUTE := viper.GetString("PublicKeyAttribute")
	if PUBLICKEYATTRIBUTE == "" {
		config.PublicKeyAttribute = "sshPublicKey"
	} else {
		config.PublicKeyAttribute = PUBLICKEYATTRIBUTE
	}
	IGNOREINSECURECERTIFICATES := viper.GetBool("IgnoreInsecureCertificates")
	config.IgnoreInsecureCertificates = IGNOREINSECURECERTIFICATES

	return nil
}
