package view

//go:generate mockery --name=View
type View interface {
	Render(name string, bind any, layouts ...string) error
}
