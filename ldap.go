package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/go-ldap/ldap"
)

func doSearch(config *Config, uid string) {
	l, err := ldap.DialURL(
		config.LdapServer.URL,
		ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: config.IgnoreInsecureCertificates}))
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
		os.Exit(-2)
	}

	if len(result.Entries) == 0 {
		log.Fatalf("Filter: [%s] produced no results\n", filter)
		os.Exit(-1)
	}

	for _, pkey := range result.Entries[0].GetAttributeValues(config.PublicKeyAttribute) {
		fmt.Printf("%s\n", pkey)
	}
}
