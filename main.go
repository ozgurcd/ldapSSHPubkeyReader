package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <username> [config-path...]\n", args[0])
		fmt.Fprintf(os.Stderr, "  username: LDAP username to search for\n")
		fmt.Fprintf(os.Stderr, "  config-path: Optional additional config search paths\n")
		fmt.Fprintf(os.Stderr, "\nEnvironment variables:\n")
		fmt.Fprintf(os.Stderr, "  LDAP_SSH_LDAP_SERVER_URL: LDAP server URL\n")
		fmt.Fprintf(os.Stderr, "  LDAP_SSH_BASE_DN: LDAP base DN\n")
		fmt.Fprintf(os.Stderr, "  LDAP_SSH_DEBUG: Enable debug mode (true/false)\n")
		os.Exit(1)
	}

	username := args[1]
	var configPaths []string
	if len(args) > 2 {
		configPaths = args[2:]
	}

	config := Config{
		ConfigPaths: configPaths,
	}

	err := readConfig(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %s\n", err.Error())
		if config.Debug {
			fmt.Fprintf(os.Stderr, "Config file used: %s\n", GetConfigFile())
		}
		os.Exit(1)
	}

	if config.Debug {
		fmt.Fprintf(os.Stderr, "Using config file: %s\n", GetConfigFile())
		fmt.Fprintf(os.Stderr, "Searching for user: %s\n", username)
		fmt.Fprintf(os.Stderr, "LDAP Server: %s\n", config.LdapServer.URL)
		fmt.Fprintf(os.Stderr, "Base DN: %s\n", config.BaseDN)
		fmt.Fprintf(os.Stderr, "Connection timeout: %v, Search timeout: %v\n", config.LdapServer.ConnectionTimeout, config.LdapServer.SearchTimeout)
		fmt.Fprintf(os.Stderr, "Max retries: %d, Retry delay: %v\n", config.LdapServer.MaxRetries, config.LdapServer.RetryDelay)
	}

	err = doSearch(&config, username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Search error: %s\n", err.Error())
		os.Exit(2)
	}
}
