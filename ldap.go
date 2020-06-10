package goldap

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

// LDAP type
type LDAP struct {
	connection *ldap.Conn
	baseDN     string
	usersCN    string
}

// GetUsersDN function
func (l *LDAP) GetUsersDN() string {
	return fmt.Sprintf("%s,%s", l.usersCN, l.baseDN)
}

// Init function
func (l *LDAP) Init(bindDN string, usersDN string) {
	l.bindDN = bindDN
	l.usersDN = usersDN
}

// GetConn function
func (l *LDAP) GetConn(ldapURL string, bindUser string, bindPass string) error {
	ldapConn, err := ldap.DialURL(ldapURL)
	if err != nil {
		return err
	}
	err = ldapConn.Bind(bindUser, bindPass)
	if err != nil {
		return err
	}
	l.connection = ldapConn
}
