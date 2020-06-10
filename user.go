package goldap

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	guuid "github.com/google/uuid"
)

// AddUser function
func (l *LDAP) AddUser(name string, password string) error {
	uid := guuid.New().String()
	userDN := fmt.Sprintf("CN=%s,%s", name, l.GetUsersDN())

	addReq := ldap.NewAddRequest(userDN, []ldap.Control{})
	var attrs []ldap.Attribute
	attr := ldap.Attribute{
		Type: "objectClass",
		Vals: []string{"top", "person", "uidObject", "simpleSecurityObject"},
	}
	attrs = append(attrs, attr)

	attr = ldap.Attribute{
		Type: "cn",
		Vals: []string{name},
	}
	attrs = append(attrs, attr)

	attr = ldap.Attribute{
		Type: "sn",
		Vals: []string{name},
	}
	attrs = append(attrs, attr)

	attr = ldap.Attribute{
		Type: "uid",
		Vals: []string{uid},
	}
	attrs = append(attrs, attr)

	attr = ldap.Attribute{
		Type: "userPassword",
		Vals: []string{password},
	}
	attrs = append(attrs, attr)

	addReq.Attributes = attrs

	if err := l.connection.Add(addReq); err != nil {
		return err
	}
	return nil
}
