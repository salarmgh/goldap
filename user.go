package goldap

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	guuid "github.com/google/uuid"
)

// AddUser function
func (l *LDAP) AddUser(name string, email string, password string) error {
	uid := guuid.New().String()
	userDN := fmt.Sprintf("CN=%s,%s", name, l.usersDN)

	addReq := ldap.NewAddRequest(userDN, []ldap.Control{})
	var attrs []ldap.Attribute
	attr := ldap.Attribute{
		Type: "objectClass",
		Vals: []string{"top", "inetOrgPerson"},
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

	attr = ldap.Attribute{
		Type: "mail",
		Vals: []string{email},
	}
	attrs = append(attrs, attr)

	addReq.Attributes = attrs

	if err := l.connection.Add(addReq); err != nil {
		return err
	}
	return nil
}

// Auth function
func (l *LDAP) Auth(loginUser string, loginPass string) (bool, error) {
	result, err := l.connection.Search(ldap.NewSearchRequest(
		l.baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=inetOrgPerson)(cn=%s))", loginUser),
		[]string{"cn"},
		nil,
	))
	fmt.Println(fmt.Sprintf("(&(objectClass=inetOrgPerson)(cn=%s))", loginUser))

	if err != nil {
		return false, fmt.Errorf("Failed to find user. %s", err)
	}

	if len(result.Entries) < 1 {
		return false, fmt.Errorf("User does not exist")
	}

	if len(result.Entries) > 1 {
		return false, fmt.Errorf("Too many entries returned")
	}

	conn, err := l.NewConn("ldap://ldap:389")
	if err != nil {
		return false, err
	}
	defer conn.Close()
	if err := conn.Bind(result.Entries[0].DN, loginPass); err != nil {
		return false, fmt.Errorf("Failed to auth. %s", err)
	}
	return true, nil
}
