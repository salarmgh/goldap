package goldap

import "github.com/go-ldap/ldap/v3"

// LDAP type
type LDAP struct {
	connection *ldap.Conn
}

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
