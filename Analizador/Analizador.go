//Paquete que contiene el parser de la consola
package Analizador

//Mkdisk -size->16 -path->"/home/randy/adios/OtroDisco" -NaMe->Disco4.dsk -uniT->M
//Mkdisk -size->l -path->"/home/mis discos/" -NaMe->Disco4.dsk -uniT->k
import (
	Code "Proyecto1MIA/Comandos"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/github.com/mitchellh/colorstring"
)

//fmt.Printf("\033[1;34m%s\033[0m", "Info")
//colorstring.Println("[blue]Hello [red]World!")

//Parser: Recibe una linea de comandos y la analiza
func Parser(Comando string) bool {
	var SourceSplit []string = Split(Comando, " ")
	//  {key}, {Value}
	i := 0
	if strings.ToUpper(strings.TrimSpace(Comando)) == "PAUSE" {
		colorstring.Println("[yellow]Presione cualquier tecla para salir del modo pausa")
		_, _, err := keyboard.GetSingleKey()
		if err != nil {
			return true
		}
		return true
	} else if strings.ToUpper(strings.TrimSpace(Comando)) == "CLEAR" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		return true
	} else if len(SourceSplit) <= 1 {
		return false
	} else if strings.ToUpper(SourceSplit[i]) == "EXEC" {
		var ExecSplit = strings.Split(SourceSplit[i+1], "->")
		if len(ExecSplit) <= 1 || ExecSplit[1] == "" {
			return false
		}
		if string(ExecSplit[1][0]) == "\"" {
			var PathSplit []string = strings.Split(ExecSplit[1], "\"")
			ExecSplit[1] = PathSplit[1]
		} else if string(ExecSplit[1][0]) == "'" {
			var PathSplit []string = strings.Split(ExecSplit[1], "'")
			ExecSplit[1] = PathSplit[1]
		}
		Code.EXEC(ExecSplit)
		return true
	} else if strings.ToUpper(SourceSplit[i]) == "MKDISK" {
		var path, size, name, unit string = "", "", "", ""
		for _, value := range SourceSplit {
			contain := strings.Contains(strings.ToUpper(value), "-PATH")
			if contain == true {
				path = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-SIZE")
			if contain == true {
				size = strings.TrimSpace(strings.Split(value, "->")[1])
				if isNumeric(size) == false {
					colorstring.Println("[red]size debe ser un numero positivo mayor a 0")
					return false
				}
			}
			contain = strings.Contains(strings.ToUpper(value), "-NAME")
			if contain == true {
				name = strings.TrimSpace(strings.Split(value, "->")[1])
				if Name(name) == false {
					colorstring.Println("[red]name debe contener solo caracteres [A-Z]|[a-z]|_|[0-9]")
					return false
				}
			}
			contain = strings.Contains(strings.ToUpper(value), "-UNIT")
			if contain == true {
				unit = strings.TrimSpace(strings.Split(value, "->")[1])
			}
		}
		//TODO borrar println
		colorstring.Println("[red]" + path + " " + size + " " + unit + " " + name)
		return Code.MKDISK(path, size, name, unit)
	} else {
		colorstring.Println("[red]El Script No existe")
		return false
	}
}  

//Para hacer Split parametro Spliter tomando en cuanta que entre comillas este spliter se ignora
func Split(Comando string, Spliter string) []string {
	var SoursesSplited []string
	var Source string = ""
	var Psplit bool = true
	for i := 0; i < len(Comando); i++ {
		if string(Comando[i]) == Spliter || len(Comando)-1 == i {
			if Psplit == true { //En true esta fuera de comillas
				if Source != "" {
					if len(Comando)-1 == i && string(Comando[i]) != Spliter {
						Source += string(Comando[i])
					}
					SoursesSplited = append(SoursesSplited, Source)
					Source = ""
				}
			} else {
				if len(Comando)-1 == i {
					if Source != "" {
						if string(Comando[i]) != Spliter {
							Source += string(Comando[i])
						}
						SoursesSplited = append(SoursesSplited, Source)
						Source = ""
					}
				}
				Source += string(Comando[i])
			}
		} else {
			if string(Comando[i]) == "\"" || string(Comando[i]) == "'" {
				if Psplit == true {
					Source += string(Comando[i]) //Aqui puedo quitar las comillas para retonrar ruta sin comillas
					Psplit = false
				} else {
					Source += string(Comando[i]) //Aqui puedo quitar las comillas para retonrar ruta sin comillas
					Psplit = true
				}
			} else {
				Source += string(Comando[i])
			}
		}
	}
	return SoursesSplited
}

func isNumeric(valor string) bool {
	Numero, err := strconv.Atoi(valor)
	if err != nil {
		return false
	}
	if Numero <= 0 {
		return false
	}
	return true
}

func Name(name string) bool {
	name2 := strings.Split(name, ".")
	if name2[1] != "dsk" {
		colorstring.Println("[red]name debe ser extension .dsk")
		return false
	}
	for _, char := range name2[0] {
		if char >= 48 && char <= 57 {
			//numeros
		} else if char >= 65 && char <= 90 {
			//alfabeto mayuscula
		} else if char >= 97 && char <= 122 {
			//alfabeto minuscula
		} else if char == 95 {
			//guion bajo
		} else {
			return false
		}
	}
	return true
}
