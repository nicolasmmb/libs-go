package uow

import "gorm.io/gorm"

type UOW interface {
	DB() *gorm.DB
	New(options ...func(*UnitOfWork)) *UnitOfWork
	AddToTx(entity any) error
	BeginTx()
	CommitTx() error
	RollbackTx() error
	MigrateDB(entities ...any) error
	ActualSchema() string
	SetSchema(schema string)
	CreateSchema(schema string) error
	DropSchema(schema string, cascade bool) error
}
