package Estructuras

//Bitacora Estructura de bitacora
type Bitacora struct {
	LogTipoOperacion [16]byte //El tipo de Operacion a Realizarse
	LogTipo          int64
	LogNombre        [100]byte
	LogContenido     [100]byte
	LogFecha         [16]byte
}
