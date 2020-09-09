package Estructuras

//DDInfo Es la estructura del arreglo de DD
type DDInfo struct {
	DDfileNombre     [16]byte
	DDfileApInodo    int64
	DDfileDateCreate [16]byte
	DDfileDateUpdate [16]byte
}

//DD Detalle de directorio
type DD struct {
	DDarrayFiles          [5]DDInfo //arreglo de la informacion de 5 archivos
	DDapDetalleDirectorio int64     //numero de byte donde esta la siguiente estructura
}
