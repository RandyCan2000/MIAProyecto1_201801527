package Estructuras

//AVD Arbol virtual de directorio
type AVD struct {
	AvdFechaCreacion            [16]byte //fecha de creacion
	AvdNombreDirectorio         [16]byte //nombre del directorio
	AvdApArraySubDirectorios    [6]int64 //arreglo de apuntadores a los sub directorios
	AvdApDetalleDirectorio      int64    //Detalle del directorio del avd
	AvdApArbolVirtualDirectorio int64    //Apuntador al mismo tipo de directorio
	AvdProper                   int64
}
