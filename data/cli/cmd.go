package main

import (
	"duyog/loader"
	"flag"
	"fmt"
	"os"
)

type startCmd struct {
	name string

	configFlag *string

	configLoader loader.Config

	*flag.FlagSet
}

func (cmd startCmd) PrintHint() {
	fmt.Println(cmd.name)
	cmd.PrintDefaults()
}

func (cmd startCmd) Parse() bool {
	cmd.FlagSet.Parse(os.Args[2:])

	if cmd.Parsed() == false || *cmd.configFlag == "" {
		cmd.PrintDefaults()
		return false
	}

	if err := loader.LoadConfig(cmd.configLoader, *cmd.configFlag, &config); err != nil {
		fmt.Println("can not load file: " + *cmd.configFlag)
		fmt.Println(err)
		return false
	}

	return true
}

func newStartCmd() startCmd {
	defaultConfig := "./config.json"
	loader := loader.NewConfig()
	cmdName := "start"
	flagSet := flag.NewFlagSet(cmdName, flag.ExitOnError)

	return startCmd{
		name: cmdName,

		configFlag: flagSet.String("config", defaultConfig, "Path of the configuration named 'config.json'"),

		configLoader: loader,

		FlagSet: flagSet,
	}
}
