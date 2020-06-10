package goldap

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

// AddGroup function
func (l *LDAP) AddGroup(name string) error {
	groupDN := fmt.Sprintf("CN=%s,%s", name, l.baseDN)
	addReq := ldap.NewAddRequest(groupDN, []ldap.Control{})
	var attrs []ldap.Attribute
	attr := ldap.Attribute{
		Type: "objectClass",
		Vals: []string{"top", "groupOfUniqueNames"},
	}
	attrs = append(attrs, attr)

	attr = ldap.Attribute{
		Type: "cn",
		Vals: []string{name},
	}
	attrs = append(attrs, attr)

	attr = ldap.Attribute{
		Type: "uniqueMember",
		Vals: []string{""},
	}
	attrs = append(attrs, attr)

	addReq.Attributes = attrs

	if err := l.connection.Add(addReq); err != nil {
		return err
	}
	return nil
}

// AddUserToGroup function
func (l *LDAP) AddUserToGroup(user string, group string) error {
	userDN := fmt.Sprintf("CN=%s,%s", user, l.GetUsersDN())
	groupDN := fmt.Sprintf("CN=%s,%s", group, l.baseDN)
	modify := ldap.NewModifyRequest(userDN, []ldap.Control{})
	modify.Add("memberOf", []string{groupDN})

	err := l.connection.Modify(modify)
	if err != nil {
		return err
	}

	modify = ldap.NewModifyRequest(groupDN, []ldap.Control{})
	modify.Add("member", []string{userDN})

	err = l.connection.Modify(modify)
	if err != nil {
		return err
	}
	return nil
}

// DelUserGroup function
func (l *LDAP) DelUserGroup(user string, group string) error {
	userDN := fmt.Sprintf("CN=%s,%s", user, l.GetUsersDN())
	groupDN := fmt.Sprintf("CN=%s,%s", group, l.baseDN)
	modify := ldap.NewModifyRequest(userDN, []ldap.Control{})
	modify.Delete("memberOf", []string{groupDN})

	err := l.connection.Modify(modify)
	if err != nil {
		return err
	}

	modify = ldap.NewModifyRequest(groupDN, []ldap.Control{})
	modify.Delete("member", []string{userDN})

	err = l.connection.Modify(modify)
	if err != nil {
		return err
	}
	return nil
}

// MemberOf function
func (l *LDAP) MemberOf(group string) ([]string, error) {
	groupDN := fmt.Sprintf("CN=%s,%s", group, l.baseDN)
	searchReq := ldap.NewSearchRequest(
		l.baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=person)(memberof=%s))", groupDN),
		[]string{"cn"},
		[]ldap.Control{})

	result, err := l.connection.Search(searchReq)
	if err != nil {
		return nil, err
	}

	var users []string
	for _, entry := range result.Entries {
		users = append(users, entry.GetAttributeValue("cn"))
	}

	return users, nil
}

// UserMemberOf function
func (l *LDAP) UserMemberOf(user string, group string) (bool, error) {
	searchReq := ldap.NewSearchRequest(
		l.baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=groupOfNames)(member=CN=%s,CN=%s,%s))", user, group, l.baseDN),
		[]string{"cn"},
		[]ldap.Control{})

	result, err := l.connection.Search(searchReq)
	if err != nil {
		return false, err
	}

	if result.Entries[0].GetAttributeValue("cn") == group {
		return true, nil
	}

	return false, nil
}
