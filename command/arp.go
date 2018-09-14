package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
)

func CmdArp(c *cli.Context) error {
	contents, _ := ioutil.ReadFile("/proc/net/arp")
	fmt.Println(string(contents))
	return nil
}
