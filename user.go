package goldap

import (
	"fmt"

	"github.com/go-ldap/ldap"
)

// AddUser function
func (l *LDAP) AddUser(name string, password string) error {
	userDN := fmt.Sprintf("CN=%s,%s", name, l.GetUsersDN())
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
