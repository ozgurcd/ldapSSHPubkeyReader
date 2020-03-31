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
		ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	filter := "(uid=" + uid + ")"
	result, err := l.Search(ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{"sshPublicKey"},
		nil,
	))

	if err != nil {
		log.Fatal(err)
		os.Exit(-2)
	}

	//users := []string{}
	for _, pkey := range result.Entries[0].GetAttributeValues("sshPublicKey") {
		fmt.Printf("%s\n", pkey)
	}
}
