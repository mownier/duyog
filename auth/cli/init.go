package main

import (
	"duyog/cmd"
	"fmt"
	"os"
)

const tokenLen = 32
const defaultConfig = "./config.json"

var config option

func init() {
	hint := `try "help" for more details`

	start := newStartCmd()
	addClient := newAddClientCmd()

	if len(os.Args) < 2 {
		fmt.Println(hint)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "help":
		cmd.PrintHint(start, addClient)
		os.Exit(0)

	case "add-client":
		cmd.Parse(addClient)
		os.Exit(0)

	case "start":
		if cmd.Parse(start) == false {
			os.Exit(0)
		}

	default:
		fmt.Println(hint)
		os.Exit(0)
	}

	if err := config.valid(); err != nil {
		fmt.Print(err.Error())
		os.Exit(0)
	}
}
