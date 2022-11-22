package grpc

import (
	"net"

	"google.golang.org/grpc"
)

type Application struct {
	Engine *grpc.Server
}

func (app *Application) Init() *Application {
	if app.Engine != nil {
		return app
	}
	app.Engine = grpc.NewServer()
	return app
}

func (app *Application) Server() *grpc.Server {
	return app.Engine
}

func (app *Application) SetServer(server *grpc.Server) {
	app.Engine = server
}

func (app *Application) Run(host string) error {
	listen, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	if err := app.Engine.Serve(listen); err != nil {
		return err
	}

	return nil
}
