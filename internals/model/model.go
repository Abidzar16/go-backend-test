package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Nasabah struct {
	gorm.Model           // Adds some metadata fields to the table
	ID         uuid.UUID `gorm:"type:uuid"` // Explicitly specify the type to be uuid
	Nama       string
	Nik        string
	No_hp      string
	Rekening   string
	Saldo      int64
}

type Transaksi struct {
	gorm.Model               // Adds some metadata fields to the table
	ID             uuid.UUID `gorm:"type:uuid"` // Explicitly specify the type to be uuid
	Rekening       string
	Kode_transaksi string
	Nominal        int64
}
