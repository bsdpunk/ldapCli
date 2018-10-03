package main

import (
	"./command"
	"fmt"
	"github.com/urfave/cli"
	"os"
)

var GlobalFlags = []cli.Flag{}

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

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
