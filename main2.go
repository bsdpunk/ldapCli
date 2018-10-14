package main

import (
	"./shell2"
	"github.com/codegangsta/cli"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		shell2.Run()

	} else {

		app := cli.NewApp()
		//app.Name = Name
		//app.Version = Version
		app.Author = "bsdpunk"
		app.Email = ""
		app.Usage = ""
		//app.BeforeFunc = connectToLdap
		//	app.Flags = GlobalFlags
		//	app.Commands = Commands
		//	app.CommandNotFound = CommandNotFound

		//fmt.Println(os.Args)
		app.Run(os.Args)
	}
}
