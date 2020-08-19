package main

//al descargar del repositorio instalar //go get -u github.com/eiannone/keyboard
//go get github.com/mitchellh/colorstring
import (
	Parser "Proyecto1MIA/Analizador"
	"bufio"
	"os"
	"os/exec"
	"strings"

	"github.com/github.com/mitchellh/colorstring"
)

const (
	darkblue  = "\033[1;34m%s\033[0m"
	lightblue = "\033[1;36m%s\033[0m"
	yellow    = "\033[1;33m%s\033[0m"
	red       = "\033[1;31m%s\033[0m"
	blue      = "\033[0;36m%s\033[0m"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var Script string = ""
	var LineNew [2]string
	for {
		Script = ""
		LineNew[0] = ""
		LineNew[1] = ""
		colorstring.Print("[green]Script: ")
		for {
			input, _ := reader.ReadByte()
			if string(input) == "\n" {
				if LineNew[0] == "/" && LineNew[1] == "*" {
					Script = strings.TrimSpace(strings.Split(Script, "/*")[0]) + " "
				} else {
					break
				}
			} else {
				Script += string(input)
			}
			LineNew[0] = LineNew[1]
			LineNew[1] = string(input)
			if LineNew[0] == "/" && LineNew[1] == "*" {
				Script = strings.TrimSpace(strings.Split(Script, "/*")[0]) + " "
			}
		}
		if Script == "" {
			print()
		} else {
			println(Script) //TODO Eliminar esto luego
			err := Parser.Parser(Script)
			if err == false {
				colorstring.Println("[red] Ocurrio un error al Ejecutar Script")
				colorstring.Println("[red] verifique que el script se encuentre escrito correctamente")
			}
		}
	}
}

func SuperUsuarios() {
	app := "sudo"
	arg0 := "su"
	cmd := exec.Command(app, arg0)
	err := cmd.Run()
	if err != nil {
		colorstring.Println("[red]No se pudo inicir sesion como super usuario")
	}
}
