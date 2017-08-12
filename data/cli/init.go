package main

import (
	"duyog/cmd"
	"fmt"
	"os"
)

func init() {
	hint := `try "help" for more details`

	start := newStartCmd()

	if len(os.Args) < 2 {
		fmt.Println(hint)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "help":
		cmd.PrintHint(start)
		os.Exit(0)

	case "start":
		if start.Parse() == false {
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
