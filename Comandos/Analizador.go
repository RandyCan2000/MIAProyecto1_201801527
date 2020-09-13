package Comandos

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/github.com/mitchellh/colorstring"
)

//fmt.Printf("\033[1;34m%s\033[0m", "Info")

//Parser Recibe una linea de comandos y la analiza
func Parser(Comando string) bool {
	var SourceSplit []string = Split(strings.TrimSpace(Comando), " ")
	//  {key}, {Value}
	i := 0
	if strings.ToUpper(strings.TrimSpace(Comando)) == "" {
		return true
	} else if strings.ToUpper(strings.TrimSpace(Comando))[0] == '#' {
		return true
	} else if strings.ToUpper(strings.TrimSpace(Comando)) == "MOUNT" {
		return MOUNT("", "")
	} else if strings.ToUpper(strings.TrimSpace(Comando)) == "PAUSE" {
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
	} else if strings.ToUpper(strings.TrimSpace(Comando)) == "LOGOUT" {
		return LOGOUT()
	} else if strings.ToUpper(strings.TrimSpace(Comando)) == "EXIT" {
		hecho := MensajeConfirmacion("¿Seguro desea terminar la ejecucion? [Y/N]: ", "Y")
		if hecho == true {
			colorstring.Println("[yellow]\tTermino la Ejecucion")
			os.Exit(0)
		}
		return true
	} else if strings.ToUpper(SourceSplit[i]) == "MOSTRARPRUEBA" {
		Mostrar()
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
		EXEC(ExecSplit)
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
		return MKDISK(removeCom(path), size, removeCom(name), unit)
	} else if strings.ToUpper(SourceSplit[i]) == "RMDISK" {
		var path string = ""
		for _, value := range SourceSplit {
			contain := strings.Contains(strings.ToUpper(value), "-PATH")
			if contain == true {
				path = strings.TrimSpace(strings.Split(value, "->")[1])
			}
		}
		return RMDISK(removeCom(path))
	} else if strings.ToUpper(SourceSplit[i]) == "FDISK" {
		var path, size, unit, tipe, fit, delete, name, add string = "", "", "", "", "", "", "", ""
		contain := false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-PATH")
			if contain == true {
				path = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-SIZE")
			if contain == true {
				size = strings.TrimSpace(strings.Split(value, "->")[1])
				if isNumeric(size) == false {
					colorstring.Println("[red]Size debe ser numero y mayor a cero")
					return false
				}
			}
			contain = strings.Contains(strings.ToUpper(value), "-UNIT")
			if contain == true {
				unit = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-TYPE")
			if contain == true {
				tipe = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-FIT")
			if contain == true {
				fit = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-DELETE")
			if contain == true {
				delete = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-NAME")
			if contain == true {
				name = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-ADD")
			if contain == true {
				add = strings.TrimSpace(strings.Split(value, "->")[1])
			}
		}
		return FDISK(removeCom(path), size, unit, tipe, fit, delete, removeCom(name), add)
	} else if strings.ToUpper(SourceSplit[i]) == "MOUNT" {
		var path, name string = "", ""
		for _, value := range SourceSplit {
			contain := strings.Contains(strings.ToUpper(value), "-PATH")
			if contain == true {
				path = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-NAME")
			if contain == true {
				name = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
		}
		return MOUNT(removeCom(strings.TrimSpace(path)), removeCom(strings.ToUpper(strings.TrimSpace(name))))
	} else if strings.ToUpper(SourceSplit[i]) == "UNMOUNT" {
		var id []string
		for _, value := range SourceSplit {
			contain := strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = append(id, removeCom(strings.TrimSpace(strings.Split(value, "->")[1])))
			}
		}
		if len(id) == 0 {
			return false
		}
		return UNMOUNT(id)
	} else if strings.ToUpper(SourceSplit[i]) == "MKFS" {
		var id, tipo, add, unit string = "", "", "", ""
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-TYPE")
			if contain == true {
				id = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-ADD")
			if contain == true {
				id = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-UNIT")
			if contain == true {
				id = strings.TrimSpace(strings.Split(value, "->")[1])
			}
		}
		return MKFS(id, tipo, unit, add)
	} else if strings.ToUpper(SourceSplit[i]) == "LOGIN" {
		var usr, pwd, id string = "", "", ""
		for _, value := range SourceSplit {
			contain := strings.Contains(strings.ToUpper(value), "-USR")
			if contain == true {
				usr = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-PWD")
			if contain == true {
				pwd = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
		}
		return LOGIN(usr, pwd, id)
	} else if strings.ToUpper(SourceSplit[i]) == "MKDIR" {
		var id, path string = "", ""
		var p bool = false
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-PATH")
			if contain == true {
				path = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-P")
			if contain == true {
				p = true
			}
		}
		hecho := MKDIR(id, removeCom(path), p)
		if hecho == true {
			colorstring.Println("[blue]\tArchivo creado con exito")
		}
		return hecho
	} else if strings.ToUpper(SourceSplit[i]) == "MKFILE" {
		var id, path, size, cont string = "", "", "", ""
		var p bool = false
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-PATH")
			if contain == true {
				path = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-SIZE")
			if contain == true {
				size = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-CONT")
			if contain == true {
				cont = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-P")
			if contain == true {
				p = true
			}
		}
		hecho := MKFILE(id, removeCom(path), p, size, cont,true)
		if hecho == true {
			colorstring.Println("[blue]\tArchivo creado con exito")
		}
		return hecho
	} else if strings.ToUpper(SourceSplit[i]) == "REP" {
		var nombre, path, id, ruta string = "", "", "", ""
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-NOMBRE")
			if contain == true {
				nombre = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-NAME")
			if contain == true {
				nombre = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-RUTA")
			if contain == true {
				ruta = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-PATH")
			if contain == true {
				path = strings.TrimSpace(strings.Split(value, "->")[1])
			}
		}
		REP(removeCom(id), removeCom(strings.TrimSpace(strings.ToUpper(nombre))), removeCom(strings.TrimSpace(removeCom(path))), removeCom(strings.TrimSpace(removeCom(ruta))))
		return true
	} else if strings.ToUpper(SourceSplit[i]) == "CAT" {
		var id string = ""
		var path []string
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-FILE")
			if contain == true {
				path = append(path, removeCom(strings.TrimSpace(strings.Split(value, "->")[1])))
			}
		}
		return CAT(path, id)
	} else if strings.ToUpper(SourceSplit[i]) == "RM" {
		var id, path string = "", ""
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-PATH")
			if contain == true {
				path = strings.TrimSpace(strings.Split(value, "->")[1])
			}
		}
		return RM(id, path)
	} else if strings.ToUpper(SourceSplit[i]) == "EDIT" {
		var id, path, size, cont string = "", "", "", ""
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-PATH")
			if contain == true {
				path = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-SIZE")
			if contain == true {
				size = strings.TrimSpace(strings.Split(value, "->")[1])
			}
			contain = strings.Contains(strings.ToUpper(value), "-CONT")
			if contain == true {
				cont = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
		}
		return MKFILE(id, removeCom(path), false, size, cont,true)
	} else if strings.ToUpper(SourceSplit[i]) == "MKGRP" {
		var id, name string = "", ""
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-NAME")
			if contain == true {
				name = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}

		}
		return MKGRP(id, removeCom(name))
	} else if strings.ToUpper(SourceSplit[i]) == "RMGRP" {
		var id, name string = "", ""
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-NAME")
			if contain == true {
				name = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}

		}
		return RMGRP(id, removeCom(name))
	} else if strings.ToUpper(SourceSplit[i]) == "MKUSR" {
		var id, user, password, grupo string = "", "", "", ""
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-USR")
			if contain == true {
				user = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-PWD")
			if contain == true {
				password = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-GRP")
			if contain == true {
				grupo = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}

		}
		return MKUSER(id, user, password, grupo)
	} else if strings.ToUpper(SourceSplit[i]) == "RMUSR" {
		var id, user string = "", ""
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-USR")
			if contain == true {
				user = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
		}
		return RMUSR(id, user)
	} else if strings.ToUpper(SourceSplit[i]) == "CP" {
		var id, path, dest string = "", "", ""
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-PATH")
			if contain == true {
				path = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-DEST")
			if contain == true {
				dest = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
		}
		return CP(id, removeCom(path), removeCom(dest))
	} else if strings.ToUpper(SourceSplit[i]) == "MV" {
		var id, path, idDestiny, dest string = "", "", "", ""
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-PATH")
			if contain == true {
				path = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-DEST")
			if contain == true {
				dest = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-IDDESTINY")
			if contain == true {
				idDestiny = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			} else {
				contain = strings.Contains(strings.ToUpper(value), "-ID")
				if contain == true {
					id = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
				}
			}
		}
		return MV(id, idDestiny, removeCom(path), removeCom(dest))
	} else if strings.ToUpper(SourceSplit[i]) == "FIND" {
		var id, path, nombre string = "", "", ""
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-PATH")
			if contain == true {
				path = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-NOMBRE")
			if contain == true {
				nombre = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-NAME")
			if contain == true {
				nombre = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
		}
		return FIND(id, removeCom(path), removeCom(nombre))
	} else if strings.ToUpper(SourceSplit[i]) == "LOSS" {
		var id string = ""
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}

		}
		return LOSS(id)
	} else if strings.ToUpper(SourceSplit[i]) == "RECOVERY" {
		var id string = ""
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}

		}
		return REC(id)
	} else if strings.ToUpper(SourceSplit[i]) == "REN" {
		var id, path, nombre string = "", "", ""
		var contain bool = false
		for _, value := range SourceSplit {
			contain = strings.Contains(strings.ToUpper(value), "-ID")
			if contain == true {
				id = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-PATH")
			if contain == true {
				path = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
			contain = strings.Contains(strings.ToUpper(value), "-NAME")
			if contain == true {
				nombre = removeCom(strings.TrimSpace(strings.Split(value, "->")[1]))
			}
		}
		return REN(id, removeCom(path), removeCom(nombre))
	} else {
		colorstring.Println("[red]El Script No existe")
		return true
	}
}

//Split parametro Spliter tomando en cuanta que entre comillas este spliter se ignora
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

func isNumeric(valor string) bool { //Verifica si es mayor a 0 y si es un numero retorna false si no cumple alguna de las condiciones
	Numero, err := strconv.Atoi(valor)
	if err != nil {
		return false
	}
	if Numero <= 0 {
		return false
	}
	return true
}

//Name valida que el name este correcto
func Name(name string) bool {
	name = removeCom(name)
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
		} else if char == ' ' {

		} else {
			return false
		}
	}
	return true
}

func removeCom(name string) string {
	if len(name) == 0 {
		return name
	} else {
		if strings.TrimSpace(string(name[0])) == "\"" {
			name2 := strings.TrimSpace(strings.Split(name, "\"")[1])
			name = name2
		}
		if strings.TrimSpace(string(name[0])) == "'" {
			name2 := strings.TrimSpace(strings.Split(name, "'")[1])
			name = name2
		}
	}
	return name
}
