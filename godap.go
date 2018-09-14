package main

import (
	//	"crypto/tls"
	//	"bufio"
	"fmt"
	"gopkg.in/ldap.v2"
	"os"
	"strings"
	//	"sync"
	"github.com/tiborvass/uniline"
)

type Prompt struct {
	ps1 string
}

func (p *Prompt) printPrompt() {
	fmt.Print(p.ps1)
}

var current = Prompt{ps1: "$ "}

const (
	ldapServer   = "ds.trozlabs.local:389"
	ldapBind     = "uid=bsdpunk,ou=People,dc=trozlabs,dc=local"
	ldapPassword = "TroCar##1999"

	filterDN      = "(objectclass=*)"
	baseDN        = "DC=trozlabs,DC=local"
	loginUsername = "bsdpunk"
	loginPassword = "TroCar##1999"
)

func runsh(args []string, conn *ldap.Conn) (err error) {

	//reader := bufio.NewReader(os.Stdin)

	scanner := uniline.DefaultScanner()
	for scanner.Scan(current.ps1) {
		line := scanner.Text()
		if len(line) > 0 {
			scanner.AddToHistory(line)

			platform := []rune(line)

			switch string(platform[0:3]) {
			case "hel":
				s := `All commands begin with a three letter prefix:
		qui or exi - for quit
		fil - for filter
		hel - for help`
				fmt.Println(s)
				current.printPrompt()

			case "fil":
				fmt.Printf("filter: ")
				//text, _ := reader.ReadString('\n')
				for scanner.Scan(current.ps1) {
					line := scanner.Text()
					if len(line) > 0 {
						scanner.AddToHistory(line)
						listFil(conn, line)

					}
				}
				if err := scanner.Err(); err != nil {
					panic(err)
				}
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
		}
	}

	return nil
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
		//fmt.Println(conn)
		return
	}()
	defer conn.Close()
	var input string
	var args []string
	for {
		fmt.Scanln(&input)
		runsh(args, conn)
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
	result, err := conn.Search(ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter(fil),
		[]string{"dn", "uid", "objectClass", "sn", "homeDirectory"},
		nil,
	))

	if err != nil {
		return fmt.Errorf("Failed to search users. %s", err)
	}
	result.PrettyPrint(0)
	for _, entry := range result.Entries {
		fmt.Println(entry)
	}
	return nil
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

func filter(needle string) string {
	res := strings.Replace(
		filterDN,
		"{username}",
		needle,
		-1,
	)

	return res
}
