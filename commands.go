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
		Usage:  "Get all Distinguished Names",
		Action: command.CmdGetAllDNs,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "arp",
		Usage:  "Show ARP table",
		Action: command.CmdArp,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
