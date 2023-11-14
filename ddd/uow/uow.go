package uow

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

type UOW struct {
	db     *gorm.DB
	tx     *gorm.DB
	schema string
}

func (uow *UOW) DB() *gorm.DB {
	return uow.db
}

func (uow *UOW) New(options ...func(*UOW)) *UOW {
	uow = &UOW{}
	for _, option := range options {
		option(uow)
	}
	return uow
}

func (uow *UOW) AddToTx(entity interface{}) error {
	if uow.tx == nil {
		return errors.New("--> Transaction not started")
	}
	return uow.tx.Create(entity).Error
}

func (uow *UOW) BeginTx() {
	if uow.schema != "" {
		uow.SetSchema(uow.schema)
	}
	uow.tx = uow.db.Begin()
}

func (uow *UOW) CommitTx() error {
	if uow.tx == nil {
		return errors.New("--> Transaction not started - Commit")
	}

	erro := uow.tx.Commit().Error
	if erro != nil {
		log.Default().Println("--> Rollback transaction after error in commit: ")
		log.Default().Println(erro)
	}
	return erro
}

func (uow *UOW) RollbackTx() error {
	return uow.tx.Rollback().Error
}

func (uow *UOW) MigrateDB(entidades ...any) error {
	return uow.db.AutoMigrate(entidades...)
}

func (uow *UOW) ActualSchema() string {
	sc := ""
	uow.db.Raw("SELECT current_schema()").Scan(&sc)
	return sc
}

func (uow *UOW) SetSchema(schema string) {
	uow.schema = schema
	uow.db.Exec("SET search_path TO " + schema + ";")
}

func (uow *UOW) CreateSchema(schema string) error {
	return uow.db.Exec("CREATE SCHEMA IF NOT EXISTS " + schema + ";").Error
}

func (uow *UOW) DropSchema(schema string, cascade bool) error {
	stmt := "DROP SCHEMA IF EXISTS " + schema
	if cascade {
		stmt += " CASCADE" + ";"
	}
	return uow.db.Exec(stmt).Error
}
