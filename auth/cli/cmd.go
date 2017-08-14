package main

import (
	"github.com/mownier/duyog/auth/rds"
	"github.com/mownier/duyog/auth/store"
	"github.com/mownier/duyog/generator"
	"github.com/mownier/duyog/loader"
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

type addClientCmd struct {
	name string

	nameFlag   *string
	roleFlag   *string
	emailFlag  *string
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

func (cmd addClientCmd) PrintHint() {
	fmt.Println(cmd.name)
	cmd.PrintDefaults()
}

func (cmd addClientCmd) Parse() bool {
	cmd.FlagSet.Parse(os.Args[2:])

	if cmd.Parsed() == false || *cmd.emailFlag == "" || *cmd.nameFlag == "" {
		cmd.PrintDefaults()
		return false
	}

	if err := loader.LoadConfig(cmd.configLoader, *cmd.configFlag, &config); err != nil {
		fmt.Println("can not load file: " + *cmd.configFlag)
		fmt.Println(err)
		return false
	}

	pool := newPool()
	keyGen := generator.XIDKey()
	apiGen := generator.APIToken(tokenLen)
	secretGen := generator.SecretToken(tokenLen)

	repo := rds.ClientRepo(keyGen, apiGen, secretGen, pool)

	client := store.Client{
		Email: *cmd.emailFlag,
		Name:  *cmd.nameFlag,
		Role:  *cmd.roleFlag,
	}

	client, err := store.CreateClient(repo, client)

	if err != nil {
		fmt.Println("can not create client:", err)
		return false
	}

	fmt.Println(toString(client))
	return true
}

func newStartCmd() startCmd {
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

func newAddClientCmd() addClientCmd {
	loader := loader.NewConfig()
	cmdName := "add-client"
	flagSet := flag.NewFlagSet(cmdName, flag.ExitOnError)

	return addClientCmd{
		name: cmdName,

		nameFlag:   flagSet.String("name", "", "Name of the client (Required)"),
		roleFlag:   flagSet.String("role", "admin", "Role of the client"),
		emailFlag:  flagSet.String("email", "", "Email of the client (Required)"),
		configFlag: flagSet.String("config", defaultConfig, "Path of the configuration named 'config.json'"),

		configLoader: loader,

		FlagSet: flagSet,
	}
}
