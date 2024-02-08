package migration

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Table struct {
	Connection  *gorm.DB
	CreateTable schema.Tabler
	RenameTable Rename
	DropTable   string
	Collate     string
	Owner       string
}

type Column struct {
	Connection    *gorm.DB
	Model         schema.Tabler
	RenameColumns []Rename
	AddColumns    []string
	DropColumns   []string
	AlterColumns  []string
	Collate       string
}

type Rename struct {
	Old string
	New string
}
