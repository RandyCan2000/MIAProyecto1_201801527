package Comandos

import (
	Estruct "Proyecto1MIA/Estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/github.com/mitchellh/colorstring"
)

type Login struct {
	ID_user  int
	Grupo    string
	User     string
	Password string
}

var UserLogueado Login

//Mostrar funcion de prueba
func Mostrar() {

}

//MKFS crea el formateo a una particion
func MKFS(id string, tipe string, unit string, add string) bool {
	//Se busca el disco montado con el id
	ParticionMontada, errID := BuscarParticionMontada(id)
	if errID == false {
		colorstring.Println("[red]\tNo se encontro el ID de la particion")
		return false
	}
	//Se lee el mbr que esta al inicio del disco ya que se encontro el id ParticionMontada.Path tiene la ruta del disco
	MBR, err := ReadMBR(ParticionMontada.Path)
	if err != nil {
		return false
	}
	var ParticionesMBR [4]Estruct.Partition
	ParticionesMBR[0] = MBR.Mbr_partition_1
	ParticionesMBR[1] = MBR.Mbr_partition_2
	ParticionesMBR[2] = MBR.Mbr_partition_3
	ParticionesMBR[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMontada.Name)
	var partStar, SizeParticion int64 = 0, 0
	//Se recorre las particion y se busca que nombre coincide y se toma ese valor de partStar osea donde inicia el superbloque
	for _, value := range ParticionesMBR {
		if value.Part_name == NAME {
			partStar = value.Part_start
			SizeParticion = value.Part_size
			break
		}
	}
	//validacion por si F
	if partStar == 0 {
		colorstring.Println("[red]\tNo se encontro el nombre de la particion")
		return false
	}
	//SizeParticion Particion este es el calculo de los datos
	//FORMULAS
	SizeSB := int64(unsafe.Sizeof(Estruct.SuperBoot{}))
	SizeAVD := int64(unsafe.Sizeof(Estruct.AVD{}))
	SizeDD := int64(unsafe.Sizeof(Estruct.DD{}))
	SizeInodo := int64(unsafe.Sizeof(Estruct.INODO{}))
	SizeBloque := int64(unsafe.Sizeof(Estruct.BD{}))
	SizeBitacora := int64(unsafe.Sizeof(Estruct.Bitacora{}))
	var NEstruct int64 = (SizeParticion - (2 * SizeSB)) / (27 + SizeAVD + SizeDD + (5*SizeInodo + (20 * SizeBloque) + SizeBitacora))
	CantidadAVD := NEstruct
	CantidadDD := NEstruct
	CantidadInodo := 5 * NEstruct
	CantidadBloque := 4 * CantidadInodo
	CantidadBitacora := NEstruct

	//SuperBoot
	SBwrite := Estruct.SuperBoot{Sb_arbol_virtual_count: CantidadAVD}
	var PathSplit []string = strings.Split(ParticionMontada.Path, "/")
	copy(SBwrite.Sb_nombre_hd[:], PathSplit[len(PathSplit)-1])
	SBwrite.Sb_detalle_directorio_count = CantidadDD
	SBwrite.Sb_inodos_count = CantidadInodo
	SBwrite.Sb_bloques_count = CantidadBloque
	SBwrite.Sb_bloques_free = SBwrite.Sb_bloques_count
	SBwrite.Sb_inodos_free = SBwrite.Sb_inodos_count
	SBwrite.Sb_arbol_virtual_free = SBwrite.Sb_arbol_virtual_count
	SBwrite.Sb_detalle_directorio_free = SBwrite.Sb_detalle_directorio_count
	copy(SBwrite.Sb_date_creacion[:], StringFechaActual())
	copy(SBwrite.Sb_date_ultimo_montaje[:], StringFechaActual())
	SBwrite.Sb_montaje_count = 1
	SBwrite.Sb_ap_bitmap_arbol_directorio = partStar + int64(unsafe.Sizeof(Estruct.SuperBoot{}))
	SBwrite.Sb_ap_arbol_directorio = SBwrite.Sb_ap_bitmap_arbol_directorio + CantidadAVD
	SBwrite.Sb_ap_bitmap_detalle_directorio = SBwrite.Sb_ap_bitmap_arbol_directorio + (SizeAVD * CantidadAVD)
	SBwrite.Sb_ap_detalle_directorio = SBwrite.Sb_ap_bitmap_detalle_directorio + CantidadDD
	SBwrite.Sb_ap_bitmap_tabla_inodo = SBwrite.Sb_ap_detalle_directorio + (SizeDD * CantidadDD)
	SBwrite.Sb_ap_tabla_inodo = SBwrite.Sb_ap_bitmap_tabla_inodo + CantidadInodo
	SBwrite.Sb_ap_bitmap_bloque = SBwrite.Sb_ap_tabla_inodo + (SizeInodo * CantidadInodo)
	SBwrite.Sb_ap_bloques = SBwrite.Sb_ap_bitmap_bloque + CantidadBloque
	SBwrite.Sb_ap_log = SBwrite.Sb_ap_bloques + (SizeBloque * CantidadBloque)
	SBwrite.Sb_size_struct_arbol_directorio = SizeAVD
	SBwrite.Sb_size_struct_detalle_directorio = SizeDD
	SBwrite.Sb_size_struct_inodo = SizeInodo
	SBwrite.Sb_size_struct_bloque = SizeBloque
	SBwrite.Sb_first_free_bit_arbol_directorio = SBwrite.Sb_ap_arbol_directorio
	SBwrite.Sb_first_free_bit_detalle_directorio = SBwrite.Sb_ap_detalle_directorio
	SBwrite.Sb_first_free_bit_table_inodo = SBwrite.Sb_ap_tabla_inodo
	SBwrite.Sb_first_free_bit_bloque = SBwrite.Sb_ap_bloques
	SBwrite.Sb_magic_num = 201801527
	echo1 := WriteSB(ParticionMontada.Path, SBwrite, partStar)
	//Escritura de la copia de superboot
	PartStarCopySB := SBwrite.Sb_ap_log + (CantidadBitacora * SizeBitacora)
	echo := WriteSB(ParticionMontada.Path, SBwrite, PartStarCopySB)
	if echo == true && echo1 == true {
		colorstring.Println("[blue]\tSe creo el sistema de archivos LWH con exito")
	} else {
		colorstring.Println("[red]\tNo se creo el sistema de archivos LWH con exito")
	}
	final := SBwrite.Sb_ap_bitmap_arbol_directorio + SBwrite.Sb_arbol_virtual_count
	Inicio := SBwrite.Sb_ap_bitmap_arbol_directorio
	for i := Inicio; i < final; i++ {
		WriteOneByteCero(ParticionMontada.Path, i)
	}
	for i := SBwrite.Sb_ap_bitmap_detalle_directorio; i < SBwrite.Sb_ap_bitmap_detalle_directorio+SBwrite.Sb_detalle_directorio_count; i++ {
		WriteOneByteCero(ParticionMontada.Path, i)
	}
	for i := SBwrite.Sb_ap_bitmap_tabla_inodo; i < SBwrite.Sb_ap_bitmap_tabla_inodo+SBwrite.Sb_inodos_count; i++ {
		WriteOneByteCero(ParticionMontada.Path, i)
	}
	for i := SBwrite.Sb_ap_bitmap_bloque; i < SBwrite.Sb_ap_bitmap_bloque+10; i++ {
		WriteOneByteCero(ParticionMontada.Path, i)
	}
	//Resetear bitacora
	Bitacora := Estruct.Bitacora{}
	Inicio = SBwrite.Sb_ap_log
	Final := (SBwrite.Sb_ap_log + (SBwrite.Sb_arbol_virtual_count * int64(unsafe.Sizeof(Estruct.Bitacora{})))) - int64(unsafe.Sizeof(Estruct.Bitacora{}))
	SizeBitacora = int64(unsafe.Sizeof(Estruct.Bitacora{}))
	var i int64
	for i = Inicio; i <= Final; i = i + SizeBitacora {
		WriteLog(ParticionMontada.Path, Bitacora, i)
	}
	//Escribir user.txt
	writeUsertxt(ParticionMontada.Path, SBwrite)
	EscribirEnBitacora(ParticionMontada.Path, SBwrite, "MKFILE", 0, "/user.txt", "1,G,root\n"+"1,U,root,root,201801527\n")
	return (echo && echo1)
}

//LOGIN inicia sesion en la particion
func LOGIN(user string, password string, idPartMontada string) bool {
	ParticionMonta, Err := BuscarParticionMontada(idPartMontada)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	AVD, ErrAVD := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
	if ErrAVD != nil {
		colorstring.Println("\t[red]Error al leer el avd raiz")
		return false
	}
	DD, ErrDD := ReadDD(ParticionMonta.Path, AVD.AvdApDetalleDirectorio)
	if ErrDD != nil {
		colorstring.Println("\t[red]Error al leer el detalle de directorio")
		return false
	}
	Inodo, ErrInodo := ReadInodo(ParticionMonta.Path, DD.DDarrayFiles[0].DDfileApInodo)
	if ErrInodo != nil {
		colorstring.Println("\t[red]Error al leer el Inodo")
		return false
	}
	var i int64 = 0
	var TEXTOUSERSTXT string = ""
	for {
		for i = 0; i < 4; i++ {
			if Inodo.IarrayBloque[i] > 0 {
				Bloque, _ := ReadBloque(ParticionMonta.Path, Inodo.IarrayBloque[i])
				TEXTOUSERSTXT += string(Bloque.BDData[:])
			}
		}
		if Inodo.IapOtroInodo == 0 || Inodo.IapOtroInodo == -1 {
			break
		} else {
			Inodo, _ = ReadInodo(ParticionMonta.Path, Inodo.IapOtroInodo)
		}
	}
	RegistroSplit := strings.Split(TEXTOUSERSTXT, "\n")
	for _, value := range RegistroSplit {
		RegSplit := strings.Split(value, ",")
		if len(RegSplit) == 5 {
			if RegSplit[3] == user && RegSplit[4] == password {
				UserLogueado.ID_user, _ = strconv.Atoi(RegSplit[0])
				UserLogueado.Password = RegSplit[4]
				UserLogueado.User = RegSplit[3]
				UserLogueado.Grupo = RegSplit[2]
				colorstring.Println("[blue]\tBienvenido Usuario: " + RegSplit[3])
				return true
			}
		}
	}
	colorstring.Println("[red]\tUsuario y Contraseña no coinciden")
	return false
}

//LOGOUT cierra la sesion iniciada
func LOGOUT() bool {
	if UserLogueado.ID_user == 0 {
		colorstring.Println("[red]\tError ningun usuario esta logueado")
		return false
	} else {
		colorstring.Println("[blue]\tAdios Usuario :" + UserLogueado.User)
		UserLogueado.ID_user = 0
		UserLogueado.Password = ""
		UserLogueado.User = ""
		UserLogueado.Grupo = ""
		return true
	}
}

//MKDIR crea un directorio en el disco especificado
func MKDIR(id string, path string, p bool) bool {
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	AVDRaiz, errAVDRaiz := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
	if errAVDRaiz != nil {
		colorstring.Println("[red]\tError al crear Carpeta")
		return false
	}
	Carpetas := strings.Split(path, "/")
	_, i := RecorrerYCrearAVD(ParticionMonta.Path, SB, Carpetas, 1, p, AVDRaiz, SB.Sb_ap_arbol_directorio)
	if i > 0 {
		EscribirEnBitacora(ParticionMonta.Path, SB, "MKDIR", 1, path, "")
	}
	return true
}

//MKFILE Crea un archivo en la carpeta path
func MKFILE(id string, path string, CrearTodo bool, size string, cont string, Pregunta bool) bool {
	ABECEDARIO := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	SIZE, err := strconv.Atoi(size)
	if size != "" {
		if err != nil {
			colorstring.Println("[red]\tError size debe ser un numero")
			return false
		}
	} else {
		SIZE = len(cont)
	}
	contador := 0
	if len(cont) < SIZE {
		for i := len(cont); i < SIZE; i++ {
			cont += string(ABECEDARIO[contador])
			contador++
			if contador == len(ABECEDARIO) {
				contador = 0
			}
		}
	} else if len(cont) > SIZE {
		contador = 0
		TXT := ""
		for _, char := range cont {
			if contador != SIZE {
				TXT += string(char)
			} else {
				break
			}
			contador++
		}
		cont = TXT
	}
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	pathConfile := strings.Split(path, "/")
	var Carpetas []string
	for key, _ := range pathConfile {
		if key != len(pathConfile)-1 {
			Carpetas = append(Carpetas, pathConfile[key])
		}
	}
	AVDRaiz, _ := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
	AVDFinal, POSAVDFinal := RecorrerYCrearAVD(ParticionMonta.Path, SB, Carpetas, 1, CrearTodo, AVDRaiz, SB.Sb_ap_arbol_directorio)
	if POSAVDFinal == 0 {
		colorstring.Println("[red]No se Creo el archivo")
		return false
	}
	if AVDFinal.AvdApDetalleDirectorio <= 0 {
		Bloques, PBloques, hecho := ReturnBloques(ParticionMonta.Path, SB, cont)
		if hecho == false {
			colorstring.Println("[red]Error al escribir el archivo")
			return false
		}
		Size, _ := strconv.Atoi(size)
		Inodos, PInodos, hecho := ReturnInodos(ParticionMonta.Path, SB, PBloques, int64(UserLogueado.ID_user), int64(Size))
		if hecho == false {
			colorstring.Println("[red]\tOcurrio un error al escrbir el archivo")
			return false
		}
		FileInfo := Estruct.DDInfo{DDfileApInodo: -1}
		copy(FileInfo.DDfileNombre[:], pathConfile[len(pathConfile)-1])
		copy(FileInfo.DDfileDateCreate[:], StringFechaActual())
		if len(Inodos) != 0 {
			FileInfo.DDfileApInodo = PInodos[0]
		}
		DD, PDD, hecho := ReturnDD(ParticionMonta.Path, SB, FileInfo)
		if hecho == false {
			colorstring.Println("[red]\tOcurrio un error al escrbir el archivo")
			return false
		}
		AVDFinal.AvdApDetalleDirectorio = PDD
		//Escribir AVD
		WriteAVD(ParticionMonta.Path, AVDFinal, POSAVDFinal)
		//Escribir DD
		WriteDD(ParticionMonta.Path, DD, PDD)
		//Escribir Inodos
		for key, _ := range Inodos {
			WriteInodo(ParticionMonta.Path, Inodos[key], PInodos[key])
		}
		for key, _ := range Bloques {
			WriteBloque(ParticionMonta.Path, Bloques[key], PBloques[key])
		}
		return true
	} else {
		DD, _ := ReadDD(ParticionMonta.Path, AVDFinal.AvdApDetalleDirectorio)
		FileInfo := Estruct.DDInfo{DDfileApInodo: -1}
		copy(FileInfo.DDfileDateUpdate[:], StringFechaActual())
		copy(FileInfo.DDfileNombre[:], pathConfile[len(pathConfile)-1])
		RecorrerYCrearDD(ParticionMonta.Path, SB, FileInfo, DD, AVDFinal.AvdApDetalleDirectorio, cont, Pregunta)
	}
	EscribirEnBitacora(ParticionMonta.Path, SB, "MKFILE", 0, path, cont)
	return true
}

//REP Script de MIA generador de reportes
func REP(id string, nombre string, path string, ruta string) bool {
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	pathDir := ""
	//spliter
	PathSplit := strings.Split(path, "/")
	Extencion := strings.Split(PathSplit[len(PathSplit)-1], ".")
	for i := 0; i < len(PathSplit)-1; i++ {
		pathDir += "/" + PathSplit[i]
	}
	//Crear Directorio
	err := os.MkdirAll(pathDir, 0775)
	if err != nil {
		colorstring.Println("[red]\tOcurrio un error al crear la carpeta donde se guarda el reporte")
		return false
	}
	//Reportes
	if nombre == "DISK" {
		Disk(ParticionMonta.Path, Extencion[len(Extencion)-1], path)
	} else if nombre == "TREE_COMPLETE" {
		TreeComplete(ParticionMonta.Path, Extencion[len(Extencion)-1], path, ParticionMonta.Name, true)
	} else if nombre == "DIRECTORIO" {
		TreeComplete(ParticionMonta.Path, Extencion[len(Extencion)-1], path, ParticionMonta.Name, false)
	} else if nombre == "BM_ARBDIR" {
		Bitmap(ParticionMonta.Path, SB.Sb_ap_bitmap_arbol_directorio, SB.Sb_arbol_virtual_count, path)
	} else if nombre == "BM_DETDIR" {
		Bitmap(ParticionMonta.Path, SB.Sb_ap_bitmap_detalle_directorio, SB.Sb_detalle_directorio_count, path)
	} else if nombre == "BM_INODE" {
		Bitmap(ParticionMonta.Path, SB.Sb_ap_bitmap_tabla_inodo, SB.Sb_inodos_count, path)
	} else if nombre == "BM_BLOCK" {
		Bitmap(ParticionMonta.Path, SB.Sb_ap_bitmap_bloque, SB.Sb_bloques_count, path)
	} else if nombre == "BITACORA" {
		LOG(ParticionMonta.Path, SB, path, Extencion[len(Extencion)-1])
	} else if nombre == "SB" {
		SBGraficador(SB, path, Extencion[len(Extencion)-1])
	} else if nombre == "TREE_FILE" {
		file, err := os.OpenFile(ParticionMonta.Path, os.O_RDWR, 0775)
		if err != nil {
			colorstring.Println("[red]\tError Al abrir el Disco")
			return false
		}
		AVDRaiz, _ := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
		RUTA := RecorrerEstructuras(file, "/", AVDRaiz, false)
		if RUTA != "EXIT" {
			RutaSplit := strings.Split(RUTA, "/")
			TREEFILE(ParticionMonta.Path, SB, path, Extencion[len(Extencion)-1], RutaSplit, true)
		}
	} else if nombre == "TREE_DIRECTORIO" {
		file, err := os.OpenFile(ParticionMonta.Path, os.O_RDWR, 0775)
		if err != nil {
			colorstring.Println("[red]\tError Al abrir el Disco")
			return false
		}
		AVDRaiz, _ := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
		RUTA := RecorrerEstructuras(file, "/", AVDRaiz, true)
		if RUTA != "EXIT" {
			RutaSplit := strings.Split(RUTA, "/")
			TREEDIRECTORIOS(ParticionMonta.Path, SB, path, Extencion[len(Extencion)-1], RutaSplit)
		}
	} else if nombre == "MBR" {
		MBRG(ParticionMonta.Path, SB, path, Extencion[len(Extencion)-1])
	}

	return true
}

//CAT Muestra el contenido de un archivo
func CAT(path []string, id string) bool {
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	AVDRaiz, ErrAVD := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
	if ErrAVD != nil {
		colorstring.Println("\t[red]Error al abrir leer el AVD")
		return false
	}
	for _, val := range path {
		colorstring.Println("[yellow]\t" + val)
		RutaSplit := strings.Split(val, "/")
		_, _, _, TEXTO := FindFile(ParticionMonta.Path, RutaSplit, 1, AVDRaiz)
		colorstring.Println("[blue]\t" + TEXTO)
	}
	return true
}

//RM Remueve el archivo del path especificado
func RM(id string, path string) bool {
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	AVDRaiz, ErrAVD := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
	if ErrAVD != nil {
		colorstring.Println("\t[red]Error al abrir leer el AVD")
		return false
	}
	RutaSplit := strings.Split(path, "/")
	POSDD, DD, KEYDDInfo, _ := FindFile(ParticionMonta.Path, RutaSplit, 1, AVDRaiz)
	BorrarFile(ParticionMonta.Path, path, SB, DD, POSDD, KEYDDInfo)
	return true
}

//MKGRP Crea un grupo en archivo user.txt
func MKGRP(id string, name string) bool {
	if UserLogueado.ID_user != 1 {
		colorstring.Println("[red]\tSolo el usuario root puede usar este script")
		return false
	}
	if name == "" {
		colorstring.Println("[red]\tError No puede estar el nombre del grupo en blanco")
		return false
	}
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	AVDRaiz, ErrAVD := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
	if ErrAVD != nil {
		colorstring.Println("\t[red]Error al abrir leer el AVD")
		return false
	}
	var RutaSplit []string
	RutaSplit = append(RutaSplit, "")
	RutaSplit = append(RutaSplit, "user.txt")
	_, _, _, TEXTO := FindFile(ParticionMonta.Path, RutaSplit, 1, AVDRaiz)
	TextoSplit := strings.Split(TEXTO, "\n")
	INDICEMAY := 0
	for _, value := range TextoSplit {
		RegistroSplit := strings.Split(value, ",")
		if len(RegistroSplit) == 3 {
			if RegistroSplit[2] == name && RegistroSplit[0] != "0" {
				colorstring.Println("[red]\tError El grupo ya existe")
				return false
			}
			IDG, _ := strconv.Atoi(RegistroSplit[0])
			if IDG > INDICEMAY {
				INDICEMAY = IDG
			}
		}
	}
	INDICEMAY++
	TEXTO += strconv.Itoa(INDICEMAY) + ",G," + name + "\n"
	MKFILE(id, "/user.txt", true, "", TEXTO, true)
	return true
}

//RMGRP Crea un grupo en archivo user.txt
func RMGRP(id string, name string) bool {
	if UserLogueado.ID_user != 1 {
		colorstring.Println("[red]\tSolo el usuario root puede usar este script")
		return false
	}
	if name == "" {
		colorstring.Println("[red]\tError No puede estar el nombre del grupo en blanco")
		return false
	}
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	AVDRaiz, ErrAVD := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
	if ErrAVD != nil {
		colorstring.Println("\t[red]Error al abrir leer el AVD")
		return false
	}
	var RutaSplit []string
	RutaSplit = append(RutaSplit, "")
	RutaSplit = append(RutaSplit, "user.txt")
	_, _, _, TEXTO := FindFile(ParticionMonta.Path, RutaSplit, 1, AVDRaiz)
	TextoSplit := strings.Split(TEXTO, "\n")
	NTEXTO := ""
	for _, value := range TextoSplit {
		RegistroSplit := strings.Split(value, ",")
		if len(RegistroSplit) == 3 {
			if RegistroSplit[2] == name {
				if RegistroSplit[0] == "0" {
					colorstring.Println("[red]\tError El grupo ya fue Eliminado")
					return false
				}
				RegistroSplit[0] = "0"
			}
		}
		for key, arg := range RegistroSplit {
			if key == len(RegistroSplit)-1 {
				NTEXTO += arg + "\n"
			} else {
				NTEXTO += arg + ","
			}
		}
	}
	MKFILE(id, "/user.txt", true, "", NTEXTO, true)
	return true
}

//MKUSER Crea el registro de usuario en el archivo user.txt
func MKUSER(id string, user string, password string, grupo string) bool {
	if UserLogueado.ID_user != 1 {
		colorstring.Println("[red]\tSolo el usuario root puede usar este script")
		return false
	}
	if user == "" || password == "" || grupo == "" {
		colorstring.Println("[red]\tError Faltan parametros")
		return false
	}
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	AVDRaiz, ErrAVD := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
	if ErrAVD != nil {
		colorstring.Println("\t[red]Error al abrir leer el AVD")
		return false
	}
	var RutaSplit []string
	RutaSplit = append(RutaSplit, "")
	RutaSplit = append(RutaSplit, "user.txt")
	_, _, _, TEXTO := FindFile(ParticionMonta.Path, RutaSplit, 1, AVDRaiz)
	TextoSplit := strings.Split(TEXTO, "\n")
	IDMay := 0
	existe := false
	NTEXTO := ""
	for _, value := range TextoSplit {
		RegistroSplit := strings.Split(value, ",")
		if len(RegistroSplit) == 3 { //Grupo
			if RegistroSplit[2] == grupo && RegistroSplit[0] != "0" {
				existe = true
			}
		} else if len(RegistroSplit) == 5 { //Usuarios
			ID, _ := strconv.Atoi(RegistroSplit[0])
			if ID > IDMay {
				IDMay = ID
			}
		}
		if value != "" {
			NTEXTO += value + "\n"
		}
	}
	if existe == false {
		colorstring.Println("[red]\tEl grupo no existe")
		return false
	}
	NTEXTO += strconv.Itoa(IDMay+1) + ",U," + grupo + "," + user + "," + password + "\n"
	MKFILE(id, "/user.txt", true, "", NTEXTO, true)
	return true
}

//RMUSR Crea un grupo en archivo user.txt
func RMUSR(id string, usuarios string) bool {
	if UserLogueado.ID_user != 1 {
		colorstring.Println("[red]\tSolo el usuario root puede usar este script")
		return false
	}
	if usuarios == "" {
		colorstring.Println("[red]\tError No puede estar el usuario del grupo en blanco")
		return false
	}
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	AVDRaiz, ErrAVD := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
	if ErrAVD != nil {
		colorstring.Println("\t[red]Error al abrir leer el AVD")
		return false
	}
	var RutaSplit []string
	RutaSplit = append(RutaSplit, "")
	RutaSplit = append(RutaSplit, "user.txt")
	_, _, _, TEXTO := FindFile(ParticionMonta.Path, RutaSplit, 1, AVDRaiz)
	TextoSplit := strings.Split(TEXTO, "\n")
	NTEXTO := ""
	for _, value := range TextoSplit {
		RegistroSplit := strings.Split(value, ",")
		if len(RegistroSplit) == 5 {
			if usuarios != RegistroSplit[3] {
				NTEXTO += value + "\n"
			}
		} else {
			if value != "" {
				NTEXTO += value + "\n"
			}
		}
	}
	MKFILE(id, "/user.txt", true, "", NTEXTO, true)
	return true
}

//CP Crea una copia del archivo
func CP(id string, path string, dest string) bool {
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	AVDRaiz, ErrAVD := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
	if ErrAVD != nil {
		colorstring.Println("\t[red]Error al abrir leer el AVD")
		return false
	}
	RutaSplit := strings.Split(path, "/")
	_, _, _, TEXTO := FindFile(ParticionMonta.Path, RutaSplit, 1, AVDRaiz)
	if dest[len(dest)-1] == '/' {
		MKFILE(id, dest+RutaSplit[len(RutaSplit)-1], false, "", TEXTO, true)
	} else {
		MKFILE(id, dest+"/"+RutaSplit[len(RutaSplit)-1], false, "", TEXTO, true)
	}
	return true
}

//MV script que mueve un archivo de un directorio a otro
func MV(id string, idDestiny string, path string, pathDestiny string) bool {
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	AVDRaiz, ErrAVD := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
	if ErrAVD != nil {
		colorstring.Println("\t[red]Error al abrir leer el AVD")
		return false
	}
	RutaSplit := strings.Split(path, "/")
	_, _, _, TEXTO := FindFile(ParticionMonta.Path, RutaSplit, 1, AVDRaiz)
	DESTINO := ""
	if pathDestiny[len(pathDestiny)-1] == '/' {
		DESTINO = pathDestiny + RutaSplit[len(RutaSplit)-1]
	} else {
		DESTINO = pathDestiny + "/" + RutaSplit[len(RutaSplit)-1]
	}
	hecho := MKFILE(idDestiny, DESTINO, false, "", TEXTO, true)
	if hecho == false {
		colorstring.Println("[red]\tNo se logro mover el archivo")
		return false
	}
	RM(id, path)
	return true
}

//FIND busca un archivo con el nombre ingresado
func FIND(id string, path string, nombre string) bool {
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	AVDRaiz, ErrAVD := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
	if ErrAVD != nil {
		colorstring.Println("\t[red]Error al abrir leer el AVD")
		return false
	}
	CARPETAS := strings.Split(path, "/")
	AVD, posAVD := RecorrerYCrearAVD(ParticionMonta.Path, SB, CARPETAS, 1, false, AVDRaiz, SB.Sb_ap_arbol_directorio)
	if posAVD <= 0 {
		colorstring.Println("[red]\tError la carpeta no existe")
		return false
	}
	if path == "/" {
		if nombre == "*" {
			Busqueda(ParticionMonta.Path, "", AVDRaiz, "", 0)
		} else {
			Busqueda(ParticionMonta.Path, "", AVDRaiz, strings.ToUpper(nombre), 0)
		}
	} else {
		if nombre == "*" {
			Busqueda(ParticionMonta.Path, path, AVD, "", 0)
		} else {
			Busqueda(ParticionMonta.Path, path, AVD, strings.ToUpper(nombre), 0)
		}
	}

	return true
}

//LOSS simula la perdida de informacion
func LOSS(id string) bool {
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	ParticionMonta.CopySB = SB.Sb_ap_log + (SB.Sb_arbol_virtual_count * int64(unsafe.Sizeof(Estruct.Bitacora{})))
	file, err := os.OpenFile(ParticionMonta.Path, os.O_WRONLY, 0775)
	if err != nil {
		colorstring.Println("[red]Error al abrir el archivo")
	}
	file.Seek(PartStartParticion+int64(int(unsafe.Sizeof(Estruct.SuperBoot{}))), 0)
	var cero byte = 0
	Ins := &cero
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, Ins)
	colorstring.Println("[blue]\tFORMATEANDO...")
	for i := PartStartParticion + int64(int(unsafe.Sizeof(Estruct.SuperBoot{}))); i < SB.Sb_ap_log; i++ {
		WriteByte(file, binario.Bytes())
	}
	colorstring.Println("[blue]\tFORMATEO COMPLETADO EXITOSAMENTE")
	file.Close()
	return true
}

//REN RENOMBRA ARCHIVOS
func REN(id string, path string, nombre string) bool {
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	AVDRaiz, ErrAVD := ReadAVD(ParticionMonta.Path, SB.Sb_ap_arbol_directorio)
	if ErrAVD != nil {
		colorstring.Println("\t[red]Error al abrir leer el AVD")
		return false
	}
	RutaSplit := strings.Split(path, "/")
	posDD, DD, keyDD, _ := FindFile(ParticionMonta.Path, RutaSplit, 1, AVDRaiz)
	if posDD == 0 {
		colorstring.Println("[red]\tNo se encontro el archivo")
		return false
	}
	copy(DD.DDarrayFiles[keyDD].DDfileNombre[:], nombre)
	WriteDD(ParticionMonta.Path, DD, posDD)
	return true
}

func REC(id string) bool {
	ParticionMonta, Err := BuscarParticionMontada(id)
	if Err == false {
		colorstring.Println("\t[red]No se encontro el id ingresado")
		return false
	}
	MBR, Errmbr := ReadMBR(ParticionMonta.Path)
	if Errmbr != nil {
		colorstring.Println("\t[red]Error al abrir el archivo")
		return false
	}
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var NAME [16]byte
	copy(NAME[:], ParticionMonta.Name)
	var PartStartParticion int64 = 0
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			PartStartParticion = Part[key].Part_start
			break
		}
	}
	if PartStartParticion == 0 {
		colorstring.Println("\t[red]No se encontro el Inicio de la particion")
		return false
	}
	SB, ErrSB := ReadSB(ParticionMonta.Path, PartStartParticion)
	if ErrSB != nil {
		colorstring.Println("\t[red]Error al abrir leer el super boot")
		return false
	}
	var mkdir [16]byte
	var mkfile [16]byte
	copy(mkdir[:], "MKDIR")
	copy(mkfile[:], "MKFILE")
	Inicio := SB.Sb_ap_log
	Final := (SB.Sb_ap_log + (SB.Sb_arbol_virtual_count * int64(unsafe.Sizeof(Estruct.Bitacora{})))) - int64(unsafe.Sizeof(Estruct.Bitacora{}))
	SizeBitacora := int64(unsafe.Sizeof(Estruct.Bitacora{}))
	AVDRAIZ := Estruct.AVD{AvdProper: int64(UserLogueado.ID_user)}
	copy(AVDRAIZ.AvdNombreDirectorio[:], "/")
	WriteAVD(ParticionMonta.Path, AVDRAIZ, SB.Sb_ap_arbol_directorio)
	WriteOneByteUno(ParticionMonta.Path, SB.Sb_ap_bitmap_arbol_directorio)
	var i int64
	for i = Inicio; i <= Final; i = i + SizeBitacora {
		LOG, _ := ReadLog(ParticionMonta.Path, i)
		if LOG.LogTipoOperacion == mkdir {
			name := ""
			for _, char := range LOG.LogNombre {
				if char != 0 {
					name += string(char)
				}
			}
			MKDIR(id, name, true)
		} else if LOG.LogTipoOperacion == mkfile {
			name := ""
			for _, char := range LOG.LogNombre {
				if char != 0 {
					name += string(char)
				}
			}
			cont := ""
			for _, char := range LOG.LogContenido {
				if char != 0 {
					cont += string(char)
				}
			}
			MKFILE(id, name, true, "", cont, false)
		} else {
			break
		}
	}
	return true
}

func Busqueda(path string, Ruta string, AVD Estruct.AVD, Filtro string, Identacion int) {
	NAMEFILE := ""
	for _, char := range AVD.AvdNombreDirectorio {
		if char != 0 {
			NAMEFILE += string(char)
		}
	}
	Identacion++
	if AVD.AvdApDetalleDirectorio > 0 {
		//colorstring.Println("[blue]ARCHIVOS:")
		DD, _ := ReadDD(path, AVD.AvdApDetalleDirectorio)
		for {
			for _, val := range DD.DDarrayFiles {
				if val.DDfileApInodo > 0 {
					NAMEFILE = ""
					for _, char := range val.DDfileNombre {
						if char != 0 {
							NAMEFILE += string(char)
						}
					}
					if Filtro == "" {
						colorstring.Println("[yellow]" + AGGI(Identacion) + Ruta + "/" + NAMEFILE)
					} else {
						if strings.Contains(strings.ToUpper(NAMEFILE), strings.ToUpper(Filtro)) {
							colorstring.Println("[yellow]" + AGGI(Identacion) + Ruta + "/" + NAMEFILE)
						}
					}
				}
			}
			if DD.DDapDetalleDirectorio > 0 {
				DD, _ = ReadDD(path, DD.DDapDetalleDirectorio)
			} else {
				break
			}
		}
	}
	AVDAux := Estruct.AVD{}
	//colorstring.Println("[blue]CARPETAS:")
	for {
		for _, val := range AVD.AvdApArraySubDirectorios {
			if val > 0 {
				AVDAux, _ = ReadAVD(path, val)
				NAMEFILE = ""
				for _, char := range AVDAux.AvdNombreDirectorio {
					if char != 0 {
						NAMEFILE += string(char)
					}
				}
				if Filtro == "" {
					colorstring.Println("[yellow]" + AGGI(Identacion) + Ruta + "/" + NAMEFILE)
				} else {
					if strings.Contains(strings.ToUpper(NAMEFILE), strings.ToUpper(Filtro)) {
						colorstring.Println("[yellow]" + AGGI(Identacion) + Ruta + "/" + NAMEFILE)
					}
				}
				Busqueda(path, Ruta+"/"+NAMEFILE, AVDAux, Filtro, Identacion+1)
			}
		}
		if AVD.AvdApArbolVirtualDirectorio > 0 {
			AVD, _ = ReadAVD(path, AVD.AvdApArbolVirtualDirectorio)
		} else {
			break
		}
	}
}

//FindFile Busca el archivo
func FindFile(path string, RutaSplit []string, indCarp int, AVD Estruct.AVD) (int64, Estruct.DD, int, string) {
	AVDAux := Estruct.AVD{}
	NAME := ""
	TEXTO := ""
	var posDD int64 = 0
	var NAME2 [16]byte
	if indCarp == len(RutaSplit)-1 {
		if AVD.AvdApDetalleDirectorio > 0 {
			posDD = AVD.AvdApDetalleDirectorio
			DD, _ := ReadDD(path, AVD.AvdApDetalleDirectorio)
			copy(NAME2[:], RutaSplit[indCarp])
			for {
				for Key, value := range DD.DDarrayFiles {
					if value.DDfileApInodo > 0 && NAME2 == value.DDfileNombre {
						Inodo, _ := ReadInodo(path, value.DDfileApInodo)
						for {
							for _, ValueInodo := range Inodo.IarrayBloque {
								if ValueInodo > 0 {
									BLOQUE, _ := ReadBloque(path, ValueInodo)
									for _, char := range BLOQUE.BDData {
										if char != 0 {
											TEXTO += string(char)
										}
									}
								}
							}
							if Inodo.IapOtroInodo > 0 {
								Inodo, _ = ReadInodo(path, Inodo.IapOtroInodo)
							} else {
								return posDD, DD, Key, TEXTO
							}
						}
					}
				}
				if DD.DDapDetalleDirectorio > 0 {
					posDD = DD.DDapDetalleDirectorio
					DD, _ = ReadDD(path, DD.DDapDetalleDirectorio)
				} else {
					break
				}
			}
		}
	}
	for {
		for _, value := range AVD.AvdApArraySubDirectorios {
			if value > 0 {
				AVDAux, _ = ReadAVD(path, value)
				NAME = ""
				for _, char := range AVDAux.AvdNombreDirectorio {
					if char != 0 {
						NAME += string(char)
					}
				}
				if NAME == RutaSplit[indCarp] {
					contadorCarpetas++
					return FindFile(path, RutaSplit, indCarp+1, AVDAux)
				}
			}
		}
		if AVD.AvdApArbolVirtualDirectorio > 0 {
			AVD, _ = ReadAVD(path, AVD.AvdApArbolVirtualDirectorio)
		} else {
			break
		}
	}
	return 0, Estruct.DD{}, 0, ""
}

//AGGI agrega la cantidad de espacios en blanco espeficada
func AGGI(cantidad int) string {
	ident := ""
	for i := 0; i < cantidad; i++ {
		ident += "  "
	}
	return ident
}

//RecorrerYCrearDD Recorre y crea Escrbie el archivo
func RecorrerYCrearDD(path string, SB Estruct.SuperBoot, FileInfo Estruct.DDInfo, DD Estruct.DD, PosDD int64, cont string, Pregunta bool) bool {
	DDAux := DD
	PosDDAux := PosDD
	for {
		for key, value := range DDAux.DDarrayFiles {
			if value.DDfileApInodo > 0 {
				if value.DDfileNombre == FileInfo.DDfileNombre {
					if Pregunta == true {
						Hacer := MensajeConfirmacion("El archivo que intenta escribir ya existe ¿Desea Sobre Escribir el Contenido? [Y/N]: ", "Y")
						if Hacer == true {
							SobreEscribirArchivo(path, SB, DDAux, FileInfo.DDfileNombre, PosDDAux, key, cont)
							return true
						} else {
							return false
						}
					} else {
						SobreEscribirArchivo(path, SB, DDAux, FileInfo.DDfileNombre, PosDDAux, key, cont)
					}
				}
			}
		}
		if DDAux.DDapDetalleDirectorio > 0 {
			PosDD = DDAux.DDapDetalleDirectorio
			DDAux, _ = ReadDD(path, DDAux.DDapDetalleDirectorio)
		} else {
			break
		}
	}
	DDAux = DD
	PosDDAux = PosDD
	for {
		for key, value := range DDAux.DDarrayFiles {
			if value.DDfileApInodo <= 0 {
				Bloques, PBloques, hecho := ReturnBloques(path, SB, cont)
				if hecho == false {
					colorstring.Println("[red]\tNo se pudo escribir el archivo")
					return false
				}
				Inodos, PInodos, hecho := ReturnInodos(path, SB, PBloques, int64(UserLogueado.ID_user), int64(len(cont)))
				if hecho == false {
					colorstring.Println("[red]\tNo se pudo escribir el archivo")
					return false
				}
				if len(PInodos) != 0 {
					FileInfo.DDfileApInodo = PInodos[0]
				}
				copy(FileInfo.DDfileDateCreate[:], StringFechaActual())
				DDAux.DDarrayFiles[key] = FileInfo
				//Escribir DD
				WriteDD(path, DDAux, PosDDAux)
				//Escribir inodos
				for key, _ := range PInodos {
					WriteInodo(path, Inodos[key], PInodos[key])
				}
				//Escribir Bloques
				for key, _ := range Bloques {
					WriteBloque(path, Bloques[key], PBloques[key])
				}
				name := ""
				for _, char := range FileInfo.DDfileNombre {
					if char != 0 {
						name += string(char)
					}
				}
				return true
			}
		}
		if DDAux.DDapDetalleDirectorio > 0 {
			PosDD = DDAux.DDapDetalleDirectorio
			DDAux, _ = ReadDD(path, DDAux.DDapDetalleDirectorio)
		} else {
			break
		}
	}
	if DDAux.DDapDetalleDirectorio <= 0 {
		DDSig := Estruct.DD{DDapDetalleDirectorio: -1}
		PosBit := PosicionBitmapLibre(path, SB.Sb_ap_bitmap_detalle_directorio, SB.Sb_detalle_directorio_count)
		if PosBit != -1 {
			EscribirDD := SB.Sb_ap_detalle_directorio + (PosBit * SB.Sb_size_struct_detalle_directorio)
			posBitMap := SB.Sb_ap_bitmap_detalle_directorio + PosBit
			WriteOneByteUno(path, posBitMap)
			WriteDD(path, DDSig, EscribirDD)
			DDAux.DDapDetalleDirectorio = EscribirDD
			WriteDD(path, DDAux, PosDDAux)
			return RecorrerYCrearDD(path, SB, FileInfo, DDSig, EscribirDD, cont, Pregunta)
		}
	}
	return true
}

//SobreEscribirArchivo Sobre escribe el archivo ya existente
func SobreEscribirArchivo(path string, SB Estruct.SuperBoot, DD Estruct.DD, name [16]byte, PosDD int64, keyDDinfor int, cont string) bool {
	Inodo, _ := ReadInodo(path, DD.DDarrayFiles[keyDDinfor].DDfileApInodo)
	BitInodoBorrar := (DD.DDarrayFiles[keyDDinfor].DDfileApInodo - SB.Sb_ap_tabla_inodo) / (SB.Sb_size_struct_inodo)
	PosEscribirByte := SB.Sb_ap_bitmap_tabla_inodo + BitInodoBorrar
	WriteOneByteCero(path, PosEscribirByte)
	for {
		for key, _ := range Inodo.IarrayBloque {
			if Inodo.IarrayBloque[key] > 0 {
				BitBloqueBorrar := (Inodo.IarrayBloque[key] - SB.Sb_ap_bloques) / (SB.Sb_size_struct_bloque)
				PosEscribirByte := SB.Sb_ap_bitmap_bloque + BitBloqueBorrar
				WriteOneByteCero(path, PosEscribirByte)
			}
		}
		if Inodo.IapOtroInodo > 0 {
			BitInodoBorrar := (Inodo.IapOtroInodo - SB.Sb_ap_tabla_inodo) / (SB.Sb_size_struct_inodo)
			PosEscribirByte := SB.Sb_ap_bitmap_tabla_inodo + BitInodoBorrar
			WriteOneByteCero(path, PosEscribirByte)

			Inodo, _ = ReadInodo(path, Inodo.IapOtroInodo)
		} else {
			break
		}
	}
	Bloques, PBloques, hecho := ReturnBloques(path, SB, cont)
	if hecho == false {
		colorstring.Println("[red]\tNo se pudo escribir el archivo")
		return false
	}
	Inodos, PInodos, hecho := ReturnInodos(path, SB, PBloques, int64(UserLogueado.ID_user), int64(len(cont)))
	if hecho == false {
		colorstring.Println("[red]\tNo se pudo escribir el archivo")
		return false
	}
	DD.DDarrayFiles[keyDDinfor].DDfileApInodo = -1
	DD.DDarrayFiles[keyDDinfor].DDfileNombre = name
	if len(PInodos) != 0 {
		DD.DDarrayFiles[keyDDinfor].DDfileApInodo = PInodos[0]
	}
	copy(DD.DDarrayFiles[keyDDinfor].DDfileDateUpdate[:], StringFechaActual())
	//Escribir DD
	WriteDD(path, DD, PosDD)
	//Escribir inodos
	for key, _ := range PInodos {
		WriteInodo(path, Inodos[key], PInodos[key])
	}
	//Escribir Bloques
	for key, _ := range Bloques {
		WriteBloque(path, Bloques[key], PBloques[key])
	}
	name2 := ""
	for _, char := range name {
		if char != 0 {
			name2 += string(char)
		}
	}
	return true
}

//BorrarFile Es el encargado de borrar el archivo
func BorrarFile(path string, Ruta string, SB Estruct.SuperBoot, DD Estruct.DD, PosDD int64, keyDDinfor int) {
	Inodo, _ := ReadInodo(path, DD.DDarrayFiles[keyDDinfor].DDfileApInodo)
	BitInodoBorrar := (DD.DDarrayFiles[keyDDinfor].DDfileApInodo - SB.Sb_ap_tabla_inodo) / (SB.Sb_size_struct_inodo)
	PosEscribirByte := SB.Sb_ap_bitmap_tabla_inodo + BitInodoBorrar
	WriteOneByteCero(path, PosEscribirByte)
	for {
		for key, _ := range Inodo.IarrayBloque {
			if Inodo.IarrayBloque[key] > 0 {
				BitBloqueBorrar := (Inodo.IarrayBloque[key] - SB.Sb_ap_bloques) / (SB.Sb_size_struct_bloque)
				PosEscribirByte := SB.Sb_ap_bitmap_bloque + BitBloqueBorrar
				WriteOneByteCero(path, PosEscribirByte)
			}
		}
		if Inodo.IapOtroInodo > 0 {
			BitInodoBorrar := (Inodo.IapOtroInodo - SB.Sb_ap_tabla_inodo) / (SB.Sb_size_struct_inodo)
			PosEscribirByte := SB.Sb_ap_bitmap_tabla_inodo + BitInodoBorrar
			WriteOneByteCero(path, PosEscribirByte)
			Inodo, _ = ReadInodo(path, Inodo.IapOtroInodo)
		} else {
			break
		}
	}
	copy(DD.DDarrayFiles[keyDDinfor].DDfileDateUpdate[:], StringFechaActual())
	DD.DDarrayFiles[keyDDinfor].DDfileApInodo = -1
	var NAMEEmpty [16]byte
	DD.DDarrayFiles[keyDDinfor].DDfileNombre = NAMEEmpty
	WriteDD(path, DD, PosDD)
}

//RecorrerYCrearAVD Recorre y crear la carpeta
func RecorrerYCrearAVD(path string, SB Estruct.SuperBoot, carpetas []string, IndCarp int, CrearTodo bool, AVD Estruct.AVD, posAVD int64) (Estruct.AVD, int64) {
	AVDAux := AVD
	PosAVDAux := posAVD
	if IndCarp <= len(carpetas)-1 {
		for {
			for _, value := range AVDAux.AvdApArraySubDirectorios {
				if value > 0 {
					SubDirAVD, _ := ReadAVD(path, value)
					var NAME [16]byte
					copy(NAME[:], carpetas[IndCarp])
					if NAME == SubDirAVD.AvdNombreDirectorio {
						return RecorrerYCrearAVD(path, SB, carpetas, IndCarp+1, CrearTodo, SubDirAVD, value)
					}
				}
			}
			if AVDAux.AvdApArbolVirtualDirectorio > 0 {
				PosAVDAux = AVDAux.AvdApArbolVirtualDirectorio
				AVDAux, _ = ReadAVD(path, AVDAux.AvdApArbolVirtualDirectorio)
			} else {
				break
			}
		}
		AVDAux = AVD
		PosAVDAux = posAVD
		for {
			for key, value := range AVDAux.AvdApArraySubDirectorios {
				if CrearTodo == true || len(carpetas)-1 == IndCarp {
					if value <= 0 {
						poslibre := PosicionBitmapLibre(path, SB.Sb_ap_bitmap_arbol_directorio, SB.Sb_arbol_virtual_count)
						if poslibre != -1 {
							posEscribir := SB.Sb_ap_arbol_directorio + (poslibre * SB.Sb_size_struct_arbol_directorio)
							AVDAux.AvdApArraySubDirectorios[key] = posEscribir
							AVDNuevo := Estruct.AVD{AvdProper: int64(UserLogueado.ID_user)}
							copy(AVDNuevo.AvdNombreDirectorio[:], carpetas[IndCarp])
							copy(AVDNuevo.AvdFechaCreacion[:], StringFechaActual())
							AVDNuevo.AvdApDetalleDirectorio = -1
							AVDNuevo.AvdApArbolVirtualDirectorio = -1
							WriteAVD(path, AVDNuevo, posEscribir)
							WriteAVD(path, AVDAux, PosAVDAux)
							posWriteUno := SB.Sb_ap_bitmap_arbol_directorio + poslibre
							//Se escribe uno en el bitmap
							WriteOneByteUno(path, posWriteUno)
							IndCarp++
							return RecorrerYCrearAVD(path, SB, carpetas, IndCarp, CrearTodo, AVDNuevo, posEscribir)
						} else {
							colorstring.Println("[red]\tNo hay espacio libre para escribir una carpeta")
							return Estruct.AVD{}, 0
						}
					}
				}
			}
			if AVDAux.AvdApArbolVirtualDirectorio > 0 {
				PosAVDAux = AVDAux.AvdApArbolVirtualDirectorio
				AVDAux, _ = ReadAVD(path, AVDAux.AvdApArbolVirtualDirectorio)
			} else {
				break
			}
		}
		if CrearTodo == true || len(carpetas)-1 == IndCarp {
			posFree := PosicionBitmapLibre(path, SB.Sb_ap_bitmap_arbol_directorio, SB.Sb_arbol_virtual_count)
			if posFree != -1 {
				AVDAUXSIG := Estruct.AVD{AvdApArbolVirtualDirectorio: -1}
				copy(AVDAUXSIG.AvdFechaCreacion[:], StringFechaActual())
				AVDAUXSIG.AvdNombreDirectorio = AVDAux.AvdNombreDirectorio
				AVDAUXSIG.AvdApDetalleDirectorio = -1
				AVDAUXSIG.AvdProper = AVDAux.AvdProper
				posInsertarAVDAUXSIG := SB.Sb_ap_arbol_directorio + (posFree * SB.Sb_size_struct_arbol_directorio)
				AVDAux.AvdApArbolVirtualDirectorio = posInsertarAVDAUXSIG
				WriteAVD(path, AVDAux, PosAVDAux)
				WriteAVD(path, AVDAUXSIG, posInsertarAVDAUXSIG)
				posWriteUno := SB.Sb_ap_bitmap_arbol_directorio + posFree
				WriteOneByteUno(path, posWriteUno)
				return RecorrerYCrearAVD(path, SB, carpetas, IndCarp, CrearTodo, AVDAUXSIG, posInsertarAVDAUXSIG)
			} else {
				colorstring.Println("[red]\tNo hay Espacio Suficiente para escribir la carpeta")
				return Estruct.AVD{}, 0
			}
		} else {
			colorstring.Println("[red]\tError la carpeta no existe")
			return Estruct.AVD{}, 0
		}
	} else {
		return AVD, posAVD
	}
}

//WriteSB Escribe el super bloque en la posicion partstart
func WriteSB(path string, SB Estruct.SuperBoot, partStart int64) bool {
	file2, errfile := os.OpenFile(path, os.O_WRONLY, 0775)
	if errfile != nil {
		colorstring.Println("[red]Error al escribir SuperBloque")
		return false
	}
	file2.Seek(partStart, 0)
	SBwrite := &SB
	var WriteStrucBinary bytes.Buffer
	err := binary.Write(&WriteStrucBinary, binary.BigEndian, SBwrite)
	if err != nil {
		colorstring.Println("[red]" + err.Error())
		return false
	}
	WriteByte(file2, WriteStrucBinary.Bytes())
	file2.Close()
	return true
}

//WriteOneByteCero Escribe un byte en la posicion a escribir
func WriteOneByteCero(path string, PosEscribir int64) {
	file, err := os.OpenFile(path, os.O_WRONLY, 0775)
	if err != nil {
		colorstring.Println("[red]Error al abrir el archivo")
	}
	file.Seek(PosEscribir, 0)
	var cero byte = 0
	Ins := &cero
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, Ins)
	WriteByte(file, binario.Bytes())
	file.Close()
}

//WriteOneByteUno Escribe un byte en la posicion a escribir
func WriteOneByteUno(path string, PosEscribir int64) {
	file, err := os.OpenFile(path, os.O_WRONLY, 0775)
	if err != nil {
		colorstring.Println("[red]Error al abrir el archivo")
	}
	file.Seek(PosEscribir, 0)
	var uno byte = 1
	Ins := &uno
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, Ins)
	WriteByte(file, binario.Bytes())
	file.Close()
}

//ReadSB retorna el super bloque de la posicion del partStart
func ReadSB(path string, partStart int64) (Estruct.SuperBoot, error) {
	file, errfile := os.Open(path)
	if errfile != nil {
		colorstring.Println("[red]No se logro abrir el disco y leer el super bloque")
		return Estruct.SuperBoot{}, errfile
	}
	file.Seek(partStart, 0)
	SBread := Estruct.SuperBoot{}
	A := ReadBytes(file, int(unsafe.Sizeof(SBread)))
	buffer := bytes.NewBuffer(A)
	err := binary.Read(buffer, binary.BigEndian, &SBread)
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
		return Estruct.SuperBoot{}, err
	}
	file.Close()
	return SBread, nil
}

//BuscarParticionMontada Recorre el arreglo ParticionesMontada y busca el ID
func BuscarParticionMontada(ID string) (Estruct.MountFisic, bool) {
	for I, value := range ParticionesMontada {
		if value.Id == ID {
			return ParticionesMontada[I], true
		}
	}
	return Estruct.MountFisic{}, false
}

//ReadMBR Lee el mbr del disco
func ReadMBR(path string) (Estruct.MBR, error) {
	file, err := os.Open(removeCom(path))
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al abrir el Disco " + path)
		return Estruct.MBR{}, err
	}
	mbr := Estruct.MBR{}
	A := ReadBytes(file, int(unsafe.Sizeof(mbr)))
	buffer := bytes.NewBuffer(A)
	err = binary.Read(buffer, binary.BigEndian, &mbr)
	file.Close()
	return mbr, err
}

//StringFechaActual Retorna la fecha actual
func StringFechaActual() string {
	date := time.Now().String()
	fecha := strings.Split(date, " ")[0] + " " + strings.Split(strings.Split(date, " ")[1], ":")[0] + ":" + strings.Split(strings.Split(date, " ")[1], ":")[1]
	return fecha
}

//PosicionBitmapLibre Retorna el numero por el cual se debe multiplicar el size de la estructura
func PosicionBitmapLibre(path string, PartStar int64, SizeArreglo int64) int64 {
	file, err := os.OpenFile(path, os.O_RDONLY, 0775)
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al abrir el disco")
		return -1
	}
	file.Seek(PartStar, 0)
	var Arreglo []byte
	A := ReadBytes(file, int(SizeArreglo))
	buffer := bytes.NewBuffer(A)
	binary.Read(buffer, binary.BigEndian, &Arreglo)
	for key, _ := range A {
		if int64(A[key]) == 0 {
			return int64(key)
		}
	}
	file.Close()
	return -1
}

//PosicionesBitmapLibre Retorna el numero por el cual se debe multiplicar el size de la estructura
func PosicionesBitmapLibre(path string, PartStar int64, SizeArreglo int64, NPos int) []int64 {
	var posiciones []int64
	file, err := os.OpenFile(path, os.O_RDONLY, 0775)
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al abrir el disco")
		return posiciones
	}
	file.Seek(PartStar, 0)
	var Arreglo []int8
	A := ReadBytes(file, int(SizeArreglo))
	buffer := bytes.NewBuffer(A)
	binary.Read(buffer, binary.BigEndian, &Arreglo)
	contador := 0
	for key, value := range A {
		if int(value) <= 0 {
			posiciones = append(posiciones, int64(key))
			contador++
			if contador == NPos {
				return posiciones
			}
		} else if value == 1 {

		}
	}
	file.Close()
	return posiciones
}

//PosicionesBitmapTOTALES retorna todos los bitmaps
func PosicionesBitmapTOTALES(path string, PartStar int64, SizeArreglo int64) []byte {
	var posiciones []byte
	file, err := os.OpenFile(path, os.O_RDONLY, 0775)
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al abrir el disco")
		return posiciones
	}
	file.Seek(PartStar, 0)
	var Arreglo []int8
	A := ReadBytes(file, int(SizeArreglo))
	buffer := bytes.NewBuffer(A)
	binary.Read(buffer, binary.BigEndian, &Arreglo)
	file.Close()
	posiciones = A
	return posiciones
}

//WriteAVD escribe el arbol virtual de directorio en la posicion partstart
func WriteAVD(path string, AVD Estruct.AVD, partStart int64) bool {
	file, errfile := os.OpenFile(path, os.O_WRONLY, 0775)
	if errfile != nil {
		colorstring.Println("[red]Error al escribir el Arbol Virtual de Directorio")
		return false
	}
	file.Seek(partStart, 0)
	write := &AVD
	var WriteStrucBinary bytes.Buffer
	err := binary.Write(&WriteStrucBinary, binary.BigEndian, write)
	if err != nil {
		colorstring.Println("[red]" + err.Error())
		return false
	}
	WriteByte(file, WriteStrucBinary.Bytes())
	file.Close()
	return true
}

//WriteDD escribe el detalle de directorio en la posicion partstart
func WriteDD(path string, DD Estruct.DD, partStart int64) bool {
	file, errfile := os.OpenFile(path, os.O_WRONLY, 0775)
	if errfile != nil {
		colorstring.Println("[red]Error al escribir el Detalle de directorio")
		return false
	}
	file.Seek(partStart, 0)
	write := &DD
	var WriteStrucBinary bytes.Buffer
	err := binary.Write(&WriteStrucBinary, binary.BigEndian, write)
	if err != nil {
		colorstring.Println("[red]" + err.Error())
		return false
	}
	WriteByte(file, WriteStrucBinary.Bytes())
	file.Close()
	return true
}

//WriteInodo escribe el Inodo en la posicion partstart
func WriteInodo(path string, Inodo Estruct.INODO, partStart int64) bool {
	file, errfile := os.OpenFile(path, os.O_WRONLY, 0775)
	if errfile != nil {
		colorstring.Println("[red]Error al escribir el Inodo")
		return false
	}
	file.Seek(partStart, 0)
	write := &Inodo
	var WriteStrucBinary bytes.Buffer
	err := binary.Write(&WriteStrucBinary, binary.BigEndian, write)
	if err != nil {
		colorstring.Println("[red]" + err.Error())
		return false
	}
	WriteByte(file, WriteStrucBinary.Bytes())
	file.Close()
	return true
}

//WriteBloque escribe el Bloque en la posicion partstart
func WriteBloque(path string, BD Estruct.BD, partStart int64) bool {
	file, errfile := os.OpenFile(path, os.O_WRONLY, 0775)
	if errfile != nil {
		colorstring.Println("[red]Error al escribir el Bloque")
		return false
	}
	defer file.Close()
	file.Seek(partStart, 0)
	write := &BD
	var WriteStrucBinary bytes.Buffer
	err := binary.Write(&WriteStrucBinary, binary.BigEndian, write)
	if err != nil {
		colorstring.Println("[red]\t" + err.Error())
		return false
	}
	WriteByte(file, WriteStrucBinary.Bytes())
	file.Close()
	return true
}

//WriteLog escribe el Log en la posicion partstart
func WriteLog(path string, LOG Estruct.Bitacora, partStart int64) bool {
	file, errfile := os.OpenFile(path, os.O_WRONLY, 0775)
	if errfile != nil {
		colorstring.Println("[red]\tError al escribir el Bloque")
		return false
	}
	defer file.Close()
	file.Seek(partStart, 0)
	write := &LOG
	var WriteStrucBinary bytes.Buffer
	err := binary.Write(&WriteStrucBinary, binary.BigEndian, write)
	if err != nil {
		colorstring.Println("[red]\t" + err.Error())
		return false
	}
	WriteByte(file, WriteStrucBinary.Bytes())
	file.Close()
	return true
}

//ReadLog lee el resitro de la bitacora.
func ReadLog(path string, partStart int64) (Estruct.Bitacora, error) {
	file, err := os.Open(removeCom(path))
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al abrir el Disco " + path)
		return Estruct.Bitacora{}, err
	}
	file.Seek(partStart, 0)
	estructura := Estruct.Bitacora{}
	A := ReadBytes(file, int(unsafe.Sizeof(estructura)))
	buffer := bytes.NewBuffer(A)
	err = binary.Read(buffer, binary.BigEndian, &estructura)
	file.Close()
	return estructura, nil
}

//ReadAVD lee el arbol virtual de directorio en la posicion partstart
func ReadAVD(path string, partStart int64) (Estruct.AVD, error) {
	file, err := os.Open(removeCom(path))
	if err != nil {
		colorstring.Println("[red]\tOcurrio un error al abrir el Disco " + path)
		return Estruct.AVD{}, err
	}
	file.Seek(partStart, 0)
	estructura := Estruct.AVD{}
	A := ReadBytes(file, int(unsafe.Sizeof(estructura)))
	buffer := bytes.NewBuffer(A)
	err = binary.Read(buffer, binary.BigEndian, &estructura)
	file.Close()
	return estructura, nil
}

//ReadDD lee el detalle de directorio en la posicion partstart
func ReadDD(path string, partStart int64) (Estruct.DD, error) {
	file, err := os.Open(removeCom(path))
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al abrir el Disco " + path)
		return Estruct.DD{}, err
	}
	file.Seek(partStart, 0)
	estructura := Estruct.DD{}
	A := ReadBytes(file, int(unsafe.Sizeof(estructura)))
	buffer := bytes.NewBuffer(A)
	err = binary.Read(buffer, binary.BigEndian, &estructura)
	file.Close()
	return estructura, nil
}

//ReadInodo lee el Inodo en la posicion partstart
func ReadInodo(path string, partStart int64) (Estruct.INODO, error) {
	file, err := os.Open(removeCom(path))
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al abrir el Disco " + path)
		return Estruct.INODO{}, err
	}
	file.Seek(partStart, 0)
	estructura := Estruct.INODO{}
	A := ReadBytes(file, int(unsafe.Sizeof(estructura)))
	buffer := bytes.NewBuffer(A)
	err = binary.Read(buffer, binary.BigEndian, &estructura)
	file.Close()
	return estructura, nil
}

//ReadBloque lee el Bloque en la posicion partstart
func ReadBloque(path string, partStart int64) (Estruct.BD, error) {
	file, err := os.Open(removeCom(path))
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al abrir el Disco " + path)
		return Estruct.BD{}, err
	}
	file.Seek(partStart, 0)
	estructura := Estruct.BD{}
	A := ReadBytes(file, int(unsafe.Sizeof(estructura)))
	buffer := bytes.NewBuffer(A)
	err = binary.Read(buffer, binary.BigEndian, &estructura)
	file.Close()
	return estructura, nil
}

//writeUsertxt escribe el archivo user.txt
func writeUsertxt(path string, SB Estruct.SuperBoot) {
	UserTxt := "1,G,root\n" + "1,U,root,root,201801527\n"
	//recuperar Bloques
	BloquesInsertar, PBloques, _ := ReturnBloques(path, SB, UserTxt)
	//Recuperar Inodos
	InodosInsertar, PInodos, _ := ReturnInodos(path, SB, PBloques, 1, int64(len(UserTxt)))
	//Recuperar DD
	FileInfo := Estruct.DDInfo{DDfileApInodo: PInodos[0]}
	copy(FileInfo.DDfileDateCreate[:], StringFechaActual())
	copy(FileInfo.DDfileNombre[:], "user.txt")
	DDInsertar, PDD, _ := ReturnDD(path, SB, FileInfo)
	//Recuperar carpeta
	AVD, PAVD, _ := ReturnAVD(path, 1, "/", SB, PDD)
	//Escribir Bloquess
	for key, _ := range BloquesInsertar {
		WriteBloque(path, BloquesInsertar[key], PBloques[key])
	}
	//Escribir direccion de bloques en inodos
	contbloque := 0
	for i := 0; i < len(InodosInsertar); i++ {
		for j := 0; j < 4; j++ {
			if contbloque != len(PBloques) {
				break
			}
			InodosInsertar[i].IarrayBloque[j] = PBloques[contbloque]
			contbloque++
		}
	}
	//Escribir Inodos
	for key, _ := range InodosInsertar {
		WriteInodo(path, InodosInsertar[key], PInodos[key])
	}
	//Escribir DD
	DDInsertar.DDarrayFiles[0].DDfileApInodo = PInodos[0]
	WriteDD(path, DDInsertar, PDD)
	//Escribir AVD
	AVD.AvdApDetalleDirectorio = PDD
	WriteAVD(path, AVD, PAVD)
}

//ReturnBloques retorna un arreglo con los bloques de texto
func ReturnBloques(path string, SB Estruct.SuperBoot, TEXTO string) ([]Estruct.BD, []int64, bool) {
	var Bloques []Estruct.BD
	var Posiciones []int64
	var contador int = 0
	var Texto25Car string = ""
	for i := 0; i < len(TEXTO); i++ {
		Texto25Car += string(TEXTO[i])
		contador++
		if contador == 25 {
			contador = 0
			newBloque := Estruct.BD{}
			copy(newBloque.BDData[:], Texto25Car)
			Bloques = append(Bloques, newBloque)
			Texto25Car = ""
		}
	}
	if len(Texto25Car) > 0 && len(Texto25Car) <= 25 {
		newBloque := Estruct.BD{}
		copy(newBloque.BDData[:], Texto25Car)
		Bloques = append(Bloques, newBloque)
	}
	POSFREE := PosicionesBitmapLibre(path, SB.Sb_ap_bitmap_bloque, SB.Sb_bloques_count, len(Bloques))
	if len(POSFREE) != len(Bloques) {
		return Bloques, Posiciones, false
	}
	for _, value := range POSFREE {
		posWrite := SB.Sb_ap_bitmap_bloque + value
		WriteOneByteUno(path, posWrite)
	}
	for key, _ := range POSFREE {
		posSiguiente := SB.Sb_ap_bloques + (POSFREE[key] * SB.Sb_size_struct_bloque)
		Posiciones = append(Posiciones, posSiguiente)
	}

	return Bloques, Posiciones, true
}

//ReturnInodos retorna una estructura de INodos y un arreglo con los part start donde se deben escribir
func ReturnInodos(path string, SB Estruct.SuperBoot, ApBloques []int64, IDproper int64, sizeArchivo int64) ([]Estruct.INODO, []int64, bool) {
	var InodosReturn []Estruct.INODO
	var PosInodo []int64
	contador := 0
	var contadorInodos int64 = 1
	InodoInsertar := Estruct.INODO{IapOtroInodo: -1}
	for key, _ := range InodoInsertar.IarrayBloque {
		InodoInsertar.IarrayBloque[key] = -1
	}
	for key, _ := range ApBloques {
		InodoInsertar.IarrayBloque[contador] = ApBloques[key]
		if contador == 3 || len(ApBloques)-1 == key {
			InodoInsertar.IcountInodo = contadorInodos
			InodoInsertar.IidProper = IDproper
			InodoInsertar.IsizeArchivo = sizeArchivo
			InodoInsertar.IcountBloqueAsignado = int64(contador + 1)
			InodosReturn = append(InodosReturn, InodoInsertar)
			InodoInsertar = Estruct.INODO{}
			InodoInsertar.IapOtroInodo = -1
			contadorInodos++
			contador = -1
			for key, _ := range InodoInsertar.IarrayBloque {
				InodoInsertar.IarrayBloque[key] = -1
			}
			//agregar inodo
		}
		contador++
	}
	POSFREE := PosicionesBitmapLibre(path, SB.Sb_ap_bitmap_tabla_inodo, SB.Sb_inodos_count, len(InodosReturn))
	if len(POSFREE) != len(InodosReturn) {
		return InodosReturn, PosInodo, false
	}
	for _, value := range POSFREE {
		posWrite := value + SB.Sb_ap_bitmap_tabla_inodo
		WriteOneByteUno(path, posWrite)
	}
	for key, _ := range POSFREE {
		if len(POSFREE)-1 != key {
			posSiguiente := SB.Sb_ap_tabla_inodo + (POSFREE[key+1] * SB.Sb_size_struct_inodo)
			InodosReturn[key].IapOtroInodo = posSiguiente
		}
		posSiguiente := SB.Sb_ap_tabla_inodo + (POSFREE[key] * SB.Sb_size_struct_inodo)
		PosInodo = append(PosInodo, posSiguiente)
	}
	return InodosReturn, PosInodo, true
}

//ReturnDD retorna el DD unicamente cuando no existe se puede utilzar este
func ReturnDD(path string, SB Estruct.SuperBoot, FileInfo Estruct.DDInfo) (Estruct.DD, int64, bool) {
	var ReturnDD Estruct.DD
	DDInsertar := Estruct.DD{DDapDetalleDirectorio: -1}
	DDInsertar.DDarrayFiles[0] = FileInfo
	posFree := PosicionBitmapLibre(path, SB.Sb_ap_bitmap_detalle_directorio, SB.Sb_detalle_directorio_count)
	if posFree == -1 {
		return ReturnDD, posFree, false
	}
	poswrite := SB.Sb_ap_bitmap_detalle_directorio + posFree
	WriteOneByteUno(path, poswrite)
	posicion := SB.Sb_ap_detalle_directorio + (posFree * SB.Sb_size_struct_detalle_directorio)
	ReturnDD = DDInsertar
	return ReturnDD, posicion, true
}

//ReturnAVD retorna la carpte AVD unicamente cuando no existe
func ReturnAVD(path string, IdPorper int64, nombreCarpeta string, SB Estruct.SuperBoot, ApDD int64) (Estruct.AVD, int64, bool) {
	AVD := Estruct.AVD{AvdProper: IdPorper}
	copy(AVD.AvdNombreDirectorio[:], nombreCarpeta)
	copy(AVD.AvdFechaCreacion[:], StringFechaActual())
	AVD.AvdApDetalleDirectorio = ApDD
	AVD.AvdApArbolVirtualDirectorio = -1
	posFree := PosicionBitmapLibre(path, SB.Sb_ap_bitmap_arbol_directorio, SB.Sb_arbol_virtual_count)
	if posFree == -1 {
		return AVD, 0, false
	}
	posFree = PosicionBitmapLibre(path, SB.Sb_ap_bitmap_arbol_directorio, SB.Sb_arbol_virtual_count)
	pos := SB.Sb_ap_bitmap_arbol_directorio + posFree
	WriteOneByteUno(path, pos)
	posicion := SB.Sb_ap_arbol_directorio + (posFree * SB.Sb_size_struct_arbol_directorio)
	return AVD, posicion, true
}

//MensajeConfirmacion Crea un mensaje de confirmacion y retorna true o false con respecto a la respuesta correcta establecida
func MensajeConfirmacion(texto string, CorrectAnswer string) bool {
	colorstring.Print("[yellow]\t" + texto)
	Answer := ""
	fmt.Scan(&Answer)
	println()
	if strings.ToUpper(Answer) == CorrectAnswer {
		return true
	} else {
		return false
	}
}

//EscribirEnBitacora Escribe el registro en la bitacora
func EscribirEnBitacora(path string, SB Estruct.SuperBoot, LTO string, CA int64, NOMBRE string, CONTENIDO string) bool {
	Bitacora := Estruct.Bitacora{}
	copy(Bitacora.LogFecha[:], StringFechaActual())
	copy(Bitacora.LogNombre[:], NOMBRE)
	copy(Bitacora.LogContenido[:], CONTENIDO)
	copy(Bitacora.LogTipoOperacion[:], LTO)
	Bitacora.LogTipo = CA
	Inicio := SB.Sb_ap_log
	Final := (SB.Sb_ap_log + (SB.Sb_arbol_virtual_count * int64(unsafe.Sizeof(Estruct.Bitacora{})))) - int64(unsafe.Sizeof(Estruct.Bitacora{}))
	SizeBitacora := int64(unsafe.Sizeof(Estruct.Bitacora{}))
	var i int64
	var libre [100]byte
	for i = Inicio; i <= Final; i = i + SizeBitacora {
		log, _ := ReadLog(path, i)
		if log.LogNombre == libre {
			WriteLog(path, Bitacora, i)
			return true
		}
	}
	return false
}

//RecorrerEstructuras Es el encargado de mostrar en cosola lo contenido
func RecorrerEstructuras(file *os.File, Ruta string, AVD Estruct.AVD, CARPETA bool) string {
	colorstring.Println("[red]Q para Salir")
	if CARPETA == true {
		colorstring.Println("[red]A para Seleccionar Carpeta")
	}
	AVDAux := make(map[string]int64)
	AVDAux["Q"] = -1
	AVDAux["q"] = -1
	if CARPETA == true {
		AVDAux["A"] = -2
	}
	if AVD.AvdApDetalleDirectorio > 0 && CARPETA == false {
		DD, _ := ReadDD(file.Name(), AVD.AvdApDetalleDirectorio)
		colorstring.Println("[yellow]Archivos:")
		for {
			for _, value := range DD.DDarrayFiles {
				if value.DDfileApInodo > 0 {
					NAME := ""
					for _, char := range value.DDfileNombre {
						if char != 0 {
							NAME += string(char)
						}
					}
					if NAME != "" {
						colorstring.Println("[blue]\t" + NAME)
						AVDAux[NAME] = 0
					}
				}
			}
			if DD.DDapDetalleDirectorio > 0 {
				DD, _ = ReadDD(file.Name(), DD.DDapDetalleDirectorio)
			} else {
				break
			}
		}
	}
	colorstring.Println("[yellow]Carpetas:")
	for {
		for _, value := range AVD.AvdApArraySubDirectorios {
			if value > 0 {
				avdInst, _ := ReadAVD(file.Name(), value)
				NAMECARPETA := ""
				for _, char := range avdInst.AvdNombreDirectorio {
					if char != 0 {
						NAMECARPETA += string(char)
					}
				}
				AVDAux[NAMECARPETA] = value
				colorstring.Println("[blue]\t" + NAMECARPETA)
			}
		}
		if AVD.AvdApArbolVirtualDirectorio > 0 {
			AVD, _ = ReadAVD(file.Name(), AVD.AvdApArbolVirtualDirectorio)
		} else {
			break
		}
	}
	for {
		if len(AVDAux) == 2 {
			return Ruta
		}
		colorstring.Print("[yellow]Ingrese el nombre de la Carpeta o Archivo: ")
		Escogido := ""
		fmt.Scanln(&Escogido)
		_, Existe := AVDAux[Escogido]
		if Existe == true {
			if AVDAux[Escogido] == 0 {
				if Ruta[len(Ruta)-1] == '/' {
					return Ruta + Escogido
				} else {
					return Ruta + "/" + Escogido
				}
			} else if AVDAux[Escogido] == -1 {
				return "EXIT"
			} else if AVDAux[Escogido] == -2 {
				return Ruta
			} else {
				AVD, _ = ReadAVD(file.Name(), AVDAux[Escogido])
				if Ruta[len(Ruta)-1] == '/' {
					return RecorrerEstructuras(file, Ruta+Escogido, AVD, CARPETA)
				} else {
					return RecorrerEstructuras(file, Ruta+"/"+Escogido, AVD, CARPETA)
				}
			}
		} else {
			colorstring.Println("[red]\tLa Carpeta o Archivo no existe")
		}
	}
	return Ruta
}
