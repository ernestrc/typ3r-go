package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	typ3r "github.com/ernestrc/typ3r-go"
)

func main() {
	var err error
	var args map[string]interface{}
	var config *typ3r.Config
	var client typ3r.Client

	usage := `Typ3r Client

Usage:
  typ3r ls
  typ3r new
  typ3r -h | --help
  typ3r --version

Options:
  -h --help     Show this screen.`

	if args, err = docopt.Parse(usage, nil, true, "typ3r CLI", false); err != nil {
		goto exit
	}

	if config, err = typ3r.Load(); err != nil {
		goto exit
	}

	client = typ3r.Client{Config: config}

	if args["ls"].(bool) {
		err = ls(&client)
	} else if args["new"].(bool) {
		err = newNote(&client)
	} else {
		fmt.Println(usage)
	}

exit:
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
