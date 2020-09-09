package Comandos

import (
	Estruct "Proyecto1MIA/Estructuras"
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unsafe"

	"github.com/github.com/mitchellh/colorstring"
)

var ListEBRS []string
var (
	contadorCarpetas int = 0
	contadorDD       int = 0
	ContadorInodo    int = 0
	contadorBloque   int = 0
)

//Disk Genera el disco particionado
func Disk(path string, ExtSalida string, pathSalida string) bool {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al abrir el Disco " + path)

	}
	mbr := Estruct.MBR{}
	A := ReadBytes(file, int(unsafe.Sizeof(mbr)))
	buffer := bytes.NewBuffer(A)
	err = binary.Read(buffer, binary.BigEndian, &mbr)
	if err != nil {
		colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
		return false
	}
	var partStart, SizeTotal int64 = 0, 0
	if mbr.Mbr_partition_1.Part_type == 'E' {
		partStart = mbr.Mbr_partition_1.Part_start
		SizeTotal = mbr.Mbr_partition_1.Part_size
	} else if mbr.Mbr_partition_2.Part_type == 'E' {
		partStart = mbr.Mbr_partition_2.Part_start
		SizeTotal = mbr.Mbr_partition_2.Part_size
	} else if mbr.Mbr_partition_3.Part_type == 'E' {
		partStart = mbr.Mbr_partition_3.Part_start
		SizeTotal = mbr.Mbr_partition_3.Part_size
	} else if mbr.Mbr_partition_4.Part_type == 'E' {
		partStart = mbr.Mbr_partition_4.Part_start
		SizeTotal = mbr.Mbr_partition_4.Part_size
	}

	//Leer EBR
	if partStart != 0 {
		file.Seek(partStart, 0)
		ebr := Estruct.EBR{}
		A = ReadBytes(file, int(unsafe.Sizeof(ebr)))
		buffer = bytes.NewBuffer(A)
		err = binary.Read(buffer, binary.BigEndian, &ebr)
		if err != nil {
			colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
			return false
		}
		ListEBRS = nil
		RecEBRS(file, ebr, SizeTotal, 0)
	}
	var lstpartMBR [4]Estruct.Partition
	lstpartMBR[0] = mbr.Mbr_partition_1
	lstpartMBR[1] = mbr.Mbr_partition_2
	lstpartMBR[2] = mbr.Mbr_partition_3
	lstpartMBR[3] = mbr.Mbr_partition_4
	os.MkdirAll("Assets", 0775)
	file2, errFile := os.Create("Assets/Disk.dot")
	if errFile != nil {
		log.Fatal(errFile)
	}
	file2.WriteString("digraph D {\n")
	file2.WriteString("node [fontname=\"Arial\"];\n")
	F := ""
	for _, fech := range mbr.Mbr_fecha_creacion {
		if fech != 0 {
			F = F + string(fech)
		}
	}
	var SizeUse int64 = 0
	nameDisc := strings.Split(file.Name(), "/")
	file2.WriteString("MBR [shape=record label=\"" + nameDisc[len(nameDisc)-1] + "\\n" + F + "|{MBR \\n " + strconv.Itoa(int(mbr.Mbr_tamaño)) + "|")
	for j, value := range lstpartMBR {
		SizeUse = SizeUse + int64(value.Part_size)
		if value.Part_type == 'E' {
			name := ""
			for _, char := range value.Part_name {
				if char != 0 {
					name = name + string(char)
				}
			}
			file2.WriteString("Extendida \\n" + name + "\\n" + strconv.Itoa(int(value.Part_size)) + "|{")
			for i, EBRS := range ListEBRS {
				if len(ListEBRS)-1 == i {
					file2.WriteString(string(strings.Split(EBRS, ",")[0]) + " \\n " + string(strings.Split(EBRS, ",")[1]))
				} else {
					file2.WriteString(string(strings.Split(EBRS, ",")[0]) + " \\n " + string(strings.Split(EBRS, ",")[1]) + "|")
				}
			}
			if len(lstpartMBR)-1 == j {
				file2.WriteString("}")
			} else {
				file2.WriteString("}|")
			}
		} else {
			name := ""
			for _, char := range value.Part_name {
				if char != 0 {
					name = name + string(char)
				}
			}
			if len(lstpartMBR)-1 == j {
				file2.WriteString(name + " \\n " + strconv.Itoa(int(value.Part_size)))
			} else {
				file2.WriteString(name + " \\n " + strconv.Itoa(int(value.Part_size)) + "|")
			}
		}
	}
	if mbr.Mbr_tamaño-SizeUse > 0 {
		file2.WriteString("|Libre\\n" + strconv.Itoa(int(mbr.Mbr_tamaño-SizeUse)))
	}
	file2.WriteString("}\"]\n")
	file2.WriteString("\n")
	file2.WriteString("}\n")
	file2.Close()
	Generador("Assets/Disk.dot", ExtSalida, pathSalida)
	return true
}

//TreeComplete genera un arbol completo
func TreeComplete(path string, ExtSalida string, pathSalida string, nombrePart string, GraficarTodo bool) {
	MBR, _ := ReadMBR(path)
	var Part [4]Estruct.Partition
	Part[0] = MBR.Mbr_partition_1
	Part[1] = MBR.Mbr_partition_2
	Part[2] = MBR.Mbr_partition_3
	Part[3] = MBR.Mbr_partition_4
	var InicioParticion int64 = 0
	var NAME [16]byte
	copy(NAME[:], nombrePart)
	for key, _ := range Part {
		if Part[key].Part_name == NAME {
			InicioParticion = Part[key].Part_start
			break
		}
	}
	SB, _ := ReadSB(path, InicioParticion)
	os.MkdirAll("Assets", 0775)
	file2, _ := os.Create("Assets/TreeComplete.dot")
	file2.WriteString(`digraph g{
		node [shape=plain]`)
	avdRaiz, _ := ReadAVD(path, SB.Sb_ap_arbol_directorio)
	contadorBloque = 0
	contadorCarpetas = 0
	contadorDD = 0
	ContadorInodo = 0
	CarpRec(path, file2, avdRaiz, GraficarTodo)
	file2.WriteString("}")
	Generador("Assets/TreeComplete.dot", ExtSalida, pathSalida)
}

//Bitmap Escribe un txt con el bitmap
func Bitmap(path string, partStartBitmap int64, Tamano int64, pathSalida string) bool {
	splitpath := strings.Split(pathSalida, ".")
	if splitpath[len(splitpath)-1] != "txt" {
		pathSalida += ".txt"
	}
	file, err := os.Create(pathSalida)
	if err != nil {
		colorstring.Println("[red]\tError al abrir archivo de bitmap")
		return false
	}
	BITMAP := PosicionesBitmapTOTALES(path, partStartBitmap, Tamano)
	index := 0
	val := 0
	for _, value := range BITMAP {
		val = int(value)
		file.WriteString(strconv.Itoa(val))
		index++
		if index == 20 {
			file.WriteString("\n")
			index = 0
		} else {
			file.WriteString("|")
		}
	}
	file.Close()
	AbrirFile(pathSalida)
	return true
}

//LOG Grafica el log
func LOG(path string, SB Estruct.SuperBoot, pathSalida string, Extension string) bool {
	os.MkdirAll("Assets", 0775)
	file, errFile := os.Create("Assets/Bitacora.dot")
	if errFile != nil {
		log.Fatal(errFile)
	}
	table := `<TR>
	<TD colspan="5">BITACORA</TD>
	</TR>
	<TR>
	<TD >OPERACION</TD>
	<TD >TIPO</TD>
	<TD >NOMBRE</TD>
	<TD >CONTENIDO</TD>
	<TD >FECHA</TD>
	</TR>
	`
	Row := ``
	file.WriteString("digraph D {\n")
	file.WriteString("node [shape=plain]\n")
	file.WriteString("a0 [label=<<TABLE>\n")
	file.WriteString(table)
	Inicio := SB.Sb_ap_log
	Final := (SB.Sb_ap_log + (SB.Sb_arbol_virtual_count * int64(unsafe.Sizeof(Estruct.Bitacora{})))) - int64(unsafe.Sizeof(Estruct.Bitacora{}))
	SizeBitacora := int64(unsafe.Sizeof(Estruct.Bitacora{}))
	var i int64
	var libre [16]byte
	for i = Inicio; i <= Final; i = i + SizeBitacora {
		log, _ := ReadLog(path, i)
		if log.LogNombre != libre {
			nombre := ""
			for _, char := range log.LogNombre {
				if char != 0 {
					nombre += string(char)
				}
			}
			Operacion := ""
			for _, char := range log.LogTipoOperacion {
				if char != 0 {
					Operacion += string(char)
				}
			}
			contenido := ""
			for _, char := range log.LogContenido {
				if char != 0 {
					contenido += string(char)
				}
			}
			fecha := ""
			for _, char := range log.LogFecha {
				if char != 0 {
					fecha += string(char)
				}
			}
			Tipo := ""
			if log.LogTipo == 0 {
				Tipo = "ARCHIVO"
			} else {
				Tipo = "CARPETA"
			}
			Row = `<TR>
			<TD>` + Operacion + `</TD>
			<TD>` + Tipo + `</TD>
			<TD>` + nombre + `</TD>
			<TD>` + contenido + `</TD>
			<TD>` + fecha + `</TD>
			</TR>
			`
			file.WriteString(Row)
		}
	}
	file.WriteString("</TABLE>>];\n")
	file.WriteString("\n}\n")
	file.Close()
	Generador("Assets/Bitacora.dot", Extension, pathSalida)
	return true
}

//SBGraficador Es el graficador del super boot
func SBGraficador(SB Estruct.SuperBoot, pathSalida string, Extension string) bool {
	os.MkdirAll("Assets", 0775)
	file, errFile := os.Create("Assets/SuperBoot.dot")
	if errFile != nil {
		colorstring.Println("[red]\tError Al escribir el super boot")
		return false
	}
	file.WriteString("digraph g{\n")
	file.WriteString("node [shape=plain] \n a0 [label=<<TABLE>")
	name := ""
	for _, char := range SB.Sb_nombre_hd {
		if char != 0 {
			name += string(char)
		}
	}
	file.WriteString("<TR><TD>Disco</TD><TD>" + name + "</TD></TR>")
	file.WriteString("<TR><TD>Cantidad AVD</TD><TD>" + strconv.Itoa(int(SB.Sb_arbol_virtual_count)) + "</TD></TR>")
	file.WriteString("<TR><TD>Cantidad DD</TD><TD>" + strconv.Itoa(int(SB.Sb_detalle_directorio_count)) + "</TD></TR>")
	file.WriteString("<TR><TD>Cantidad INODOS</TD><TD>" + strconv.Itoa(int(SB.Sb_inodos_count)) + "</TD></TR>")
	file.WriteString("<TR><TD>Cantidad BLOQUES</TD><TD>" + strconv.Itoa(int(SB.Sb_bloques_count)) + "</TD></TR>")
	file.WriteString("<TR><TD>AVD Libre</TD><TD>" + strconv.Itoa(int(SB.Sb_arbol_virtual_free)) + "</TD></TR>")
	file.WriteString("<TR><TD>DD Libre</TD><TD>" + strconv.Itoa(int(SB.Sb_detalle_directorio_free)) + "</TD></TR>")
	file.WriteString("<TR><TD>INODO Libre</TD><TD>" + strconv.Itoa(int(SB.Sb_inodos_free)) + "</TD></TR>")
	file.WriteString("<TR><TD>BLOQUES Libre</TD><TD>" + strconv.Itoa(int(SB.Sb_bloques_free)) + "</TD></TR>")
	fecha := ""
	for _, char := range SB.Sb_date_creacion {
		if char != 0 {
			fecha += string(char)
		}
	}
	file.WriteString("<TR><TD>Fecha Creacion</TD><TD>" + fecha + "</TD></TR>")
	fecha = ""
	for _, char := range SB.Sb_date_ultimo_montaje {
		if char != 0 {
			fecha += string(char)
		}
	}
	file.WriteString("<TR><TD>Fecha Ultimo Montaje</TD><TD>" + fecha + "</TD></TR>")
	file.WriteString("<TR><TD>Montajes</TD><TD>" + strconv.Itoa(int(SB.Sb_montaje_count)) + "</TD></TR>")
	file.WriteString("<TR><TD>Ap Bitmap AVD</TD><TD>" + strconv.Itoa(int(SB.Sb_ap_bitmap_arbol_directorio)) + "</TD></TR>")
	file.WriteString("<TR><TD>Ap Bitmap DD</TD><TD>" + strconv.Itoa(int(SB.Sb_ap_bitmap_detalle_directorio)) + "</TD></TR>")
	file.WriteString("<TR><TD>Ap Bitmap INODOS</TD><TD>" + strconv.Itoa(int(SB.Sb_ap_bitmap_tabla_inodo)) + "</TD></TR>")
	file.WriteString("<TR><TD>Ap Bitmap BLOQUES</TD><TD>" + strconv.Itoa(int(SB.Sb_ap_bitmap_bloque)) + "</TD></TR>")
	file.WriteString("<TR><TD>Ap AVD</TD><TD>" + strconv.Itoa(int(SB.Sb_ap_arbol_directorio)) + "</TD></TR>")
	file.WriteString("<TR><TD>Ap DD</TD><TD>" + strconv.Itoa(int(SB.Sb_ap_detalle_directorio)) + "</TD></TR>")
	file.WriteString("<TR><TD>Ap INODOS</TD><TD>" + strconv.Itoa(int(SB.Sb_ap_tabla_inodo)) + "</TD></TR>")
	file.WriteString("<TR><TD>Ap BLOQUES</TD><TD>" + strconv.Itoa(int(SB.Sb_ap_bloques)) + "</TD></TR>")
	file.WriteString("<TR><TD>Ap LOG</TD><TD>" + strconv.Itoa(int(SB.Sb_ap_log)) + "</TD></TR>")
	file.WriteString("<TR><TD>Size AVD</TD><TD>" + strconv.Itoa(int(SB.Sb_size_struct_arbol_directorio)) + "</TD></TR>")
	file.WriteString("<TR><TD>Size DD</TD><TD>" + strconv.Itoa(int(SB.Sb_size_struct_detalle_directorio)) + "</TD></TR>")
	file.WriteString("<TR><TD>Size INODOS</TD><TD>" + strconv.Itoa(int(SB.Sb_size_struct_inodo)) + "</TD></TR>")
	file.WriteString("<TR><TD>Size BLOQUES</TD><TD>" + strconv.Itoa(int(SB.Sb_size_struct_bloque)) + "</TD></TR>")
	file.WriteString("\n</TABLE>>];\n")
	file.WriteString("\n}\n")
	file.Close()
	Generador("Assets/SuperBoot.dot", Extension, pathSalida)
	return true
}

//TREEFILE Crea el reporte tree File
func TREEFILE(path string, SB Estruct.SuperBoot, pathSalida string, Extension string, RutaSplit []string, GraphINODOS bool) {
	os.MkdirAll("Assets", 0775)
	file2, _ := os.Create("Assets/TreeFile.dot")
	file2.WriteString(`digraph g{
		node [shape=plain]`)
	Filas := `<TR>
	<TD >1</TD>
	<TD >2</TD>
	<TD >3</TD>
	<TD >4</TD>
	<TD >5</TD>
	<TD >6</TD>
	<TD >DD</TD>
	<TD >SG</TD>
	</TR>
  
	<TR>
	<TD port="1"></TD>
	<TD port="2"></TD>
	<TD port="3"></TD>
	<TD port="4"></TD>
	<TD port="5"></TD>
	<TD port="6"></TD>
	<TD port="7" bgcolor="skyblue"></TD>
	<TD port="8" bgcolor="lightgreen"></TD>
	</TR>
	`
	nombre := "AVD" + strconv.Itoa(contadorCarpetas)
	NAME := ""
	AVD, _ := ReadAVD(path, SB.Sb_ap_arbol_directorio)
	for _, char := range AVD.AvdNombreDirectorio {
		if char != 0 {
			NAME += string(char)
		}
	}
	file2.WriteString(nombre + " [label=<\n")
	file2.WriteString("<TABLE BGCOLOR=\"#E3EE67\">\n")
	file2.WriteString("<TR><TD colspan=\"8\">" + NAME + "</TD></TR>\n")
	file2.WriteString(Filas)
	file2.WriteString("</TABLE>\n")
	file2.WriteString(">];\n")
	contadorCarpetas++
	NOMBRE := RecTreeFile(path, file2, AVD, RutaSplit, 1, GraphINODOS)
	if NOMBRE[0] == 'D' {
		file2.WriteString(nombre + ":7->" + NOMBRE)
	} else {
		file2.WriteString(nombre + ":1->" + NOMBRE)
	}
	file2.WriteString("}")
	Generador("Assets/TreeFile.dot", Extension, pathSalida)
	file2.Close()
}

//TREEDIRECTORIOS Crea el reporte de TreeDirectorio
func TREEDIRECTORIOS(path string, SB Estruct.SuperBoot, pathSalida string, Extension string, RutaSplit []string) {
	os.MkdirAll("Assets", 0775)
	file2, _ := os.Create("Assets/TreeDirectorio.dot")
	file2.WriteString(`digraph g{
		node [shape=plain]`)
	Filas := `<TR>
	<TD >1</TD>
	<TD >2</TD>
	<TD >3</TD>
	<TD >4</TD>
	<TD >5</TD>
	<TD >6</TD>
	<TD >DD</TD>
	<TD >SG</TD>
	</TR>
  
	<TR>
	<TD port="1"></TD>
	<TD port="2"></TD>
	<TD port="3"></TD>
	<TD port="4"></TD>
	<TD port="5"></TD>
	<TD port="6"></TD>
	<TD port="7" bgcolor="skyblue"></TD>
	<TD port="8" bgcolor="lightgreen"></TD>
	</TR>
	`
	nombre := "AVD" + strconv.Itoa(contadorCarpetas)
	NAME := ""
	AVD, _ := ReadAVD(path, SB.Sb_ap_arbol_directorio)
	for _, char := range AVD.AvdNombreDirectorio {
		if char != 0 {
			NAME += string(char)
		}
	}
	file2.WriteString(nombre + " [label=<\n")
	file2.WriteString("<TABLE BGCOLOR=\"#E3EE67\">\n")
	file2.WriteString("<TR><TD colspan=\"8\">" + NAME + "</TD></TR>\n")
	file2.WriteString(Filas)
	file2.WriteString("</TABLE>\n")
	file2.WriteString(">];\n")
	contadorCarpetas++
	NOMBRE := RecTreeDirectorio(path, file2, AVD, RutaSplit, 1, true)
	if NOMBRE[0] == 'D' {
		file2.WriteString(nombre + ":7->" + NOMBRE)
	} else {
		file2.WriteString(nombre + ":1->" + NOMBRE)
	}
	file2.WriteString("}")
	Generador("Assets/TreeDirectorio.dot", Extension, pathSalida)
	file2.Close()
}

//RecTreeFile metodo recursivo que genera el reporte treefile
func RecTreeFile(path string, file *os.File, AVD Estruct.AVD, RutaSplit []string, indCarp int, GraphINODOS bool) string {
	nombre := "AVD" + strconv.Itoa(contadorCarpetas)
	Filas := `<TR>
	<TD >1</TD>
	<TD >2</TD>
	<TD >3</TD>
	<TD >4</TD>
	<TD >5</TD>
	<TD >6</TD>
	<TD >DD</TD>
	<TD >SG</TD>
	</TR>
  
	<TR>
	<TD port="1"></TD>
	<TD port="2"></TD>
	<TD port="3"></TD>
	<TD port="4"></TD>
	<TD port="5"></TD>
	<TD port="6"></TD>
	<TD port="7" bgcolor="skyblue"></TD>
	<TD port="8" bgcolor="lightgreen"></TD>
	</TR>
	`
	AVDAux := Estruct.AVD{}
	NAME := ""
	if indCarp == len(RutaSplit)-1 {
		if AVD.AvdApDetalleDirectorio > 0 {
			DD, _ := ReadDD(path, AVD.AvdApDetalleDirectorio)
			return DDRec(path, file, DD, true, RutaSplit[indCarp], true)
		} else {
			return "DD"
		}
	} else {
		for {
			for key, value := range AVD.AvdApArraySubDirectorios {
				if value > 0 {
					AVDAux, _ = ReadAVD(path, value)
					NAME = ""
					for _, char := range AVDAux.AvdNombreDirectorio {
						if char != 0 {
							NAME += string(char)
						}
					}
					if NAME == RutaSplit[indCarp] {
						file.WriteString(nombre + " [label=<\n")
						file.WriteString("<TABLE BGCOLOR=\"#E3EE67\">\n")
						file.WriteString("<TR><TD colspan=\"8\">" + NAME + "</TD></TR>\n")
						file.WriteString(Filas)
						file.WriteString("</TABLE>\n")
						file.WriteString(">];\n")
						AVD, _ = ReadAVD(path, value)
						Relacion := ""
						if indCarp == len(RutaSplit)-2 {
							Relacion = nombre + ":7->"
						} else {
							Relacion = nombre + ":" + strconv.Itoa(key+1) + "->"
						}
						contadorCarpetas++
						NAMEApunt := RecTreeFile(path, file, AVD, RutaSplit, indCarp+1, GraphINODOS)
						file.WriteString(Relacion + NAMEApunt + "\n")
						return nombre
					}
				}
			}
			if AVD.AvdApArbolVirtualDirectorio > 0 {
				AVD, _ = ReadAVD(path, AVD.AvdApArbolVirtualDirectorio)
			} else {
				break
			}
		}
	}
	return "DD"
}

//RecTreeDirectorio metodo recursivo que genear el reporte treeDirectorio
func RecTreeDirectorio(path string, file *os.File, AVD Estruct.AVD, RutaSplit []string, indCarp int, GraphINODOS bool) string {
	nombre := "AVD" + strconv.Itoa(contadorCarpetas)
	Filas := `<TR>
	<TD >1</TD>
	<TD >2</TD>
	<TD >3</TD>
	<TD >4</TD>
	<TD >5</TD>
	<TD >6</TD>
	<TD >DD</TD>
	<TD >SG</TD>
	</TR>
  
	<TR>
	<TD port="1"></TD>
	<TD port="2"></TD>
	<TD port="3"></TD>
	<TD port="4"></TD>
	<TD port="5"></TD>
	<TD port="6"></TD>
	<TD port="7" bgcolor="skyblue"></TD>
	<TD port="8" bgcolor="lightgreen"></TD>
	</TR>
	`
	AVDAux := Estruct.AVD{}
	NAME := ""

	for {
		for key, value := range AVD.AvdApArraySubDirectorios {
			if value > 0 {
				AVDAux, _ = ReadAVD(path, value)
				NAME = ""
				for _, char := range AVDAux.AvdNombreDirectorio {
					if char != 0 {
						NAME += string(char)
					}
				}
				if RutaSplit[indCarp] == "" && indCarp == len(RutaSplit)-1 {
					if AVD.AvdApDetalleDirectorio > 0 {
						DD, _ := ReadDD(path, AVD.AvdApDetalleDirectorio)
						NOMBRE := DDRec(path, file, DD, true, "", false)
						return NOMBRE
					}
				}
				if indCarp == len(RutaSplit)-1 && NAME == RutaSplit[indCarp] {
					file.WriteString(nombre + " [label=<\n")
					file.WriteString("<TABLE BGCOLOR=\"#E3EE67\">\n")
					file.WriteString("<TR><TD colspan=\"8\">" + NAME + "</TD></TR>\n")
					file.WriteString(Filas)
					file.WriteString("</TABLE>\n")
					file.WriteString(">];\n")
					if AVDAux.AvdApDetalleDirectorio > 0 {
						DD, _ := ReadDD(path, AVDAux.AvdApDetalleDirectorio)
						NOMBRE := DDRec(path, file, DD, true, "", false)
						file.WriteString(nombre + ":7->" + NOMBRE + "\n")
					}
					return nombre
				}
				if NAME == RutaSplit[indCarp] {
					file.WriteString(nombre + " [label=<\n")
					file.WriteString("<TABLE BGCOLOR=\"#E3EE67\">\n")
					file.WriteString("<TR><TD colspan=\"8\">" + NAME + "</TD></TR>\n")
					file.WriteString(Filas)
					file.WriteString("</TABLE>\n")
					file.WriteString(">];\n")
					AVD, _ = ReadAVD(path, value)
					Relacion := ""
					Relacion = nombre + ":" + strconv.Itoa(key+1) + "->"
					contadorCarpetas++
					NAMEApunt := RecTreeDirectorio(path, file, AVD, RutaSplit, indCarp+1, GraphINODOS)
					file.WriteString(Relacion + NAMEApunt + "\n")
					return nombre
				}
			}
		}
		if AVD.AvdApArbolVirtualDirectorio > 0 {
			AVD, _ = ReadAVD(path, AVD.AvdApArbolVirtualDirectorio)
		} else {
			break
		}
	}
	return "DD"
}

//CarpRec Escribe la estructura de carpetas retorna el nombre de la tabla AVD
func CarpRec(path string, file *os.File, Carpeta Estruct.AVD, GraficarDD bool) string {
	nombre := "AVD" + strconv.Itoa(contadorCarpetas)
	var SC1, SC2, SC3, SC4, SC5, SC6 string = "", "", "", "", "", ""
	if Carpeta.AvdApArraySubDirectorios[0] > 0 {
		SubCarp, _ := ReadAVD(path, Carpeta.AvdApArraySubDirectorios[0])
		contadorCarpetas++
		SC1 = CarpRec(path, file, SubCarp, GraficarDD)
	}
	if Carpeta.AvdApArraySubDirectorios[1] > 0 {
		SubCarp, _ := ReadAVD(path, Carpeta.AvdApArraySubDirectorios[1])
		contadorCarpetas++
		SC2 = CarpRec(path, file, SubCarp, GraficarDD)
	}
	if Carpeta.AvdApArraySubDirectorios[2] > 0 {
		SubCarp, _ := ReadAVD(path, Carpeta.AvdApArraySubDirectorios[2])
		contadorCarpetas++
		SC3 = CarpRec(path, file, SubCarp, GraficarDD)
	}
	if Carpeta.AvdApArraySubDirectorios[3] > 0 {
		SubCarp, _ := ReadAVD(path, Carpeta.AvdApArraySubDirectorios[3])
		contadorCarpetas++
		SC4 = CarpRec(path, file, SubCarp, GraficarDD)
	}
	if Carpeta.AvdApArraySubDirectorios[4] > 0 {
		SubCarp, _ := ReadAVD(path, Carpeta.AvdApArraySubDirectorios[4])
		contadorCarpetas++
		SC5 = CarpRec(path, file, SubCarp, GraficarDD)
	}
	if Carpeta.AvdApArraySubDirectorios[5] > 0 {
		SubCarp, _ := ReadAVD(path, Carpeta.AvdApArraySubDirectorios[5])
		contadorCarpetas++
		SC6 = CarpRec(path, file, SubCarp, GraficarDD)
	}
	Sig := ""
	if Carpeta.AvdApArbolVirtualDirectorio > 0 {
		SubCarp, _ := ReadAVD(path, Carpeta.AvdApArbolVirtualDirectorio)
		contadorCarpetas++
		Sig = CarpRec(path, file, SubCarp, GraficarDD)
	}
	DD := ""
	if GraficarDD == true {
		if Carpeta.AvdApDetalleDirectorio > 0 {
			DDsig, _ := ReadDD(path, Carpeta.AvdApDetalleDirectorio)
			contadorDD++
			DD = DDRec(path, file, DDsig, GraficarDD, "", false)
		}
	}
	Filas := `<TR>
	<TD >1</TD>
	<TD >2</TD>
	<TD >3</TD>
	<TD >4</TD>
	<TD >5</TD>
	<TD >6</TD>
	<TD >DD</TD>
	<TD >SG</TD>
	</TR>
  
	<TR>
	<TD port="1"></TD>
	<TD port="2"></TD>
	<TD port="3"></TD>
	<TD port="4"></TD>
	<TD port="5"></TD>
	<TD port="6"></TD>
	<TD port="7" bgcolor="skyblue"></TD>
	<TD port="8" bgcolor="lightgreen"></TD>
	</TR>\n`
	NameDirectorio := ""
	for _, char := range Carpeta.AvdNombreDirectorio {
		if char != 0 {
			NameDirectorio += string(char)
		}
	}
	file.WriteString(nombre + " [label=<\n")
	file.WriteString("<TABLE BGCOLOR=\"#E3EE67\">\n")
	file.WriteString("<TR><TD colspan=\"8\">" + NameDirectorio + "</TD></TR>\n")
	file.WriteString(Filas)
	file.WriteString("</TABLE>\n")
	file.WriteString(">];\n")
	if SC1 != "" {
		file.WriteString(nombre + ":1->" + SC1 + "\n")
	}
	if SC2 != "" {
		file.WriteString(nombre + ":2->" + SC2 + "\n")
	}
	if SC3 != "" {
		file.WriteString(nombre + ":3->" + SC3 + "\n")
	}
	if SC4 != "" {
		file.WriteString(nombre + ":4->" + SC4 + "\n")
	}
	if SC5 != "" {
		file.WriteString(nombre + ":5->" + SC5 + "\n")
	}
	if SC6 != "" {
		file.WriteString(nombre + ":6->" + SC6 + "\n")
	}
	if DD != "" {
		file.WriteString(nombre + ":7->" + DD + "\n")
	}
	if Sig != "" {
		file.WriteString(nombre + ":8->" + Sig + "\n")
	}
	return nombre
}

//DDRec escribe la estructura dd en el archivo y retorna el nombre de la tabla
func DDRec(path string, file *os.File, DD Estruct.DD, GraficarInodo bool, NombreIgualar string, SoloUno bool) string {
	Nombre := "DD" + strconv.Itoa(contadorDD)
	var DD1, DD2, DD3, DD4, DD5 string = "", "", "", "", ""
	var NDD1, NDD2, NDD3, NDD4, NDD5 string = "", "", "", "", ""
	if GraficarInodo == true {
		if DD.DDarrayFiles[0].DDfileApInodo > 0 {
			for _, char := range DD.DDarrayFiles[0].DDfileNombre {
				if char != 0 {
					NDD1 += string(char)
				}
			}
			if SoloUno == true && NDD1 != NombreIgualar {
				NDD1 = ""
			}
			if NDD1 != "" {
				Inodo, _ := ReadInodo(path, DD.DDarrayFiles[0].DDfileApInodo)
				ContadorInodo++
				DD1 = InodoRec(path, file, Inodo, GraficarInodo)
			}
		}
		if DD.DDarrayFiles[1].DDfileApInodo > 0 {
			for _, char := range DD.DDarrayFiles[1].DDfileNombre {
				if char != 0 {
					NDD2 += string(char)
				}
			}
			if SoloUno == true && NDD2 != NombreIgualar {
				NDD2 = ""
			}
			if NDD2 != "" {
				Inodo, _ := ReadInodo(path, DD.DDarrayFiles[1].DDfileApInodo)
				ContadorInodo++
				DD2 = InodoRec(path, file, Inodo, GraficarInodo)
			}

		}
		if DD.DDarrayFiles[2].DDfileApInodo > 0 {
			for _, char := range DD.DDarrayFiles[2].DDfileNombre {
				if char != 0 {
					NDD3 += string(char)
				}
			}
			if SoloUno == true && NDD3 != NombreIgualar {
				NDD3 = ""
			}
			if NDD3 != "" {
				Inodo, _ := ReadInodo(path, DD.DDarrayFiles[2].DDfileApInodo)
				ContadorInodo++
				DD3 = InodoRec(path, file, Inodo, GraficarInodo)
			}
		}
		if DD.DDarrayFiles[3].DDfileApInodo > 0 {
			for _, char := range DD.DDarrayFiles[3].DDfileNombre {
				if char != 0 {
					NDD4 += string(char)
				}
			}
			if SoloUno == true && NDD4 != NombreIgualar {
				NDD4 = ""
			}
			if NDD4 != "" {
				Inodo, _ := ReadInodo(path, DD.DDarrayFiles[3].DDfileApInodo)
				ContadorInodo++
				DD4 = InodoRec(path, file, Inodo, GraficarInodo)
			}
		}
		if DD.DDarrayFiles[4].DDfileApInodo > 0 {
			for _, char := range DD.DDarrayFiles[4].DDfileNombre {
				if char != 0 {
					NDD5 += string(char)
				}
			}
			if SoloUno == true && NDD5 != NombreIgualar {
				NDD5 = ""
			}
			if NDD5 != "" {
				Inodo, _ := ReadInodo(path, DD.DDarrayFiles[4].DDfileApInodo)
				ContadorInodo++
				DD5 = InodoRec(path, file, Inodo, GraficarInodo)
			}
		}
	}

	Sig := ""
	if DD.DDapDetalleDirectorio > 0 {
		SigDD, _ := ReadDD(path, DD.DDapDetalleDirectorio)
		contadorDD++
		Sig = DDRec(path, file, SigDD, GraficarInodo, NombreIgualar, SoloUno)
	}
	Table := `<TR><TD >1</TD><TD port="1">` + NDD1 + `</TD></TR>
	<TR><TD >2</TD><TD port="2">` + NDD2 + `</TD></TR>
	<TR><TD >3</TD><TD port="3">` + NDD3 + `</TD></TR>
	<TR><TD >4</TD><TD port="4">` + NDD4 + `</TD></TR>
	<TR><TD >5</TD><TD port="5">` + NDD5 + `</TD></TR>
	<TR><TD >SG</TD><TD port="6" bgcolor="lightgreen"></TD></TR>\n`
	file.WriteString(Nombre + " [label=<\n")
	file.WriteString("<TABLE BGCOLOR=\"#EEAD67\">\n")
	file.WriteString("<TR><TD colspan=\"2\">DD</TD></TR>\n")
	file.WriteString(Table)
	file.WriteString("</TABLE>\n")
	file.WriteString(">];\n")
	if Sig != "" {
		file.WriteString(Nombre + ":6->" + Sig + "\n")
	}
	if GraficarInodo == true {
		if DD1 != "" {
			file.WriteString(Nombre + ":1->" + DD1 + "\n")
		}
		if DD2 != "" {
			file.WriteString(Nombre + ":2->" + DD2 + "\n")
		}
		if DD3 != "" {
			file.WriteString(Nombre + ":3->" + DD3 + "\n")
		}
		if DD4 != "" {
			file.WriteString(Nombre + ":4->" + DD4 + "\n")
		}
		if DD5 != "" {
			file.WriteString(Nombre + ":5->" + DD5 + "\n")
		}
	}
	return Nombre
}

//InodosRec escribe la estructura Inodo en el archivo y retorna el nombre de la tabla
func InodoRec(path string, file *os.File, Inodo Estruct.INODO, GraficarBloque bool) string {
	Nombre := "Inodo" + strconv.Itoa(ContadorInodo)
	var B1, B2, B3, B4, SG string = "", "", "", "", ""
	if GraficarBloque == true {
		if Inodo.IarrayBloque[0] > 0 {
			Bloque, _ := ReadBloque(path, Inodo.IarrayBloque[0])
			contadorBloque++
			B1 = BloqueGraficar(file, Bloque)
		}
		if Inodo.IarrayBloque[1] > 0 {
			Bloque, _ := ReadBloque(path, Inodo.IarrayBloque[1])
			contadorBloque++
			B2 = BloqueGraficar(file, Bloque)
		}
		if Inodo.IarrayBloque[2] > 0 {
			Bloque, _ := ReadBloque(path, Inodo.IarrayBloque[2])
			contadorBloque++
			B3 = BloqueGraficar(file, Bloque)
		}
		if Inodo.IarrayBloque[3] > 0 {
			Bloque, _ := ReadBloque(path, Inodo.IarrayBloque[3])
			contadorBloque++
			B4 = BloqueGraficar(file, Bloque)
		}
	}
	if Inodo.IapOtroInodo > 0 {
		InodoSG, _ := ReadInodo(path, Inodo.IapOtroInodo)
		ContadorInodo++
		SG = InodoRec(path, file, InodoSG, GraficarBloque)
	}
	Table := `<TR><TD >1</TD><TD port="1"></TD></TR>
	<TR><TD >2</TD><TD port="2"></TD></TR>
	<TR><TD >3</TD><TD port="3"></TD></TR>
	<TR><TD >4</TD><TD port="4"></TD></TR>
	<TR><TD >SG</TD><TD port="5" bgcolor="lightgreen"></TD></TR>\n`
	file.WriteString(Nombre + " [label=<<TABLE BGCOLOR=\"#F99381\"><TR><TD colspan=\"2\">INODO</TD></TR>\n")
	file.WriteString(Table)
	file.WriteString("</TABLE>>];\n")
	if B1 != "" {
		file.WriteString(Nombre + ":1->" + B1 + "\n")
	}
	if B2 != "" {
		file.WriteString(Nombre + ":2->" + B2 + "\n")
	}
	if B3 != "" {
		file.WriteString(Nombre + ":3->" + B3 + "\n")
	}
	if B4 != "" {
		file.WriteString(Nombre + ":4->" + B4 + "\n")
	}
	if SG != "" {
		file.WriteString(Nombre + ":5->" + SG + "\n")
	}
	return Nombre
}

//BloqueGraficar escribe la estructura de bloque y retorna el nombre de la tabla
func BloqueGraficar(file *os.File, Bloque Estruct.BD) string {
	Nombre := "BLOQUE" + strconv.Itoa(contadorBloque)
	TextoBloque := ""
	for _, char := range Bloque.BDData {
		if char != 0 {
			TextoBloque += string(char)
		}
	}
	file.WriteString(Nombre + " [label=<<TABLE><TR><TD bgcolor=\"skyblue\">" + TextoBloque + "</TD></TR></TABLE>>];\n")
	return Nombre
}

//name,porcentaje
func RecEBRS(file *os.File, ebr Estruct.EBR, SizeTotal int64, SumSizes int64) {
	if ebr.Part_status == 'A' {
		name := ""
		for _, char := range ebr.Part_name {
			if char != 0 {
				name = name + string(char)
			}
		}
		ListEBRS = append(ListEBRS, name+","+strconv.Itoa(int(ebr.Part_size)))
	} else {
		if ebr.Part_next == -1 {
			if SumSizes > SizeTotal {
				ListEBRS = append(ListEBRS, "Libre,0")
			} else {
				ListEBRS = append(ListEBRS, "Libre,"+strconv.Itoa(int(SizeTotal-SumSizes)))
			}

		} else {
			ListEBRS = append(ListEBRS, "Libre,"+strconv.Itoa(int(ebr.Part_size)))
		}
	}
	if ebr.Part_next != -1 {
		file.Seek(ebr.Part_next, 0)
		ebr2 := Estruct.EBR{}
		A := ReadBytes(file, int(unsafe.Sizeof(ebr2)))
		buffer := bytes.NewBuffer(A)
		err := binary.Read(buffer, binary.BigEndian, &ebr2)
		if err != nil {
			colorstring.Println("[red]Ocurrio un error al leer el Disco " + err.Error())
		}
		RecEBRS(file, ebr2, SizeTotal, SumSizes+ebr.Part_size)
	}
}

//Generador Es el generador de la grafica
func Generador(pathdot string, Extension string, pathSalida string) {
	app := "dot"
	arg0 := "-T" + Extension
	arg1 := pathdot
	arg2 := "-o"
	arg3 := pathSalida
	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	err := cmd.Run()
	if err != nil {
		colorstring.Println("[red]\tOcurrio error al generar el reporte")
	} else {
		colorstring.Println("[green]\tSe genero con exito el reporte")
		AbrirFile(pathSalida)
	}
}

//AbrirFile abre en su visualizador predeterminado cualquier archivo
func AbrirFile(pathSalida string) {
	app := "xdg-open"
	arg0 := pathSalida
	cmd := exec.Command(app, arg0)
	err := cmd.Run()
	if err != nil {
		colorstring.Println("[red]\tOcurrio error al abrir el reporte")
	}
}

func floattostr(input_num float64) string {

	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'g', 1, 64)
}
