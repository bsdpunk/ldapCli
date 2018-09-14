package shell

import "../explore"
import "fmt"
import "os"
import "strings"
import "github.com/urfave/cli"
import "../command"
import "github.com/gobs/readline"

import "gopkg.in/ldap.v2"

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

var quit string = "quit"
var GlobalFlags = []cli.Flag{}

var found string = "no"

var conn *ldap.Conn
var err error

//conn, err := connect()

var words []string

var matches = make([]string, 0, len(words))

func AttemptedCompletion(text string, start, end int) []string {
	if start == 0 { // this is the command to match
		return readline.CompletionMatches(text, CompletionEntry)
	} else {
		return nil
	}
}

func CompletionEntry(prefix string, index int) string {
	if index == 0 {
		matches = matches[:0]

		for _, w := range words {
			if strings.HasPrefix(w, prefix) {
				matches = append(matches, w)
			}
		}
	}

	if index < len(matches) {
		return matches[index]
	} else {
		return ""
	}
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
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

var Commands = []cli.Command{
	{
		Name:   "GetAllDNs",
		Usage:  "Get All DNs",
		Action: command.CmdGetAllDNs,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "GetAllThirds",
		Usage:  "Get All DNs",
		Action: command.CmdGetAllThirds,
		Flags:  []cli.Flag{},
	},

	{
		Name:   "GetAllAttr",
		Usage:  "Get All Attributes",
		Action: command.CmdGetAllAttr,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "Search",
		Usage:  "Search LDAP",
		Action: command.CmdSearch,
		Flags:  []cli.Flag{},
	},

	{
		Name:   "arp",
		Usage:  "",
		Action: command.CmdArp,
		Flags:  []cli.Flag{},
	},
	//	{
	//		Name:   "GetAllDNs",
	//		Usage:  "",
	//		Action: command.CmdHeyo,
	//		Flags:  []cli.Flag{},
	//	},
}

func Run() {
	//	conn, err = connect()
	command.InitLDAP()

	for _, c := range Commands {
		words = append(words, c.Name)
	}
	words = append(words, "quit")
	words = append(words, "ls")

	prompt := "goldap> "
	matches = make([]string, 0, len(words))

L:
	for {
		found = "no"
		readline.SetCompletionEntryFunction(CompletionEntry)
		readline.SetAttemptedCompletionFunction(nil)
		result := readline.ReadLine(&prompt)
		if result == nil { // exit loop
			break L
		}

		input := *result
		input = strings.TrimSpace(input)
		if input == quit {
			os.Exit(1)
		} else if input == "ls" {
			fmt.Println(Commands)
		} else if input == "Explore" {
			prompt = "Explore> "
			ns := command.Explore()
			for _, newWord := range ns.ReturnThird() {
				words = append(words, newWord)
			}
			explore.Extui()

		} else {

			for _, c := range Commands {
				splitInput := strings.Split(input, " ")
				if c.HasName(splitInput[0]) {

					var command []string
					command = append(command, "")
					for _, i := range splitInput {

						command = append(command, i)
					}

					app := cli.NewApp()
					app.Author = "bsdpunk"
					app.Email = ""
					app.Usage = ""
					app.Name = splitInput[0]
					app.Version = "0.1.0"
					//app.Arg
					app.Flags = GlobalFlags
					app.Commands = Commands
					//app.CommandNotFound = CommandNotFound

					app.Run(command)
					found = "yes"

				}
			}
			if found == "no" {
				fmt.Println("Invalid Command")
			}
			readline.AddHistory(input)
		}

	}
}

func PrintSlice(slice []string) {
	fmt.Printf("Slice length = %d\r\n", len(slice))
	for i := 0; i < len(slice); i++ {
		fmt.Printf("[%d] := %s\r\n", i, slice[i])
	}
}
