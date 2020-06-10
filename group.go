package goldap

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

// AddGroup function
func (l *LDAP) AddGroup(name string) error {
	gropDN = fmt.Sprintf("CN=%s,DC=digiops,DC=com", name)
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

	if err := l.Add(addReq); err != nil {
		return err
	}
}

// AddUserToGroup function
func (l *LDAP) AddUserToGroup(userDN string, groupDN string) error {
	modify = ldap.NewModifyRequest(userDN, []ldap.Control{})
	modify.Add("memberOf", []string{groupDN})

	err = l.Modify(modify)
	if err != nil {
		return err
	}

	modify := ldap.NewModifyRequest(groupDN, []ldap.Control{})
	modify.Add("member", []string{userDN})

	err = l.Modify(modify)
	if err != nil {
		return err
	}
}

// DelUserGroup function
func (l *LDAP) DelUserGroup(userDN string, groupDN string) error {
	modify = ldap.NewModifyRequest(userDN, []ldap.Control{})
	modify.Delete("memberOf", []string{groupDN})

	err = l.Modify(modify)
	if err != nil {
		return err
	}

	modify := ldap.NewModifyRequest(groupDN, []ldap.Control{})
	modify.Delete("member", []string{userDN})

	err = l.Modify(modify)
	if err != nil {
		return err
	}
}

// MemberOf function
func (l *LDAP) MemberOf(group string) ([]string, error) {
	searchReq := ldap.NewSearchRequest(
		"dc=digiops,dc=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=person)(memberof=CN=%s,DC=digiops,DC=com))", group),
		[]string{"cn"},
		[]ldap.Control{})

	result, err := l.Search(searchReq)
	if err != nil {
		return nil, err
	}

	var users []string
	for _, entry := range result.Entries {
		append(users, entry.GetAttributeValue("cn"))
	}

	return users, nil
}

// UserMemberOf function
func (l *LDAP) UserMemberOf(user string, group string) (bool, error) {
	searchReq := ldap.NewSearchRequest(
		"dc=digiops,dc=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=groupOfNames)(member=CN=%s,CN=folan,DC=digiops,DC=com))", group),
		[]string{"cn"},
		[]ldap.Control{})

	result, err := l.Search(searchReq)
	if err != nil {
		return nil, err
	}

	if result.Entries[0].GetAttributeValue("cn") == group {
		return true, nil
	}

	return false, nil
}
