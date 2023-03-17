package console

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/oarkflow/framework/contracts/console"
	"github.com/oarkflow/framework/contracts/console/command"
	"github.com/oarkflow/framework/support/file"
	"github.com/oarkflow/framework/support/str"

	"github.com/gookit/color"
)

type RequestMakeCommand struct {
}

// Signature The name and signature of the console command.
func (receiver *RequestMakeCommand) Signature() string {
	return "make:request"
}

// Description The console command description.
func (receiver *RequestMakeCommand) Description() string {
	return "Create a new request class"
}

// Extend The console command extend.
func (receiver *RequestMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			{
				Name:    "path",
				Value:   "/app/http/requests/",
				Aliases: []string{"p"},
				Usage:   "Path for request file",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *RequestMakeCommand) Handle(ctx console.Context) error {
	path := ctx.Option("path")
	name := ctx.Argument(0)
	if name == "" {
		return errors.New("Not enough arguments (missing: name) ")
	}
	if path == "" {
		path = "/app/http/requests/"
	}

	file.Create(receiver.getPath(path, name), receiver.populateStub(receiver.getStub(), name))
	color.Greenln("Request created successfully")

	return nil
}

func (receiver *RequestMakeCommand) getStub() string {
	return Stubs{}.Request()
}

// populateStub Populate the place-holders in the command stub.
func (receiver *RequestMakeCommand) populateStub(stub string, name string) string {
	stub = strings.ReplaceAll(stub, "DummyRequest", str.Case2Camel(name))
	stub = strings.ReplaceAll(stub, "DummyField", "Name string `form:\"name\" json:\"name\"`")

	return stub
}

// getPath Get the full path to the command.
func (receiver *RequestMakeCommand) getPath(path, name string) string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, path, str.Camel2Case(name)+".go")
}
