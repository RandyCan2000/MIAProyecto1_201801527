package Estructuras

type EBR struct {
	Part_status byte
	Part_fit    byte
	Part_start  int
	Part_size   int
	Part_next   int
	Part_name   [16]byte
}

