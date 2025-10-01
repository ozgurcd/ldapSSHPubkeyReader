package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
)

func doSearch(config *Config, username string) error {
	// Input sanitization - prevent LDAP injection
	username = sanitizeLDAPInput(username)
	if username == "" {
		return fmt.Errorf("invalid username provided")
	}

	var l *ldap.Conn
	var err error

	// Retry logic for connection
	for attempt := 0; attempt <= config.LdapServer.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(config.LdapServer.RetryDelay)
		}

		l, err = connectToLDAP(config)
		if err == nil {
			break
		}

		if attempt == config.LdapServer.MaxRetries {
			return fmt.Errorf("failed to connect to LDAP server after %d attempts: %w", config.LdapServer.MaxRetries+1, err)
		}
	}
	defer l.Close()

	// Set connection timeout for operations
	l.SetTimeout(config.LdapServer.ConnectionTimeout)

	// Bind if credentials are provided
	if config.LdapServer.BindDN != "" && config.LdapServer.BindPassword != "" {
		err = l.Bind(config.LdapServer.BindDN, config.LdapServer.BindPassword)
		if err != nil {
			return fmt.Errorf("failed to bind to LDAP server: %w", err)
		}
	}

	// Build search filter
	filter := fmt.Sprintf(config.SearchFilter, config.UserAttribute, username)

	if config.Debug {
		fmt.Fprintf(os.Stderr, "Performing LDAP search with filter: %s, timeout: %v\n", filter, config.LdapServer.SearchTimeout)
	}

	// Create search request with timeout
	searchRequest := ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		int(config.LdapServer.SearchTimeout.Seconds()),
		false,
		filter,
		[]string{config.PublicKeyAttribute},
		nil,
	)

	result, err := l.Search(searchRequest)
	if err != nil {
		return fmt.Errorf("failed to search LDAP server: %w", err)
	}

	if len(result.Entries) == 0 {
		return fmt.Errorf("no results found for user: %s", username)
	}

	// Output all public keys for the user
	for _, pkey := range result.Entries[0].GetAttributeValues(config.PublicKeyAttribute) {
		if strings.TrimSpace(pkey) != "" {
			fmt.Printf("%s\n", pkey)
		}
	}

	return nil
}

// connectToLDAP establishes a connection to the LDAP server with timeout
func connectToLDAP(config *Config) (*ldap.Conn, error) {
	if config.Debug {
		fmt.Fprintf(os.Stderr, "Connecting to LDAP server with timeout: %v\n", config.LdapServer.ConnectionTimeout)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: config.TLS.InsecureSkipVerify,
	}

	// Load custom certificates if provided
	if config.TLS.CAFile != "" {
		// TODO: Implement custom CA certificate loading
	}

	// Create a custom dialer with timeout
	dialer := &net.Dialer{
		Timeout: config.LdapServer.ConnectionTimeout,
	}

	// Use DialURL with timeout configuration
	l, err := ldap.DialURL(
		config.LdapServer.URL,
		ldap.DialWithTLSConfig(tlsConfig),
		ldap.DialWithDialer(dialer),
	)
	if err != nil {
		return nil, err
	}

	if config.Debug {
		fmt.Fprintf(os.Stderr, "Successfully connected to LDAP server\n")
	}

	return l, nil
}

// sanitizeLDAPInput sanitizes user input to prevent LDAP injection attacks
func sanitizeLDAPInput(input string) string {
	// Remove potentially dangerous characters
	dangerous := []string{
		"*", "(", ")", "\\", "/", "\x00", "&", "|", "!",
		"=", "<", ">", "~", ";", ",", "+", "\"", "'",
	}

	sanitized := input
	for _, char := range dangerous {
		sanitized = strings.ReplaceAll(sanitized, char, "")
	}

	return strings.TrimSpace(sanitized)
}
