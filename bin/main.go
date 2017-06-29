package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/ernestrc/typ3r"
)

func exit(o fmt.Stringer, err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(o.String())
}

func main() {
	var err error
	var config *typ3r.Config

	usage := `Typ3r Client

Usage:
  typ3r ls
  typ3r -h | --help
  typ3r --version

Options:
  -h --help     Show this screen.`

	arguments, err := docopt.Parse(usage, nil, true, "typ3r CLI", false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config, err = typ3r.Load()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := typ3r.TPClient{Config: config}

	if arguments["ls"].(bool) {
		exit(client.ListNotes())
	} else {
		fmt.Println(usage)
	}
}
