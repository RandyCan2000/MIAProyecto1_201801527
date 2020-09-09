package Estructuras

//INODO Estructura de inodo
type INODO struct {
	IcountInodo          int64
	IsizeArchivo         int64
	IcountBloqueAsignado int64
	IarrayBloque         [4]int64
	IapOtroInodo         int64
	IidProper            int64
}
