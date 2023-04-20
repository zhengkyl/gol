package common

type Common struct {
	Width  int
	Height int
}

type CommonModel interface {
	SetSize(width, height int)
}
