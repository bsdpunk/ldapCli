package main

import (
	//	"crypto/tls"
	"bufio"
	"fmt"
	"gopkg.in/ldap.v2"
	"os"
	"strings"
	//	"sync"
	"./sets"
	"log"
	//	"regexp"
)

type Prompt struct {
	ps1 string
}

func (p *Prompt) printPrompt() {
	fmt.Print(p.ps1)
}

var current = Prompt{ps1: "$ "}

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

func runsh(command string, args []string, conn *ldap.Conn) (s string, err error) {

	reader := bufio.NewReader(os.Stdin)
	platform := []rune(command)

	switch string(platform[0:3]) {
	case "hel":
		s := `All commands begin with a three letter prefix:
		qui or exi - for quit
		fil - for filter
		hel - for help`
		fmt.Println(s)
		current.printPrompt()
	case "del":
		fmt.Printf("delete attribute of dn: ")
		text, _ := reader.ReadString('\n')
		delSh(conn, text)
		current.printPrompt()

	case "rep":
		fmt.Printf("replace attribute of dn: ")
		text, _ := reader.ReadString('\n')
		repSh(conn, text)
		current.printPrompt()
	case "add":
		fmt.Printf("add attribute of dn: ")
		text, _ := reader.ReadString('\n')
		addSh(conn, text)
		current.printPrompt()
	case "fil":
		fmt.Printf("filter: ")
		text, _ := reader.ReadString('\n')
		listFil(conn, text)
		current.printPrompt()
	case "get":
		var attributes *sets.Set = getAllAttr(conn)
		var objectClasses *sets.Set = getAllOC(conn)
		attributes.PrintAll()
		objectClasses.PrintAll()
		current.printPrompt()
	case "dns":
		var dns *sets.Set = getAllDN(conn)
		dns.PrintAll()
		current.printPrompt()

	case "exi":
		os.Exit(0)
	case "qui":
		os.Exit(0)
	case "cle":
		fmt.Print("\033[H\033[2J")
		current.printPrompt()
	default:
		fmt.Println("Unrecognized Command")
		current.printPrompt()
	}
	return command, nil
}

func main() {
	current.printPrompt()
	var conn *ldap.Conn
	var err error
	go func() {
		conn, err = connect()
		if err != nil {
			fmt.Errorf("thing", err)
			return
		}
		current = Prompt{ps1: "> "}
		//fmt.Print("\033[H\033[2J")
		current.printPrompt()

		//fmt.Println(conn)
		return
	}()
	defer conn.Close()
	var input string
	var args []string
	for {
		fmt.Scanln(&input)
		runsh(input, args, conn)
	}

}

func connect() (*ldap.Conn, error) {
	//tlsConfig := &tls.Config{InsecureSkipVerify: true}

	conn, err := ldap.Dial("tcp", ldapServer)

	if err != nil {
		return nil, fmt.Errorf("Failed to connect. %s", err)
	}

	if err := conn.Bind(ldapBind, ldapPassword); err != nil {
		return nil, fmt.Errorf("Failed to bind. %s", err)
	}

	return conn, nil
}

func list(conn *ldap.Conn) error {
	result, err := conn.Search(ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter("*"),
		[]string{"dn", "uid", "objectClass", "sn", "homeDirectory"},
		nil,
	))

	if err != nil {
		return fmt.Errorf("Failed to search users. %s", err)
	}

	for _, entry := range result.Entries {
		fmt.Printf(
			"%s: %s %s -- %v -- %v\n",
			entry.DN,
			entry.GetAttributeValue("uid"),
			entry.GetAttributeValue("sn"),
			entry.GetAttributeValue("objectClass"),
			entry.GetAttributeValue("homeDirectory"),
		)
	}

	return nil
}

func listFil(conn *ldap.Conn, fil string) error {
	l := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fil[:len(fil)-1],
		[]string{},
		nil,
	)
	sr, err := conn.Search(l)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range sr.Entries {
		entry.Print()
		fmt.Println("")
	}
	return nil
}
func getAllDN(conn *ldap.Conn) (s *sets.Set) {
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
	sr, err := conn.Search(l)
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
	return ns
}

func getAllAttr(conn *ldap.Conn) (s *sets.Set) {
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
	sr, err := conn.Search(l)
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
	return ns
}
func getAllOC(conn *ldap.Conn) (s *sets.Set) {
	l := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(objectClass=*)",
		[]string{"objectClass"},
		nil,
	)
	sr, err := conn.Search(l)
	if err != nil {
		log.Fatal(err)
	}
	ns := sets.NewSet()
	for _, entry := range sr.Entries {
		//entry.PrettyPrint(1)
		//fmt.Println(entry.DN)
		for _, val := range entry.GetAttributeValues("objectClass") {
			//fmt.Println(attr.Name)
			ns.Add(val)
		}

	}

	//var removeErrant = regexp.MustCompile(`[a-zA-Z: 0-9=]+`)

	//theSplit := strings.Split(strings.Join(removeErrant.FindAllString(fil, -1), ""), "=")

	//fmt.Println(strings.Join(theSplit, ", "))
	ns.PrintAll()
	return ns
}

func addSh(conn *ldap.Conn, dn string) error {

	reader := bufio.NewReader(os.Stdin)
	// DN := dn + baseDN
	modify := ldap.NewModifyRequest(dn)
	fmt.Printf("Add to dn " + dn + "(Attribute Name):")

	att, _ := reader.ReadString('\n')

	fmt.Printf("Add to dn " + dn + "(Value):")
	val, _ := reader.ReadString('\n')
	value := []string{val, ""}
	modify.Add(att[:len(att)-1], value[:len(value)-1])

	//modify.Replace("mail", []string{"user@example.org"})
	err := conn.Modify(modify)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
func delSh(conn *ldap.Conn, dn string) error {

	reader := bufio.NewReader(os.Stdin)
	// DN := dn + baseDN
	modify := ldap.NewModifyRequest(dn)
	fmt.Printf("Delete attribute on dn " + dn + "(Attribute Name):")

	att, _ := reader.ReadString('\n')

	fmt.Printf("Delete attribute on dn " + dn + "(Value):")
	val, _ := reader.ReadString('\n')
	value := []string{val, ""}
	modify.Delete(att[:len(att)-1], value[:len(value)-1])

	//modify.Replace("mail", []string{"user@example.org"})
	err := conn.Modify(modify)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func repSh(conn *ldap.Conn, dn string) error {

	reader := bufio.NewReader(os.Stdin)
	// DN := dn + baseDN
	modify := ldap.NewModifyRequest(dn)
	fmt.Printf("Replace attribute on dn " + dn + "(Attribute Name):")

	att, _ := reader.ReadString('\n')

	fmt.Printf("Replace attribute on dn " + dn + "(Value):")
	val, _ := reader.ReadString('\n')
	value := []string{val, ""}
	modify.Replace(att[:len(att)-1], value[:len(value)-1])

	//modify.Replace("mail", []string{"user@example.org"})
	err := conn.Modify(modify)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func auth(conn *ldap.Conn) error {
	result, err := conn.Search(ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter(loginUsername),
		[]string{"dn"},
		nil,
	))

	if err != nil {
		return fmt.Errorf("Failed to find user. %s", err)
	}

	if len(result.Entries) < 1 {
		return fmt.Errorf("User does not exist")
	}

	if len(result.Entries) > 1 {
		return fmt.Errorf("Too many entries returned")
	}

	if err := conn.Bind(result.Entries[0].DN, loginPassword); err != nil {
		fmt.Printf("Failed to auth. %s", err)
	} else {
		fmt.Printf("Authenticated successfuly!")
	}

	return nil
}

func filterTwo(needle string, fil string) string {
	res := strings.Replace(
		needle,
		"{username}",
		fil,
		-1,
	)
	fmt.Println(res)

	return res
}

func filter(needle string) string {
	res := strings.Replace(
		filterDN,
		"{username}",
		needle,
		-1,
	)
	fmt.Println(res)
	return res
}
