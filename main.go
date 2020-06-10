package main

import (
	"log"

	"github.com/go-ldap/ldap/v3"
)

func main() {
	modify := ldap.NewModifyRequest("CN=folan,DC=digiops,DC=com", []ldap.Control{})
	modify.Add("member", []string{"CN=haha,CN=folan,DC=digiops,DC=com"})

	err = l.Modify(modify)
	if err != nil {
		log.Fatal(err)
	}

	modify = ldap.NewModifyRequest("CN=haha,CN=folan,DC=digiops,DC=com", []ldap.Control{})
	modify.Add("memberOf", []string{"CN=folan,DC=digiops,DC=com"})

	err = l.Modify(modify)
	if err != nil {
		log.Fatal(err)
	}
}
