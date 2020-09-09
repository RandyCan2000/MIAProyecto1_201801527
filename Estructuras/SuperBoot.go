package Estructuras

//SuperBoot Estructura de super bloque
type SuperBoot struct {
	Sb_nombre_hd                         [16]byte //Sb_nombre_hd Nombre del disco duro virtual
	Sb_arbol_virtual_count               int64    //Cantidad de estructuras en el arbol virtual de directorio
	Sb_detalle_directorio_count          int64    //Cantidad de estructuras en el detalle de directorio
	Sb_inodos_count                      int64    //Cantidad de inodos en la tabla de i nodos
	Sb_bloques_count                     int64    //cantidad de bloque de datos
	Sb_arbol_virtual_free                int64    //cantidad de estructuras en el arbol virtual de dictorio libre
	Sb_detalle_directorio_free           int64    //Cantidad de estructuras en el detalle de directorio libre
	Sb_inodos_free                       int64    //cantidad de inodos en la tabla de inodos libres
	Sb_bloques_free                      int64    //Cantidad de bloque de datos libre
	Sb_date_creacion                     [16]byte //Fecha de creacion del sistema el formato dd/mm/yy hh:mm
	Sb_date_ultimo_montaje               [16]byte //Ultima fecha de montaje del sistema el formato dd/mm/yy hh:mm
	Sb_montaje_count                     int64    //contador de montajes del sistema de archivos LWH
	Sb_ap_bitmap_arbol_directorio        int64    //Apuntador al inicio del bitmap del arbol virtual de directorio
	Sb_ap_arbol_directorio               int64    //Apuntador al inicio del arbol virtual del directorio
	Sb_ap_bitmap_detalle_directorio      int64    //apuntador al inicio del bitmap de detalle de directorio
	Sb_ap_detalle_directorio             int64    //apuntador al inicio del detalle de directorio
	Sb_ap_bitmap_tabla_inodo             int64    //apuntador al inicio del birmap de la tabla de inodos
	Sb_ap_tabla_inodo                    int64    //apuntador al inicio de tabla de inodos
	Sb_ap_bitmap_bloque                  int64    //Apuntador al inicio del bitmap de bloque de datos
	Sb_ap_bloques                        int64    //apuntador al inicio del bloque de datos
	Sb_ap_log                            int64    //Apuntador al inicio del log o bitacoras
	Sb_size_struct_arbol_directorio      int64    //Tamaño de una estructura del arbol virtual de directorio
	Sb_size_struct_detalle_directorio    int64    //Tamaño de la estructura de una detalle de directorio
	Sb_size_struct_inodo                 int64    //Tamaño de la estructura de un inodo
	Sb_size_struct_bloque                int64    //Tamaño de la estructura de un bloque de datos
	Sb_first_free_bit_arbol_directorio   int64    //primer bit libre en el bitmap arbol de directorio
	Sb_first_free_bit_detalle_directorio int64    //Primero bit libre en el bitmap detalle de directorio
	Sb_first_free_bit_table_inodo        int64    //primero bit libre en el bitmap de inodo
	Sb_first_free_bit_bloque             int64    //Primer bit libre en el bitmap del bloque de datos
	Sb_magic_num                         int64    //Numero de carnet Estudiante
}
