package Estructuras

var (
	ABC         string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumeroLetra        = 0
)

type Mount struct {
	Path   string
	Letra  string
	Numero int64
}

type MountFisic struct {
	Id     string
	Path   string
	Name   string
	CopySB int64
}
