package command

import (
	"../sets"
	"fmt"
	"gopkg.in/ldap.v2"
	"os"
	//	"strings"
	"github.com/urfave/cli"
	"log"
)

var (
	//ldapServer = "ds.trozlabs.local:389"
	ldapServer   = string(os.Getenv("LDAPServer"))
	ldapBind     = "CN=Administrator,CN=Users,DC=trozlabs,DC=local"
	ldapPassword = string(os.Getenv("LDAPPassword"))

	filterDN      = "(objectClass=*)"
	baseDN        = string(os.Getenv("LDAPBase"))
	loginUsername = string(os.Getenv("LDAPUser"))
	loginPassword = string(os.Getenv("LDAPPassword"))
)
var ld *ldap.Conn

func InitLDAP() (*ldap.Conn, error) {
	//tlsConfig := &tls.Config{InsecureSkipVerify: true}

	conn, err := ldap.Dial("tcp", ldapServer)

	if err != nil {
		return nil, fmt.Errorf("Failed to connect. %s", err)
	}

	if err := conn.Bind(ldapBind, ldapPassword); err != nil {
		return nil, fmt.Errorf("Failed to bind. %s", err)
	}

	ld = conn
	return ld, nil
}

func CmdGetAllAttr(c *cli.Context) {
	l := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(objectClass=*)",
		[]string{},
		nil,
	)
	sr, err := ld.Search(l)
	if err != nil {
		log.Fatal(err)
	}
	ns := sets.NewSet()
	for _, entry := range sr.Entries {
		//entry.PrettyPrint(1)
		//fmt.Println(entry.DN)
		for _, attr := range entry.Attributes {
			//fmt.Println(attr.Name)
			ns.Add(attr.Name)
		}

	}

	//var removeErrant = regexp.MustCompile(`[a-zA-Z: 0-9=]+`)

	//theSplit := strings.Split(strings.Join(removeErrant.FindAllString(fil, -1), ""), "=")

	//fmt.Println(strings.Join(theSplit, ", "))
	ns.PrintAll()
	//	return ns
}
func CmdGetAllDNs(c *cli.Context) {
	//conn = ld
	l := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(objectClass=*)",
		[]string{},
		nil,
	)
	//fmt.Println(c.Ldap())
	//fmt.Println(c)
	//conn := c.Conn
	sr, err := ld.Search(l)
	if err != nil {
		log.Fatal(err)
	}
	ns := sets.NewSet()
	for _, entry := range sr.Entries {
		//entry.PrettyPrint(1)
		ns.Add(entry.DN)
		fmt.Println(entry.DN)
	}

	//var removeErrant = regexp.MustCompile(`[a-zA-Z: 0-9=]+`)

	//theSplit := strings.Split(strings.Join(removeErrant.FindAllString(fil, -1), ""), "=")

	//fmt.Println(strings.Join(theSplit, ", "))
	ns.PrintAll()
}

func CmdGetAllThirds(c *cli.Context) {
	//conn = ld
	l := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(objectClass=*)",
		[]string{},
		nil,
	)
	//fmt.Println(c.Ldap())
	//fmt.Println(c)
	//conn := c.Conn
	sr, err := ld.Search(l)
	if err != nil {
		log.Fatal(err)
	}
	ns := sets.NewSet()
	for _, entry := range sr.Entries {
		//entry.PrettyPrint(1)
		ns.Add(entry.DN)
		//	fmt.Println(entry.DN)
	}

	ns.PrintThird()
}
