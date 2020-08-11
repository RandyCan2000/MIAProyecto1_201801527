package Comandos

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/github.com/mitchellh/colorstring"
)

func EXEC(ExecSplited []string) {
	if len(ExecSplited) >= 2 {
		if strings.ToUpper(ExecSplited[0]) == "-PATH" {
			if ExecSplited[1] == "" {
				colorstring.Println("[red]No hay una direccion")
				colorstring.Println("[red]Se Esperaba exec -path->\"path\"")
			} else {
				var Extension []string = strings.Split(ExecSplited[1], ".")
				if strings.TrimSpace(Extension[len(Extension)-1]) == "mia" {
					colorstring.Println("[green]" + ExecSplited[1])
					//Abrir archivos con comando y llamar a parser de analizador
				} else {
					colorstring.Println("[red]El archivo no es extension .mia")
				}
			}
		}
	} else {
		colorstring.Println("[red]Faltan Parametros en el comando Exec")
		colorstring.Println("[red]Se Esperaba exec -path->\"path\"")
	}
}

func MKDISK(path string, size string, name string, unit string) bool {
	if path == "" || size == "" || name == "" {
		colorstring.Println("[red]Faltan parametros")
		return false
	} else {
		if strings.TrimSpace(string(path[0])) == "\"" {
			path2 := strings.TrimSpace(strings.Split(path, "\"")[1])
			path = path2
		}
		if strings.TrimSpace(string(path[0])) == "'" {
			path2 := strings.TrimSpace(strings.Split(path, "'")[1])
			path = path2
		}

		//os.chown, (path, int(os.getenv('SUDO_UID')), int(os.getenv('SUDO_GID')))
		/*
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				errDir := os.MkdirAll(path, os.ModeDir)
				if errDir != nil {
					colorstring.Println("[red]No se creo el disco con exito")
				}

			}
		*/
	}

	app := "sudo"
	arg0 := "mkdir"
	arg1 := "-p"
	arg2 := path

	cmd := exec.Command(app, arg0, arg1, arg2)
	_, err := cmd.Output()
	//dd if=/dev/zero of=/tmp/archivo_grande bs=1024 count=1024
	arg0 = "dd"
	arg1 = "if=/dev/zero"
	arg2 = "of=" + path + "/" + name
	arg3 := ""
	if strings.ToUpper(unit) == "K" {
		value, _ := strconv.Atoi(size)
		arg3 = "bs=" + strconv.Itoa(value*1000)
	} else if strings.ToUpper(unit) == "M" {
		arg3 = "bs=" + size + "MB"
	} else {
		arg3 = "bs=" + size + "MB"
	}
	arg4 := "count=1"

	if err != nil {
		println("[red]La carpeta no se Creo: " + err.Error())
		return false
	}
	colorstring.Println("[green]Se creo con exito la carpeta")
	cmd = exec.Command(app, arg0, arg1, arg2, arg3, arg4)
	_, err = cmd.Output()
	if err != nil {
		println("[red]El archivo binario no se creo con exito: " + err.Error())
		return false
	}
	return true
}

/*
func WritePPMHeader(outf *os.File) {
	fheader := []byte("P6\n")
	binary.Write(outf, binary.LittleEndian, fheader)
	fheader2 := []byte("640 480\n")
	binary.Write(outf, binary.LittleEndian, fheader2)
	fheader3 := []byte("255\n")
	binary.Write(outf, binary.LittleEndian, fheader3)
}

func write() {
	outf, _ := os.Create("/home/randy/Escritorio/nuevo.dsk")
	println(outf.Name())
	WritePPMHeader(outf)
	var ii int8 = 0
	err := binary.Write(outf, binary.LittleEndian, ii)
	if err != nil {
		fmt.Println("err!", err)
	}
	outf.Close()
}

*/
