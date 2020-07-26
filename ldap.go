package goldap

import (
	"log"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// LDAP type
type LDAP struct {
	connection *ldap.Conn
	baseDN     string
	usersDN    string
	groupsDN   string
}

// Init function
func (l *LDAP) Init(baseDN string, usersDN string, groupsDN string) {
	l.baseDN = baseDN
	l.usersDN = usersDN
	l.groupsDN = groupsDN
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
		log.Println(err)
	}
	if !usersDNExists {
		err = l.AddGroup(strings.Split(strings.Split(l.usersDN, ",")[0], "=")[1])
		if err != nil {
			return err
		}
	}
	groupsDNExists, err := l.GroupExists(l.groupsDN)
	if err != nil {
		log.Println(err)
	}
	if !groupsDNExists {
		err = l.AddGroup(strings.Split(strings.Split(l.groupsDN, ",")[0], "=")[1])
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
