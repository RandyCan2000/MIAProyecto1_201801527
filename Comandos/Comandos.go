package Comandos

import (
	Estruct "Proyecto1MIA/Estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"github.com/github.com/mitchellh/colorstring"
)

type Prueba struct {
	Numero  int
	Numero2 int
}

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
	var size2 int64 = 0
	if path == "" || size == "" || name == "" {
		colorstring.Println("[red]Faltan parametros")
		return false
	} else {
		path = removeCom(path)
	}
	if strings.ToUpper(strings.TrimSpace(unit)) == "K" {
		s, _ := strconv.Atoi(size)
		s2 := (s * 1024) - 1
		S3 := strconv.Itoa(s2)
		size2, _ = strconv.ParseInt(S3, 10, 64)
	} else if strings.ToUpper(strings.TrimSpace(unit)) == "M" || strings.ToUpper(strings.TrimSpace(unit)) == "" {
		s, _ := strconv.Atoi(size)
		s2 := (s * 1024 * 1024) - 1
		S3 := strconv.Itoa(s2)
		size2, _ = strconv.ParseInt(S3, 10, 64)
	} else {
		colorstring.Println("[red]La unidad debe ser K, M o en blanco")
		return false
	}
	err := os.MkdirAll(path+"/", 0775)
	if err != nil {
		colorstring.Println("[red]No se creo la carpeta correctamente")
		return false
	}
	colorstring.Println("[green]Carpeta encontrada")
	file, err2 := os.Create(name)
	defer file.Close()
	if err2 != nil {
		colorstring.Println("[red]No se creo el disco correctamente")
		return false
	}
	colorstring.Println("[green]Disco Creado con exito")
	var otro int8 = 0
	s := &otro
	fmt.Println(unsafe.Sizeof(otro))
	//Escribimos un 0 en el inicio del archivo.
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)
	WriteByte(file, binario.Bytes())
	println(size2)
	file.Seek(size2, 0) // segundo parametro: 0, 1, 2.     0 -> Inicio, 1-> desde donde esta el puntero, 2 -> Del fin para atras

	//Escribimos un 0 al final del archivo.
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s)
	WriteByte(file, binario2.Bytes())
	return true
}

func RMDISK(path string) bool {
	path = removeCom(path)
	var extencion []string = strings.Split(path, ".")
	if strings.ToUpper(extencion[len(extencion)-1]) != "DSK" {
		colorstring.Println("[red]El archivo no es extencion .dsk")
		return false
	}
	//TODO crear mensaje de confirmacion para eliminar la particion
	err := os.Remove(path)
	if err != nil {
		colorstring.Println("[red]El disco no se pudo borra con exito" + err.Error())
		return false
	}
	colorstring.Println("[green]Se removio el disco con exito")
	return true
}

func FDISK(path string, size string, unit string, tipe string, fit string, delete string, name string, add string) bool {
	file, err := os.Open(removeCom(path))
	defer file.Close()
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al abrir el Disco " + path)
		return false
	}
	mbr := Estruct.MBR{}
	A := ReadBytes(file, int(unsafe.Sizeof(mbr)))
	fmt.Println(A)
	buffer := bytes.NewBuffer(A)
	err = binary.Read(buffer, binary.BigEndian, &mbr)
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
		return false
	}
	fmt.Println(mbr)
	return true
}

func removeCom(path string) string {
	if strings.TrimSpace(string(path[0])) == "\"" {
		path2 := strings.TrimSpace(strings.Split(path, "\"")[1])
		path = path2
	}
	if strings.TrimSpace(string(path[0])) == "'" {
		path2 := strings.TrimSpace(strings.Split(path, "'")[1])
		path = path2
	}
	return path
}

func WriteByte(file *os.File, bytes []byte) {
	println(bytes)
	_, err := file.Write(bytes)
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al escribir en el archivo " + file.Name())
	}

}
func ReadBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number) //array de bytes
	_, err := file.Read(bytes)    // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
