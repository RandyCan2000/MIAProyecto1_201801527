package comandos

import (
	Estruct "Proyecto1MIA/Estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/github.com/mitchellh/colorstring"
)

var (
	//DiscosMontados Discos con letra asignada
	DiscosMontados []Estruct.Mount
	//ParticionesMontada Lista de particiones montadas
	ParticionesMontada []Estruct.MountFisic
)

//EXEC Comando script de MIA
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

//MKDISK Comando script de MIA
func MKDISK(path string, size string, name string, unit string) bool {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
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
	file, err2 := os.Create(path + "/" + name)
	defer file.Close()
	if err2 != nil {
		colorstring.Println("[red]No se creo el disco correctamente")
		return false
	}
	colorstring.Println("[green]Disco Creado con exito")
	var otro int8 = 0
	s := &otro
	//Cero al inicio
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)
	WriteByte(file, binario.Bytes())
	//cero al final
	file.Seek(size2, 0)
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s)
	WriteByte(file, binario2.Bytes())
	//Escribir estructura
	file.Seek(0, 0)

	//Partition Estruct
	prt := Estruct.Partition{Part_fit: 0}
	prt.Part_size = 0
	prt.Part_status = 0
	prt.Part_type = 0
	//MBR
	mbr := Estruct.MBR{Mbr_tamaño: size2 + 1}
	date := time.Now()
	fecha := strconv.Itoa(date.Day()) + "/" + date.Month().String() + "/" + strconv.Itoa(date.Year()) + "-H:" + strconv.Itoa(date.Hour()) + "-M:" + strconv.Itoa(date.Minute())
	copy(mbr.Mbr_fecha_creacion[:], fecha)
	mbr.Mbr_disk_signature = int64(r1.Intn(1000))

	//asignar particiones
	startP1 := unsafe.Sizeof(mbr.Mbr_tamaño) + unsafe.Sizeof(mbr.Mbr_fecha_creacion) + unsafe.Sizeof(mbr.Mbr_disk_signature)
	prt.Part_start = int64(startP1)
	mbr.Mbr_partition_1 = prt
	println(startP1)
	//part2
	startP2 := startP1 + unsafe.Sizeof(Estruct.Partition{})
	println(startP2)
	prt.Part_start = int64(startP2)
	mbr.Mbr_partition_2 = prt
	//part3
	startP3 := startP2 + unsafe.Sizeof(Estruct.Partition{})
	println(startP3)
	prt.Part_start = int64(startP3)
	mbr.Mbr_partition_3 = prt
	//part3
	startP4 := startP3 + unsafe.Sizeof(Estruct.Partition{})
	println(startP4)
	prt.Part_start = int64(startP4)
	mbr.Mbr_partition_4 = prt

	var WriteStrucBinary bytes.Buffer
	sMBR := &mbr
	binary.Write(&WriteStrucBinary, binary.BigEndian, sMBR)
	WriteByte(file, WriteStrucBinary.Bytes())

	return true
}

//RMDISK Comando script de MIA
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

//FDISK Comando script de MIA
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
	writeParticion(mbr)
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
		return false
	}

	if strings.ToUpper(strings.TrimSpace(unit)) == "B" || strings.ToUpper(strings.TrimSpace(unit)) == "K" || strings.ToUpper(strings.TrimSpace(unit)) == "M" || strings.ToUpper(strings.TrimSpace(unit)) == "" {
		if strings.ToUpper(strings.TrimSpace(unit)) == "" {
			unit = "M"
		}
	} else {
		colorstring.Println("[red]Valor de unit incorrecto debe ser B, K, M o vacio")
		return false
	}
	if strings.ToUpper(strings.TrimSpace(tipe)) == "P" || strings.ToUpper(strings.TrimSpace(tipe)) == "E" || strings.ToUpper(strings.TrimSpace(tipe)) == "L" || strings.ToUpper(strings.TrimSpace(tipe)) == "" {
		if strings.ToUpper(strings.TrimSpace(tipe)) == "" {
			tipe = "P"
		}
		if strings.ToUpper(strings.TrimSpace(tipe)) == "L" {
			if string(mbr.Mbr_partition_1.Part_type) != "E" &&
				string(mbr.Mbr_partition_2.Part_type) != "E" &&
				string(mbr.Mbr_partition_3.Part_type) != "E" &&
				string(mbr.Mbr_partition_4.Part_type) != "E" {
				colorstring.Println("[red]No existe ninguna particion Extendida para crear una particion logica")
				return false
			}
		}
	} else {
		colorstring.Println("[red]Valor de type incorrecto debe ser P, E, L o vacio")
		return false
	}

	if strings.TrimSpace(strings.ToUpper(fit)) == "B" || strings.TrimSpace(strings.ToUpper(fit)) == "F" || strings.TrimSpace(strings.ToUpper(fit)) == "W" || strings.TrimSpace(strings.ToUpper(fit)) == "" {
		if strings.TrimSpace(strings.ToUpper(fit)) == "" {
			fit = "W"
		}
	} else {
		colorstring.Println("[red]Valor de fit incorrecto debe ser BF, FF, WF o vacio")
		return false
	}

	//Creacion de particiones

	if size != "" && name != "" {
		var NAME [16]byte
		copy(NAME[:], strings.TrimSpace(strings.ToUpper(name)))
		if NAME == mbr.Mbr_partition_1.Part_name {
			colorstring.Println("[red]El nombre de la particion no puede repetirse")
			return false
		} else if NAME == mbr.Mbr_partition_2.Part_name {
			colorstring.Println("[red]El nombre de la particion no puede repetirse")
			return false
		} else if NAME == mbr.Mbr_partition_3.Part_name {
			colorstring.Println("[red]El nombre de la particion no puede repetirse")
			return false
		} else if NAME == mbr.Mbr_partition_4.Part_name {
			colorstring.Println("[red]El nombre de la particion no puede repetirse")
			return false
		}
		if strings.ToUpper(tipe) == "P" || strings.ToUpper(tipe) == "E" {
			if strings.ToUpper(tipe) == "E" {
				if string(mbr.Mbr_partition_1.Part_type) == "E" || string(mbr.Mbr_partition_2.Part_type) == "E" ||
					string(mbr.Mbr_partition_3.Part_type) == "E" || string(mbr.Mbr_partition_4.Part_type) == "E" {
					colorstring.Println("[red]Ya existe una particion Extendida")
					return false
				}
			}
			var Size int = 0
			if strings.ToUpper(unit) == "B" {
				SizeAux, _ := strconv.Atoi(size)
				Size = SizeAux
			} else if strings.ToUpper(unit) == "K" {
				SizeAux, _ := strconv.Atoi(size)
				Size = SizeAux * 1024
			} else if strings.ToUpper(unit) == "M" {
				SizeAux, _ := strconv.Atoi(size)
				Size = SizeAux * 1024 * 1024
			} else {
				colorstring.Println("[red]Unit indefinido")
				return false
			}
			sizeP1 := int64(mbr.Mbr_partition_1.Part_size)
			sizeP2 := int64(mbr.Mbr_partition_2.Part_size)
			sizeP3 := int64(mbr.Mbr_partition_3.Part_size)
			sizeP4 := int64(mbr.Mbr_partition_4.Part_size)
			sizeTotal := int64(Size) + sizeP1 + sizeP2 + sizeP3 + sizeP4
			if sizeTotal > int64(mbr.Mbr_tamaño) && strings.ToUpper(tipe)[0] != 'L' {
				colorstring.Println("[red]El tamanio de la particion sobrepasa la del disco")
				return false
			}
			if string(mbr.Mbr_partition_1.Part_status) != "A" {
				copy(mbr.Mbr_partition_1.Part_name[:], strings.ToUpper(name))
				mbr.Mbr_partition_1.Part_fit = strings.ToUpper(fit)[0]
				mbr.Mbr_partition_1.Part_size = int64(Size)
				mbr.Mbr_partition_1.Part_status = 'A'
				mbr.Mbr_partition_1.Part_type = strings.ToUpper(tipe)[0]
				mbr.Mbr_partition_1.Part_start = int64(unsafe.Sizeof(mbr)) + 20
				createPartition(file, mbr)
				if mbr.Mbr_partition_1.Part_type == 'E' {
					ebr := Estruct.EBR{Part_next: -1}
					ebr.Part_start = int64(unsafe.Sizeof(mbr)) + 20
					createPartitionLogic(file, ebr, int64(unsafe.Sizeof(mbr))+20)
				}
				return true
			} else if string(mbr.Mbr_partition_2.Part_status) != "A" {
				copy(mbr.Mbr_partition_2.Part_name[:], strings.ToUpper(name))
				mbr.Mbr_partition_2.Part_fit = strings.ToUpper(fit)[0]
				mbr.Mbr_partition_2.Part_size = int64(Size)
				mbr.Mbr_partition_2.Part_status = 'A'
				mbr.Mbr_partition_2.Part_type = strings.ToUpper(tipe)[0]
				mbr.Mbr_partition_2.Part_start = int64(unsafe.Sizeof(mbr)) + 20
				createPartition(file, mbr)
				if mbr.Mbr_partition_2.Part_type == 'E' {
					ebr := Estruct.EBR{Part_next: -1}
					ebr.Part_start = int64(unsafe.Sizeof(mbr)) + 20
					createPartitionLogic(file, ebr, int64(unsafe.Sizeof(mbr))+20)
				}
				return true
			} else if string(mbr.Mbr_partition_3.Part_status) != "A" {
				copy(mbr.Mbr_partition_3.Part_name[:], strings.ToUpper(name))
				mbr.Mbr_partition_3.Part_fit = strings.ToUpper(fit)[0]
				mbr.Mbr_partition_3.Part_size = int64(Size)
				mbr.Mbr_partition_3.Part_status = 'A'
				mbr.Mbr_partition_3.Part_type = strings.ToUpper(tipe)[0]
				mbr.Mbr_partition_3.Part_start = int64(unsafe.Sizeof(mbr)) + 20
				createPartition(file, mbr)
				if mbr.Mbr_partition_3.Part_type == 'E' {
					ebr := Estruct.EBR{Part_next: -1}
					ebr.Part_start = int64(unsafe.Sizeof(mbr)) + 20
					createPartitionLogic(file, ebr, int64(unsafe.Sizeof(mbr))+20)
				}
				return true
			} else if string(mbr.Mbr_partition_4.Part_status) != "A" {
				copy(mbr.Mbr_partition_4.Part_name[:], strings.ToUpper(name))
				mbr.Mbr_partition_4.Part_fit = strings.ToUpper(fit)[0]
				mbr.Mbr_partition_4.Part_size = int64(Size)
				mbr.Mbr_partition_4.Part_status = 'A'
				mbr.Mbr_partition_4.Part_type = strings.ToUpper(tipe)[0]
				mbr.Mbr_partition_4.Part_start = int64(unsafe.Sizeof(mbr)) + 20
				createPartition(file, mbr)
				if mbr.Mbr_partition_4.Part_type == 'E' {
					ebr := Estruct.EBR{Part_next: -1}
					ebr.Part_start = int64(unsafe.Sizeof(mbr)) + 20
					createPartitionLogic(file, ebr, int64(unsafe.Sizeof(mbr))+20)
				}
				return true
			} else {
				colorstring.Println("[red]No se encontraron particiones vacias")
				return false
			}
		} else if strings.ToUpper(tipe)[0] == 'L' {
			var Size int = 0
			if strings.ToUpper(unit) == "B" {
				SizeAux, _ := strconv.Atoi(size)
				Size = SizeAux
			} else if strings.ToUpper(unit) == "K" {
				SizeAux, _ := strconv.Atoi(size)
				Size = SizeAux * 1024
			} else if strings.ToUpper(unit) == "M" {
				SizeAux, _ := strconv.Atoi(size)
				Size = SizeAux * 1024 * 1024
			} else {
				colorstring.Println("[red]Unit indefinido")
				return false
			}
			var extSize int64 = 0
			if mbr.Mbr_partition_1.Part_type == 'E' {
				file.Seek(mbr.Mbr_partition_1.Part_start, 0)
				extSize = mbr.Mbr_partition_1.Part_size
			} else if mbr.Mbr_partition_2.Part_type == 'E' {
				file.Seek(mbr.Mbr_partition_2.Part_start, 0)
				extSize = mbr.Mbr_partition_2.Part_size
			} else if mbr.Mbr_partition_3.Part_type == 'E' {
				file.Seek(mbr.Mbr_partition_3.Part_start, 0)
				extSize = mbr.Mbr_partition_3.Part_size
			} else if mbr.Mbr_partition_4.Part_type == 'E' {
				file.Seek(mbr.Mbr_partition_4.Part_start, 0)
				extSize = mbr.Mbr_partition_4.Part_size
			} else {
				colorstring.Println("[red]No hay particion Extendida")
				return false
			}
			ebr := Estruct.EBR{}
			A := ReadBytes(file, int(unsafe.Sizeof(ebr)))
			buffer := bytes.NewBuffer(A)
			err = binary.Read(buffer, binary.BigEndian, &ebr)
			if err != nil {
				colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
				return false
			}
			ebrInsert := Estruct.EBR{Part_next: -1}
			ebrInsert.Part_fit = strings.ToUpper(fit)[0]
			copy(ebrInsert.Part_name[:], strings.ToUpper(name))
			ebrInsert.Part_size = int64(Size)
			ebrInsert.Part_status = 'A'
			RecursivoFindLogic(file, ebr, ebrInsert, extSize, ebrInsert.Part_size+ebr.Part_size)
			showLogicRec(file, ebr)
			return true
		}
	}

	if strings.ToUpper(strings.TrimSpace(delete)) == "FAST" || strings.ToUpper(strings.TrimSpace(delete)) == "FULL" {
		var partStart int64
		if mbr.Mbr_partition_1.Part_type == 'E' {
			partStart = mbr.Mbr_partition_1.Part_start
		} else if mbr.Mbr_partition_2.Part_type == 'E' {
			partStart = mbr.Mbr_partition_2.Part_start
		} else if mbr.Mbr_partition_3.Part_type == 'E' {
			partStart = mbr.Mbr_partition_3.Part_start
		} else if mbr.Mbr_partition_4.Part_type == 'E' {
			partStart = mbr.Mbr_partition_4.Part_start
		}
		var NAME [16]byte
		copy(NAME[:], strings.ToUpper(name))
		if mbr.Mbr_partition_1.Part_name == NAME {
			if strings.ToUpper(strings.TrimSpace(delete)) == "FAST" {
				mbr.Mbr_partition_1.Part_fit = ' '
				mbr.Mbr_partition_1.Part_size = 0
				mbr.Mbr_partition_1.Part_start = 0
				mbr.Mbr_partition_1.Part_status = ' '
				mbr.Mbr_partition_1.Part_type = ' '
				copy(mbr.Mbr_partition_1.Part_name[:], "                ")
			} else {
				mbr.Mbr_partition_1 = Estruct.Partition{}
			}
			createPartition(file, mbr)
			colorstring.Println("[green]Eliminado con exito")
			return true
		} else if mbr.Mbr_partition_2.Part_name == NAME {
			if strings.ToUpper(strings.TrimSpace(delete)) == "FAST" {
				mbr.Mbr_partition_2.Part_fit = ' '
				mbr.Mbr_partition_2.Part_size = 0
				mbr.Mbr_partition_2.Part_start = 0
				mbr.Mbr_partition_2.Part_status = ' '
				mbr.Mbr_partition_2.Part_type = ' '
				copy(mbr.Mbr_partition_2.Part_name[:], "                ")
			} else {
				mbr.Mbr_partition_2 = Estruct.Partition{}
			}
			createPartition(file, mbr)
			colorstring.Println("[green]Eliminado con exito")
			return true
		} else if mbr.Mbr_partition_3.Part_name == NAME {
			if strings.ToUpper(strings.TrimSpace(delete)) == "FAST" {
				mbr.Mbr_partition_3.Part_fit = ' '
				mbr.Mbr_partition_3.Part_size = 0
				mbr.Mbr_partition_3.Part_start = 0
				mbr.Mbr_partition_3.Part_status = ' '
				mbr.Mbr_partition_3.Part_type = ' '
				copy(mbr.Mbr_partition_3.Part_name[:], "                ")
			} else {
				mbr.Mbr_partition_3 = Estruct.Partition{}
			}
			createPartition(file, mbr)
			colorstring.Println("[green]Eliminado con exito")
			return true
		} else if mbr.Mbr_partition_4.Part_name == NAME {
			if strings.ToUpper(strings.TrimSpace(delete)) == "FAST" {
				mbr.Mbr_partition_4.Part_fit = ' '
				mbr.Mbr_partition_4.Part_size = 0
				mbr.Mbr_partition_4.Part_start = 0
				mbr.Mbr_partition_4.Part_status = ' '
				mbr.Mbr_partition_4.Part_type = ' '
				copy(mbr.Mbr_partition_4.Part_name[:], "                ")
			} else {
				mbr.Mbr_partition_4 = Estruct.Partition{}
			}
			createPartition(file, mbr)
			colorstring.Println("[green]Eliminado con exito")
			return true
		}
		if partStart != 0 {
			file.Seek(partStart, 0)
			ebr := Estruct.EBR{}
			A := ReadBytes(file, int(unsafe.Sizeof(ebr)))
			buffer := bytes.NewBuffer(A)
			err = binary.Read(buffer, binary.BigEndian, &ebr)
			if err != nil {
				colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
				return false
			}
			valor := RecDeletePartitionLogic(file, ebr, NAME, strings.ToUpper(delete))
			if valor != 0 {
				if mbr.Mbr_partition_1.Part_type == 'E' {
					mbr.Mbr_partition_1.Part_start = valor
				} else if mbr.Mbr_partition_2.Part_type == 'E' {
					mbr.Mbr_partition_2.Part_start = valor
				} else if mbr.Mbr_partition_3.Part_type == 'E' {
					mbr.Mbr_partition_3.Part_start = valor
				} else if mbr.Mbr_partition_4.Part_type == 'E' {
					mbr.Mbr_partition_4.Part_start = valor
				}
				createPartition(file, mbr)
			}
			colorstring.Println("[green]Eliminado")
			return true
		} else {
			colorstring.Println("[red]No se econtro ninguna particion con el nombre indicado")
			return false
		}
	}
	Add, erradd := strconv.Atoi(add)
	if erradd != nil {
		colorstring.Println("[red]Add debe ser un numero")
		return false
	} else {
		var NAME [16]byte
		copy(NAME[:], name)
		if Add < 0 {
			if mbr.Mbr_partition_1.Part_name == NAME {
				if mbr.Mbr_partition_1.Part_size-int64(Add) <= 0 {
					colorstring.Println("[red]No se puede reducir a menos de 0 el size de la particion")
				} else {
					newsize := mbr.Mbr_partition_1.Part_size - int64(Add)
					mbr.Mbr_partition_1.Part_size = newsize
					createPartition(file, mbr)
				}
			} else if mbr.Mbr_partition_2.Part_name == NAME {
				if mbr.Mbr_partition_2.Part_size-int64(Add) <= 0 {
					colorstring.Println("[red]No se puede reducir a menos de 0 el size de la particion")
				} else {
					newsize := mbr.Mbr_partition_2.Part_size - int64(Add)
					mbr.Mbr_partition_2.Part_size = newsize
					createPartition(file, mbr)
				}
			} else if mbr.Mbr_partition_3.Part_name == NAME {
				if mbr.Mbr_partition_3.Part_size-int64(Add) <= 0 {
					colorstring.Println("[red]No se puede reducir a menos de 0 el size de la particion")
				} else {
					newsize := mbr.Mbr_partition_3.Part_size - int64(Add)
					mbr.Mbr_partition_3.Part_size = newsize
					createPartition(file, mbr)
				}
			} else if mbr.Mbr_partition_4.Part_name == NAME {
				if mbr.Mbr_partition_4.Part_size-int64(Add) <= 0 {
					colorstring.Println("[red]No se puede reducir a menos de 0 el size de la particion")
				} else {
					newsize := mbr.Mbr_partition_4.Part_size - int64(Add)
					mbr.Mbr_partition_4.Part_size = newsize
					createPartition(file, mbr)
				}
			}
			return true
		} else if Add > 0 {
			newsize := mbr.Mbr_partition_1.Part_size + mbr.Mbr_partition_2.Part_size + mbr.Mbr_partition_3.Part_size + mbr.Mbr_partition_4.Part_size + int64(Add)
			if mbr.Mbr_partition_1.Part_name == NAME {
				if mbr.Mbr_tamaño >= newsize {
					newsize = mbr.Mbr_partition_1.Part_size + int64(Add)
					mbr.Mbr_partition_1.Part_size = newsize
					createPartition(file, mbr)
				}
			} else if mbr.Mbr_partition_2.Part_name == NAME {
				if mbr.Mbr_tamaño >= newsize {
					newsize = mbr.Mbr_partition_2.Part_size + int64(Add)
					mbr.Mbr_partition_2.Part_size = newsize
					createPartition(file, mbr)
				}
			} else if mbr.Mbr_partition_3.Part_name == NAME {
				if mbr.Mbr_tamaño >= newsize {
					newsize = mbr.Mbr_partition_3.Part_size + int64(Add)
					mbr.Mbr_partition_3.Part_size = newsize
					createPartition(file, mbr)
				}
			} else if mbr.Mbr_partition_4.Part_name == NAME {
				if mbr.Mbr_tamaño >= newsize {
					newsize = mbr.Mbr_partition_4.Part_size + int64(Add)
					mbr.Mbr_partition_4.Part_size = newsize
					createPartition(file, mbr)
				}
			}
			return true
		} else {
			colorstring.Println("[red]Add debe ser mayor o menor a 0")
			return false
		}
	}
	return false
}

//MOUNT Comando script de MIA
func MOUNT(path string, name string) bool {
	println(Estruct.NumeroLetra)
	if path == "" && name == "" {
		for _, Value := range ParticionesMontada {
			if Value.Name != "" {
				colorstring.Println("[green]\tID: " + Value.Id + " PATH: " + Value.Path + " NAME: " + Value.Name)
			}
		}
		return true
	} else {
		file, err := os.Open(removeCom(path))
		defer file.Close()
		if err != nil {
			colorstring.Println("[red]Ocurrio un error al abrir el Disco " + path)
			return false
		}
		mbr := Estruct.MBR{}
		A := ReadBytes(file, int(unsafe.Sizeof(mbr)))
		buffer := bytes.NewBuffer(A)
		err = binary.Read(buffer, binary.BigEndian, &mbr)
		if err != nil {
			colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
			return false
		}
		for i, Disco := range DiscosMontados {
			if path == Disco.Path {
				for _, Particiones := range ParticionesMontada {
					if Particiones.Name == name && path == Particiones.Path {
						colorstring.Println("[red]La particion ya se encuentra montada")
						return false
					}
				}
				var NAME [16]byte
				copy(NAME[:], name)
				if mbr.Mbr_partition_1.Part_name == NAME && mbr.Mbr_partition_1.Part_type != 'E' {
					PartMontNew := Estruct.MountFisic{Path: path}
					PartMontNew.Name = name
					PartMontNew.Id = "vd" + Disco.Letra + strconv.Itoa(int(Disco.Numero))
					DiscosMontados[i].Numero = Disco.Numero + 1
					ParticionesMontada = append(ParticionesMontada, PartMontNew)
					colorstring.Println("[green]Se monto con exito")
					return true
				} else if mbr.Mbr_partition_2.Part_name == NAME && mbr.Mbr_partition_2.Part_type != 'E' {
					PartMontNew := Estruct.MountFisic{Path: path}
					PartMontNew.Name = name
					PartMontNew.Id = "vd" + Disco.Letra + strconv.Itoa(int(Disco.Numero))
					DiscosMontados[i].Numero = Disco.Numero + 1
					ParticionesMontada = append(ParticionesMontada, PartMontNew)
					println(len(DiscosMontados))
					println(len(ParticionesMontada))
					colorstring.Println("[green]Se monto con exito")
					return true
				} else if mbr.Mbr_partition_3.Part_name == NAME && mbr.Mbr_partition_3.Part_type != 'E' {
					PartMontNew := Estruct.MountFisic{Path: path}
					PartMontNew.Name = name
					PartMontNew.Id = "vd" + Disco.Letra + strconv.Itoa(int(Disco.Numero))
					DiscosMontados[i].Numero = Disco.Numero + 1
					ParticionesMontada = append(ParticionesMontada, PartMontNew)
					colorstring.Println("[green]Se monto con exito")
					return true
				} else if mbr.Mbr_partition_4.Part_name == NAME && mbr.Mbr_partition_4.Part_type != 'E' {
					PartMontNew := Estruct.MountFisic{Path: path}
					PartMontNew.Name = name
					PartMontNew.Id = "vd" + Disco.Letra + strconv.Itoa(int(Disco.Numero))
					DiscosMontados[i].Numero = Disco.Numero + 1
					ParticionesMontada = append(ParticionesMontada, PartMontNew)
					colorstring.Println("[green]Se monto con exito")
					return true
				} else {
					var partStar int64 = 0
					if mbr.Mbr_partition_1.Part_type == 'E' {
						partStar = mbr.Mbr_partition_1.Part_start
					} else if mbr.Mbr_partition_2.Part_type == 'E' {
						partStar = mbr.Mbr_partition_2.Part_start
					} else if mbr.Mbr_partition_3.Part_type == 'E' {
						partStar = mbr.Mbr_partition_3.Part_start
					} else if mbr.Mbr_partition_4.Part_type == 'E' {
						partStar = mbr.Mbr_partition_4.Part_start
					}
					if partStar != 0 {
						file.Seek(partStar, 0)
						ebr := Estruct.EBR{}
						A := ReadBytes(file, int(unsafe.Sizeof(Estruct.EBR{})))
						buffer := bytes.NewBuffer(A)
						err := binary.Read(buffer, binary.BigEndian, &ebr)
						if err != nil {
							colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
						}
						montar := FinPartLogic(file, ebr, NAME)
						if montar == true {
							PartMontNew := Estruct.MountFisic{Path: path}
							PartMontNew.Name = name
							PartMontNew.Id = "vd" + Disco.Letra + strconv.Itoa(int(Disco.Numero))
							DiscosMontados[i].Numero = Disco.Numero + 1
							ParticionesMontada = append(ParticionesMontada, PartMontNew)
							colorstring.Println("[green]Se monto con exito")
							return true
						}
					}
					colorstring.Println("[red]No se pudo montar la particion")
					return false
				}
			}
		}
		newDisc := Estruct.Mount{Letra: string(Estruct.ABC[Estruct.NumeroLetra])}
		Estruct.NumeroLetra++
		newDisc.Numero = 2
		newDisc.Path = path
		DiscosMontados = append(DiscosMontados, newDisc)
		var NAME [16]byte
		copy(NAME[:], name)
		if mbr.Mbr_partition_1.Part_name == NAME && mbr.Mbr_partition_1.Part_type != 'E' {
			PartMontNew := Estruct.MountFisic{Path: path}
			PartMontNew.Name = name
			PartMontNew.Id = "vd" + newDisc.Letra + "1"
			ParticionesMontada = append(ParticionesMontada, PartMontNew)
			colorstring.Println("[green]Se monto con exito")
			return true
		} else if mbr.Mbr_partition_2.Part_name == NAME && mbr.Mbr_partition_2.Part_type != 'E' {
			PartMontNew := Estruct.MountFisic{Path: path}
			PartMontNew.Name = name
			PartMontNew.Id = "vd" + newDisc.Letra + "1"
			ParticionesMontada = append(ParticionesMontada, PartMontNew)
			colorstring.Println("[green]Se monto con exito")
			return true
		} else if mbr.Mbr_partition_3.Part_name == NAME && mbr.Mbr_partition_3.Part_type != 'E' {
			PartMontNew := Estruct.MountFisic{Path: path}
			PartMontNew.Name = name
			PartMontNew.Id = "vd" + newDisc.Letra + "1"
			ParticionesMontada = append(ParticionesMontada, PartMontNew)
			colorstring.Println("[green]Se monto con exito")
			return true
		} else if mbr.Mbr_partition_4.Part_name == NAME && mbr.Mbr_partition_4.Part_type != 'E' {
			PartMontNew := Estruct.MountFisic{Path: path}
			PartMontNew.Name = name
			PartMontNew.Id = "vd" + newDisc.Letra + "1"
			ParticionesMontada = append(ParticionesMontada, PartMontNew)
			colorstring.Println("[green]Se monto con exito")
			return true
		} else {
			var partStar int64 = 0
			if mbr.Mbr_partition_1.Part_type == 'E' {
				partStar = mbr.Mbr_partition_1.Part_start
			} else if mbr.Mbr_partition_2.Part_type == 'E' {
				partStar = mbr.Mbr_partition_2.Part_start
			} else if mbr.Mbr_partition_3.Part_type == 'E' {
				partStar = mbr.Mbr_partition_3.Part_start
			} else if mbr.Mbr_partition_4.Part_type == 'E' {
				partStar = mbr.Mbr_partition_4.Part_start
			}
			if partStar != 0 {
				file.Seek(partStar, 0)
				ebr := Estruct.EBR{}
				A := ReadBytes(file, int(unsafe.Sizeof(Estruct.EBR{})))
				buffer := bytes.NewBuffer(A)
				err := binary.Read(buffer, binary.BigEndian, &ebr)
				if err != nil {
					colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
				}
				montar := FinPartLogic(file, ebr, NAME)
				if montar == true {
					PartMontNew := Estruct.MountFisic{Path: path}
					PartMontNew.Name = name
					PartMontNew.Id = "vd" + newDisc.Letra + "1"
					ParticionesMontada = append(ParticionesMontada, PartMontNew)
					colorstring.Println("[green]Se monto con exito")
					return true
				}
			}
			colorstring.Println("[red]No se pudo montar la particion")
			return false
		}

	}
	colorstring.Println("[red]No se pudo montar la particion")
	return false
}

//UNMOUNT Comando script de MIA
func UNMOUNT(Ids []string) bool {
	for i, value := range ParticionesMontada {
		for j, id := range Ids {
			if value.Id == id {
				colorstring.Println("[green]\tSe Desmonto la particion con Id: " + value.Id + " Path: " + value.Path + " Name: " + value.Name)
				ParticionesMontada[i] = Estruct.MountFisic{}
				Ids[j] = ""
			}
		}
	}
	for _, id := range Ids {
		if id != "" {
			colorstring.Println("[red]\tNo se encontro el Id: " + id)
		}
	}
	return true
}

//removeCom es el encargado de quitarle las comillas a un texto
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

//writeParticion encargado de mostrar contenido de mbr
func writeParticion(mbr Estruct.MBR) {
	println(mbr.Mbr_disk_signature)
	println(string(mbr.Mbr_fecha_creacion[:]))
	println(mbr.Mbr_tamaño)
	println(string(mbr.Mbr_partition_1.Part_name[:]))
	println(string(mbr.Mbr_partition_1.Part_fit))
	println(mbr.Mbr_partition_1.Part_size)
	println(mbr.Mbr_partition_1.Part_start)
	println(string(mbr.Mbr_partition_1.Part_status))
	println(string(mbr.Mbr_partition_1.Part_type))
	println(string(mbr.Mbr_partition_2.Part_name[:]))
	println(string(mbr.Mbr_partition_2.Part_fit))
	println(mbr.Mbr_partition_2.Part_size)
	println(mbr.Mbr_partition_2.Part_start)
	println(string(mbr.Mbr_partition_2.Part_status))
	println(string(mbr.Mbr_partition_2.Part_type))
	println(string(mbr.Mbr_partition_3.Part_name[:]))
	println(string(mbr.Mbr_partition_3.Part_fit))
	println(mbr.Mbr_partition_3.Part_size)
	println(mbr.Mbr_partition_3.Part_start)
	println(string(mbr.Mbr_partition_3.Part_status))
	println(string(mbr.Mbr_partition_3.Part_type))
	println(string(mbr.Mbr_partition_4.Part_name[:]))
	println(string(mbr.Mbr_partition_4.Part_fit))
	println(mbr.Mbr_partition_4.Part_size)
	println(mbr.Mbr_partition_4.Part_start)
	println(string(mbr.Mbr_partition_4.Part_status))
	println(string(mbr.Mbr_partition_4.Part_type))
}

//createPartition Es el que Escribe el mbr al inicio del archivo
func createPartition(file *os.File, mbr Estruct.MBR) {
	file2, _ := os.OpenFile(file.Name(), os.O_WRONLY, 0775)
	defer file2.Close()
	file2.Seek(0, 0)
	mbr2 := &mbr
	var WriteStrucBinary bytes.Buffer
	err := binary.Write(&WriteStrucBinary, binary.BigEndian, mbr2)
	if err != nil {
		println(err.Error())
	}
	WriteByte(file2, WriteStrucBinary.Bytes())
}

//createPartitionLogic Es el que Escribe el ebr en la posicion de partstar
func createPartitionLogic(file *os.File, ebr Estruct.EBR, partStar int64) {
	file2, _ := os.OpenFile(file.Name(), os.O_WRONLY, 0775)
	defer file2.Close()
	file2.Seek(partStar, 0)
	ebr2 := &ebr
	var WriteStrucBinary bytes.Buffer
	err := binary.Write(&WriteStrucBinary, binary.BigEndian, ebr2)
	if err != nil {
		println(err.Error())
	}
	WriteByte(file2, WriteStrucBinary.Bytes())
}

//RecursivoFindLogic Busca e inserta donde se pueda la particion logica nueva
func RecursivoFindLogic(file *os.File, ebr Estruct.EBR, ebrInsert Estruct.EBR, ExtSize int64, SumLogic int64) {
	if ebr.Part_name == ebrInsert.Part_name {
		colorstring.Println("[red]En las particiones logicas no se puede repetir el nombre")
	} else {
		if ebr.Part_next == -1 {
			if SumLogic <= ExtSize {
				if ebr.Part_status != 'A' && (ebr.Part_next == -1) {
					ebr.Part_name = ebrInsert.Part_name
					ebr.Part_fit = ebrInsert.Part_fit
					ebr.Part_size = ebrInsert.Part_size
					ebr.Part_status = ebrInsert.Part_status
					ebr.Part_next = ebr.Part_start + int64(unsafe.Sizeof(Estruct.EBR{})) + 1
					createPartitionLogic(file, ebr, ebr.Part_start)
					ebr2 := Estruct.EBR{Part_next: -1}
					ebr2.Part_start = ebr.Part_start + int64(unsafe.Sizeof(Estruct.EBR{})) + 1
					createPartitionLogic(file, ebr2, ebr2.Part_start)
				}
			} else {
				colorstring.Println("[red]No existe espacio suficiente en la particion Extendida")
			}
		} else if ebr.Part_status != 'A' && ebr.Part_size >= ebrInsert.Part_size {
			ebr.Part_name = ebrInsert.Part_name
			ebr.Part_fit = ebrInsert.Part_fit
			ebr.Part_size = ebrInsert.Part_size
			ebr.Part_status = ebrInsert.Part_status
			createPartitionLogic(file, ebr, ebr.Part_start)
		} else {
			file.Seek(ebr.Part_next, 0)
			ebr2 := Estruct.EBR{}
			A := ReadBytes(file, int(unsafe.Sizeof(Estruct.EBR{})))
			buffer := bytes.NewBuffer(A)
			err := binary.Read(buffer, binary.BigEndian, &ebr2)
			if err != nil {
				colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
			} else {
				RecursivoFindLogic(file, ebr2, ebrInsert, ExtSize, SumLogic+ebr2.Part_size)
			}
		}
	}
}

func showLogicRec(file *os.File, ebr Estruct.EBR) {
	if ebr.Part_next == -1 {
		println(string(ebr.Part_name[:]))
		println(ebr.Part_size)
	} else {
		println(string(ebr.Part_name[:]))
		println(ebr.Part_size)
		file.Seek(ebr.Part_next, 0)
		ebr2 := Estruct.EBR{}
		A := ReadBytes(file, int(unsafe.Sizeof(Estruct.EBR{})))
		buffer := bytes.NewBuffer(A)
		err := binary.Read(buffer, binary.BigEndian, &ebr2)
		if err != nil {
			colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
		} else {
			showLogicRec(file, ebr2)
		}
	}
}

//FinPartLogic es la encargada de retonar si existe la particion para montarla
func FinPartLogic(file *os.File, ebr Estruct.EBR, NAME [16]byte) bool {
	if ebr.Part_name == NAME {
		return true
	} else {
		if ebr.Part_next != -1 {
			file.Seek(ebr.Part_next, 0)
			ebr2 := Estruct.EBR{}
			A := ReadBytes(file, int(unsafe.Sizeof(Estruct.EBR{})))
			buffer := bytes.NewBuffer(A)
			err := binary.Read(buffer, binary.BigEndian, &ebr2)
			if err != nil {
				colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
			}
			return FinPartLogic(file, ebr2, NAME)
		} else if ebr.Part_next == -1 {
			return false
		}
	}
	return false
}

//RecDeletePartitionLogic metodo recursivo que elimina las particiones logicas
func RecDeletePartitionLogic(file *os.File, ebr Estruct.EBR, name [16]byte, tipeDelete string) int64 {
	if ebr.Part_name == name {
		if tipeDelete == "FAST" {
			ebr2 := Estruct.EBR{Part_next: ebr.Part_next}
			ebr2.Part_size = ebr.Part_size
			ebr2.Part_start = ebr.Part_start
			createPartitionLogic(file, ebr2, ebr.Part_start)
			return 0
		} else if tipeDelete == "FULL" {
			return ebr.Part_next
		}
	} else {
		if ebr.Part_next == -1 {
			colorstring.Println("[red]No se encontro ninguna particion con ese nombre")
			return 0
		} else {
			if ebr.Part_next != -1 {
				file.Seek(ebr.Part_next, 0)
			}
			ebr2 := Estruct.EBR{}
			A := ReadBytes(file, int(unsafe.Sizeof(ebr2)))
			buffer := bytes.NewBuffer(A)
			err := binary.Read(buffer, binary.BigEndian, &ebr2)
			if err != nil {
				colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
				return 0
			}
			Valor := RecDeletePartitionLogic(file, ebr2, name, tipeDelete)
			if Valor != 0 {
				ebr.Part_next = Valor
				createPartitionLogic(file, ebr, ebr.Part_start)
				return 0
			} else {
				return 0
			}
		}
	}
	return 0
}

//WriteByte metodo que Escribe en el archivo
func WriteByte(file *os.File, bytes []byte) {
	println(bytes)
	_, err := file.Write(bytes)
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al escribir en el archivo " + file.Name() + "\n" + err.Error())
	}

}

//ReadBytes metodo que Lee en el archivo
func ReadBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number) //array de bytes
	_, err := file.Read(bytes)    // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
