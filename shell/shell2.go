package shell2

import "strconv"

//import "../explore"
import "fmt"
import "os"
import "strings"
import "github.com/urfave/cli"
import "../command"
import "github.com/gobs/readline"
import "github.com/lobocv/itree/ctx"
import "github.com/nsf/termbox-go"
import "gopkg.in/ldap.v2"
import "math"
import "bytes"
import "errors"

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

type ScreenState int

const (
	Directory ScreenState = iota
	Help
)

type Screen struct {
	SearchString []rune
	CurrentDir   *ctx.Directory
	state        ScreenState
	captureInput bool

	highlightedColor termbox.Attribute
	filteredColor    termbox.Attribute
	directoryColor   termbox.Attribute
	fileColor        termbox.Attribute
}

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
	words = append(words, "Explore")

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
			//ns := command.Explore()
			//for _, newWord := range ns.ReturnThird() {
			//	words = append(words, newWord)
			//}
			//explore.Extui()

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

func (s *Screen) draw() {
	switch s.state {
	case Help:
		s.clearScreen()
		s.Print(0, 0, termbox.ColorWhite, termbox.ColorDefault, "itree - An interactive tree application for file system navigation.")
		s.Print(0, 2, termbox.ColorWhite, termbox.ColorDefault, "Calvin Lobo, 2018")
		s.Print(0, 3, termbox.ColorWhite, termbox.ColorDefault, "https://github.com/lobocv/itree")
		s.Print(0, 5, termbox.ColorWhite, termbox.ColorDefault, "Usage:")
		s.Print(0, 7, termbox.ColorWhite, termbox.ColorDefault, "h - Toggle hidden files and folders.")
		s.Print(0, 8, termbox.ColorWhite, termbox.ColorDefault, "e - Log2 skip up.")
		s.Print(0, 9, termbox.ColorWhite, termbox.ColorDefault, "d - Log2 skip down.")
		s.Print(0, 10, termbox.ColorWhite, termbox.ColorDefault, "c - Toggle position between first and last file.")
		s.Print(0, 11, termbox.ColorWhite, termbox.ColorDefault, "/ - Goes into input mode for file searching. Press ESC / CTRL+C to exit input mode.")
	case Directory:
		upperLevels, err := strconv.Atoi(os.Getenv("MaxUpperLevels"))
		if err != nil {
			upperLevels = 3
		}
		for {
			s.clearScreen()
			// Print the current path
			s.Print(0, 0, termbox.ColorRed, termbox.ColorDefault, s.CurrentDir.AbsPath)
			if s.captureInput {
				instruction := "Enter a search string:"
				s.Print(0, 1, termbox.ColorWhite, termbox.ColorDefault, instruction)
				s.Print(len(instruction)+2, 1, termbox.ColorWhite, termbox.ColorDefault, string(s.SearchString))
			}
			dirlist := s.getDirView(upperLevels)
			err := s.drawDirContents(0, 2, dirlist)
			if err == nil {
				break
			} else {
				upperLevels -= 1
			}
		}
	}

	termbox.Flush()
}

func (s *Screen) getDirView(upperLevels int) ctx.DirView {
	// Create a slice of the directory chain containing upperLevels number of parents
	dir := s.CurrentDir
	dirlist := make([]*ctx.Directory, 0, 1+upperLevels)
	//	ns := Explore()
	dirlist = append(dirlist, dir)

	//ns := command.Explore()
	//for _, newWord := range ns.ReturnThird() {
	//	words = append(words, newWord)
	//}
	next := dir.Parent
	for ii := 0; next != nil; ii++ {
		if ii >= upperLevels {
			break
		}
		dirlist = append([]*ctx.Directory{next}, dirlist...)
		next = next.Parent
	}
	return dirlist
}

func (s *Screen) clearScreen() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func (s *Screen) Print(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (s *Screen) drawDirContents(x0, y0 int, dirlist ctx.DirView) error {
	var levelOffsetX, levelOffsetY int // draw position offset
	var stretch int                    // Length of line connecting subdirectories
	var maxLineWidth int               // Length of longest item in the directory
	var scrollOffsety int              // Offset to scroll the visible directory text by
	var subDirSpacing = 2              // Spacing between subdirectories (on top of max item length)

	screenWidth, screenHeight := termbox.Size()

	levelOffsetX = x0
	levelOffsetY = y0

	// Determine the scrolling offset
	scrollOffsety = levelOffsetY
	for _, dir := range dirlist {
		scrollOffsety += dir.FileIdx
	}
	// If the selected item is off the screen then shift the entire view up in order
	// to make it visible.
	scrollOffsety -= screenHeight - levelOffsetY
	if scrollOffsety < 0 {
		scrollOffsety = 0
	} else {
		pagejump := float64(screenHeight) / 5
		scrollOffsety = int(math.Ceil(float64(scrollOffsety)/pagejump) * pagejump)
	}

	// Iterate through the directory list, drawing a tree structure
	for level, dir := range dirlist {
		maxLineWidth = 0

		for ii, f := range dir.Files {

			// Keep track of the longest length item in the directory
			filenameLen := len(f.Name())
			if filenameLen > maxLineWidth {
				maxLineWidth = filenameLen
			}

			// Determine the color of the currently printing directory item
			var color termbox.Attribute
			if dir.FileIdx == ii && level == len(dirlist)-1 {
				color = s.highlightedColor
			} else {
				if _, ok := dir.FilteredFiles[ii]; ok {
					color = s.filteredColor
				} else if f.IsDir() {
					color = s.directoryColor
				} else {
					color = s.fileColor
				}

			}

			// Start creating the line to be printed
			line := bytes.Buffer{}
			if ii == 0 {
				line.WriteString(strings.Repeat("─", stretch))
			}

			switch ii {
			case 0:
				if level > 0 {
					if len(dir.Files) < 2 {
						line.WriteString(strings.Repeat("─", subDirSpacing))
					} else {
						line.WriteString(strings.Repeat("─", subDirSpacing))
						line.WriteString("┬─")
					}
				} else {
					line.WriteString(strings.Repeat(" ", subDirSpacing))
					line.WriteString("├─")
				}
			case len(dir.Files) - 1:
				line.WriteString(strings.Repeat(" ", subDirSpacing))
				line.WriteString("└─")
			default:
				line.WriteString(strings.Repeat(" ", subDirSpacing))
				line.WriteString("├─")
			}

			// Create the item label, add / if it is a directory
			line.WriteString(f.Name())
			if f.IsDir() {
				line.WriteString("/")
			}

			// Calculate the draw position
			y := levelOffsetY + ii - scrollOffsety
			x := levelOffsetX
			if ii == 0 {
				// The first item is connected to the parent directory with a line
				// shift the position left to account for this line
				x -= stretch
			}
			if x+len(line.String()) > screenWidth && len(dirlist) > 1 {
				return errors.New("DisplayOverflow")
			}
			if y < y0 {
				y = y0
			}
			s.Print(x, y, color, termbox.ColorDefault, line.String())
		}

		// Determine the length of line we need to draw to connect to the next directory
		if len(dir.Files) > 0 {
			stretch = maxLineWidth - len(dir.Files[dir.FileIdx].Name())
		}

		// Shift the draw position in preparation for the next directory
		levelOffsetY += dir.FileIdx
		levelOffsetX += maxLineWidth + 2 + subDirSpacing

	}

	return nil
}
