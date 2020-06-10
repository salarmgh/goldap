package goldap

import (
	"fmt"

	"github.com/go-ldap/ldap"
)

func (l *LDAP) AddUser(name string, password string) error {
	userDN := fmt.Sprintf("CN=%s,CN=people,DC=digiops,DC=com", name)
	addReq := ldap.NewAddRequest(userDN, []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "person", "uidObject", "simpleSecurityObject"})
	addReq.Attribute("cn", []string{name})
	addReq.Attribute("sn", []string{name})
	addReq.Attribute("uid", []string{uid})
	addReq.Attribute("userPassword", []string{password})

	if err := l.Add(addReq); err != nil {
		return err
	}
}
