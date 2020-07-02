package goldap

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// AddGroup function
func (l *LDAP) AddGroup(name string) error {
	groupDN := fmt.Sprintf("CN=%s,%s", name, l.baseDN)
	addReq := ldap.NewAddRequest(groupDN, []ldap.Control{})
	var attrs []ldap.Attribute
	attr := ldap.Attribute{
		Type: "member",
		Vals: []string{""},
	}
	attrs = append(attrs, attr)
	attr = ldap.Attribute{
		Type: "objectClass",
		Vals: []string{"top", "groupOfNames"},
	}
	attrs = append(attrs, attr)

	attr = ldap.Attribute{
		Type: "cn",
		Vals: []string{groupDN},
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
	userDN := fmt.Sprintf("cn=%s,%s", user, l.usersDN)
	groupDN := fmt.Sprintf("cn=%s,%s", group, l.baseDN)
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
	userDN := fmt.Sprintf("CN=%s,%s", user, l.usersDN)
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

func (l *LDAP) UserGroups(user string) ([]string, error) {
	userDN := fmt.Sprintf("cn=%s,%s", user, l.usersDN)
	searchReq := ldap.NewSearchRequest(
		l.baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=groupOfNames)(member=%s))", userDN),
		[]string{"cn"},
		[]ldap.Control{})

	result, err := l.connection.Search(searchReq)
	if err != nil {
		return nil, err
	}

	var groups []string
	for _, entry := range result.Entries {
		groups = append(groups, strings.Split(strings.Split(entry.GetAttributeValue("cn"), ",")[0], "=")[1])
	}

	return groups, nil
}

// MemberOf function
func (l *LDAP) MemberOf(group string) ([]string, error) {
	groupDN := fmt.Sprintf("CN=%s,%s", group, l.baseDN)
	searchReq := ldap.NewSearchRequest(
		l.baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=inetOrgPerson)(memberof=%s))", groupDN),
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
		fmt.Sprintf("(&(objectClass=groupOfNames)(member=CN=%s,%s)(cn=%s))", user, l.usersDN, group),
		[]string{"cn"},
		[]ldap.Control{})

	result, err := l.connection.Search(searchReq)
	if err != nil {
		return false, err
	}

	if len(result.Entries) > 0 {
		if result.Entries[0].GetAttributeValue("cn") == group {
			return true, nil
		}
	}

	return false, nil
}

func (l *LDAP) Groups() ([]string, error) {
	searchReq := ldap.NewSearchRequest(
		l.baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=groupOfNames)",
		[]string{"cn"},
		[]ldap.Control{})

	result, err := l.connection.Search(searchReq)
	if err != nil {
		return nil, err
	}

	var groups []string
	for _, entry := range result.Entries {
		groups = append(groups, strings.Split(strings.Split(entry.GetAttributeValue("cn"), ",")[0], "=")[1])
	}

	return groups, nil
}

func (l *LDAP) GroupExists(groupName string) (bool, error) {
	searchReq := ldap.NewSearchRequest(
		l.baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=groupOfNames)(cn=%s))", groupName),
		[]string{"cn"},
		[]ldap.Control{})

	result, err := l.connection.Search(searchReq)
	if err != nil {
		return false, err
	}

	if len(result.Entries) != 0 {
		return true, nil
	}
	return false, nil
}

// DelGroup function
func (l *LDAP) DelGroup(name string) error {
	groupDN := fmt.Sprintf("CN=%s,%s", name, l.baseDN)
	deleteGroup := ldap.NewDelRequest(groupDN, []ldap.Control{})
	err := l.connection.Del(deleteGroup)

	if err != nil {
		return err
	}

	return nil
}
