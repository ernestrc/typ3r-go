package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	typ3r "github.com/ernestrc/typ3r-go"
)

var (
	err    error
	args   map[string]interface{}
	config *typ3r.Config
	client typ3r.Client
	usage  = `Typ3r Client

Usage:
  typ3r
  typ3r new
  typ3r -h | --help
  typ3r --version

Options:
  -h --help     Show this screen.`
)

func main() {

	if args, err = docopt.Parse(usage, nil, true, "typ3r CLI", false); err != nil {
		goto exit
	}

	if config, err = typ3r.LoadConfig(); err != nil {
		goto exit
	}

	client = typ3r.Client{Config: config}

	if args["new"].(bool) {
		err = newNote(&client)
	} else {
		err = console(&client)
	}

exit:
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
