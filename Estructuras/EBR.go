package Estructuras

//EBR estructura que usa el disco
type EBR struct {
	Part_status byte     //Indica si la particion esta activa o no
	Part_fit    byte     //Tipo de ajuste B(best) F(Fist) W(worst)
	Part_start  int64    //indica el byte de inicio de la particion
	Part_size   int64    //contiene el tamano total de la particion en bytes
	Part_next   int64    //byte en el que esta el proximo EBR -1 si no hay siguiente
	Part_name   [16]byte //Nombre de la particion
}
