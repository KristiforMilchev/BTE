package sending

type Model struct {
	width  int
	height int
}

func New() *Model { return &Model{} }
