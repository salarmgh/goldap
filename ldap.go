package goldap

import (
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// LDAP type
type LDAP struct {
	connection *ldap.Conn
	baseDN     string
	usersDN    string
}

// Init function
func (l *LDAP) Init(baseDN string, usersDN string) {
	l.baseDN = baseDN
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
	usersDNExists, err := l.GroupExists(l.usersDN)
	if err != nil {
		return err
	}
	if !usersDNExists {
		err = l.AddGroup(strings.Split(strings.Split(l.usersDN, ",")[0], "=")[1])
		if err != nil {
			return err
		}
	}
	return nil
}

// NewConn function
func (l *LDAP) NewConn(ldapURL string) (*ldap.Conn, error) {
	ldapConn, err := ldap.DialURL(ldapURL)
	if err != nil {
		return nil, err
	}
	return ldapConn, nil
}
