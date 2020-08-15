package Estructuras

import (
	"time"
)

type MBR struct {
	Mbr_tamaño         int
	Mbr_fecha_creacion time.Time
	Mbr_disk_signature int
	Mbr_partition_1    Partition
	Mbr_partition_2    Partition
	Mbr_partition_3    Partition
	Mbr_partition_4    Partition
}

