package main

import (
	"crypto/tls"
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

func doSearch(config *Config, uid string) error {
	if config.LdapServer.URL == "" || config.BaseDN == "" || config.PublicKeyAttribute == "" {
		return fmt.Errorf("invalid LDAP configuration")
	}

	l, err := ldap.DialURL(
		config.LdapServer.URL,
		ldap.DialWithTLSConfig(
			&tls.Config{
				InsecureSkipVerify: config.IgnoreInsecureCertificates,
			},
		),
	)
	if err != nil {
		return fmt.Errorf("failed to dial LDAP server: %w", err)
	}
	defer l.Close()

	filter := fmt.Sprintf("(uid=%s)", uid)
	result, err := l.Search(ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{config.PublicKeyAttribute},
		nil,
	))
	if err != nil {
		return fmt.Errorf("failed to search LDAP server: %w", err)
	}

	if len(result.Entries) == 0 {
		return fmt.Errorf("Filter: [%s] produced no results\n", filter)
	}

	for _, pkey := range result.Entries[0].GetAttributeValues(config.PublicKeyAttribute) {
		fmt.Printf("%s\n", pkey)
	}

	return nil
}
