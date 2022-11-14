package console

import (
	"errors"
	"github.com/sujit-baniya/framework/contracts/console"
	"github.com/sujit-baniya/framework/contracts/console/command"
	"github.com/sujit-baniya/framework/support/file"
	"github.com/sujit-baniya/framework/support/str"
	"os"
	"strings"
)

type MakeCommand struct {
}

// Signature The name and signature of the console command.
func (receiver *MakeCommand) Signature() string {
	return "make:command"
}

// Description The console command description.
func (receiver *MakeCommand) Description() string {
	return "Create a new Artisan command"
}

// Extend The console command extend.
func (receiver *MakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			{
				Name:    "signature",
				Value:   "",
				Aliases: []string{"s"},
				Usage:   "signature of the command",
			},
			{
				Name:    "description",
				Value:   "",
				Aliases: []string{"d"},
				Usage:   "Command description",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *MakeCommand) Handle(ctx console.Context) error {
	name := ctx.Argument(0)
	signature := ctx.Option("signature")
	description := ctx.Option("description")
	if name == "" {
		return errors.New("Not enough arguments (missing: name) ")
	}
	stubContent := receiver.populateStub(receiver.getStub(), name)
	if signature != "" {
		stubContent = strings.ReplaceAll(stubContent, "command:name", signature)
	} else {
		stubContent = strings.ReplaceAll(stubContent, "command:name", str.ToCommandSignature(name))
	}
	if description != "" {
		stubContent = strings.ReplaceAll(stubContent, "Command description", description)
	}
	file.Create(receiver.getPath(name), stubContent)

	return nil
}

func (receiver *MakeCommand) getStub() string {
	return Stubs{}.Command()
}

// populateStub Populate the place-holders in the command stub.
func (receiver *MakeCommand) populateStub(stub string, name string) string {
	return strings.ReplaceAll(stub, "DummyCommand", str.Case2Camel(name))
}

// getPath Get the full path to the command.
func (receiver *MakeCommand) getPath(name string) string {
	pwd, _ := os.Getwd()

	return pwd + "/app/console/commands/" + str.Camel2Case(name) + ".go"
}
