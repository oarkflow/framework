package view

//go:generate mockery --name=Mail
type View interface {
	Render(name string, bind any, layouts ...string) error
}
