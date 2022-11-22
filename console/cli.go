package console

import (
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"

	"github.com/sujit-baniya/framework/contracts/console"
	"github.com/sujit-baniya/framework/contracts/console/command"
	"github.com/sujit-baniya/framework/support"
)

type Cli struct {
	instance *cli.App
}

func NewCli(name ...string) console.Artisan {
	cliName := "Goravel Framework"
	if len(name) > 0 {
		cliName = name[0]
	}
	instance := cli.NewApp()
	instance.Name = cliName
	instance.Usage = support.Version
	instance.UsageText = "artisan [global options] command [options] [arguments...]"

	return &Cli{instance}
}

func (c *Cli) Register(commands []console.Command) {
	for _, item := range commands {
		item := item
		cliCommand := cli.Command{
			Name:  item.Signature(),
			Usage: item.Description(),
			Action: func(ctx *cli.Context) error {
				return item.Handle(&CliContext{ctx})
			},
		}

		cliCommand.Category = item.Extend().Category
		cliCommand.Flags = flagsToCliFlags(item.Extend().Flags)
		c.instance.Commands = append(c.instance.Commands, &cliCommand)
	}
}

func (c *Cli) Unregister(command string) {
	for idx, v := range c.instance.Commands {
		if v.Name == command {
			c.instance.Commands = append(c.instance.Commands[0:idx], c.instance.Commands[idx+1:]...)
			break
		}
	}
}

// Call Run an Artisan console command by name.
func (c *Cli) Call(command string) {
	c.Run(append([]string{os.Args[0], "artisan"}, strings.Split(command, " ")...), false)
}

// CallAndExit Run an Artisan console command by name and exit.
func (c *Cli) CallAndExit(command string) {
	c.Run(append([]string{os.Args[0], "artisan"}, strings.Split(command, " ")...), true)
}

// Run a command. Args come from os.Args.
func (c *Cli) Run(args []string, exitIfArtisan bool) {
	if len(args) >= 2 {
		if args[1] == "artisan" {
			if len(args) == 2 {
				args = append(args, "--help")
			}

			if args[2] != "-V" && args[2] != "--version" {
				cliArgs := append([]string{args[0]}, args[2:]...)
				if err := c.instance.Run(cliArgs); err != nil {
					panic(err.Error())
				}
			}

			printResult(args[2])

			if exitIfArtisan {
				os.Exit(0)
			}
		}
	}
}

func flagsToCliFlags(flags []command.Flag) []cli.Flag {
	var cliFlags []cli.Flag
	for _, flag := range flags {
		cliFlags = append(cliFlags, &cli.StringFlag{
			Name:     flag.Name,
			Aliases:  flag.Aliases,
			Usage:    flag.Usage,
			Required: flag.Required,
			Value:    flag.Value,
		})
	}

	return cliFlags
}

func printResult(command string) {
	switch command {
	case "make:command":
		color.Greenln("Console command created successfully")
	case "-V", "--version":
		color.Greenln("Goravel Framework " + support.Version)
	}
}
