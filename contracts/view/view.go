package view

type View interface {
	Render(name string, bind any, layouts ...string) error
}
