package Estructuras

type Partition struct {
	Part_status byte     /*Indicar si la particion esta activa o no*/
	Part_type   byte     /*Indica el tipo de particion primario o extendida P o E*/
	Part_fit    byte     /*Tipo de ajuste de la particion tendra los valores B(best) ,f(firs),W(worst)*/
	Part_start  int64    /*Indica en que byte del disco inicia la particion*/
	Part_size   int64    /*Contiene el tamanio total de la particion en bytes*/
	Part_name   [16]byte /*Nombre de la particion*/
}
