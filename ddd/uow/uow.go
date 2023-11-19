package uow

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

type UnitOfWork struct {
	db     *gorm.DB
	tx     *gorm.DB
	schema string
}

func (uow *UnitOfWork) DB() *gorm.DB {
	return uow.db
}

func (uow *UnitOfWork) New(options ...func(*UnitOfWork)) *UnitOfWork {
	uow = &UnitOfWork{}
	for _, option := range options {
		option(uow)
	}
	return uow
}

func (uow *UnitOfWork) AddToTx(entity interface{}) error {
	if uow.tx == nil {
		return errors.New("--> Transaction not started")
	}
	return uow.tx.Create(entity).Error
}

func (uow *UnitOfWork) BeginTx() {
	if uow.schema != "" {
		uow.SetSchema(uow.schema)
	}
	uow.tx = uow.db.Begin()
}

func (uow *UnitOfWork) CommitTx() error {
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

func (uow *UnitOfWork) RollbackTx() error {
	return uow.tx.Rollback().Error
}

func (uow *UnitOfWork) MigrateDB(entidades ...any) error {
	return uow.db.AutoMigrate(entidades...)
}

func (uow *UnitOfWork) ActualSchema() string {
	sc := ""
	uow.db.Raw("SELECT current_schema();").Scan(&sc)
	return sc
}

func (uow *UnitOfWork) SetSchema(schema string) {
	uow.schema = schema
	uow.db.Exec("SET search_path TO " + schema + ";")
}

func (uow *UnitOfWork) CreateSchema(schema string) error {
	return uow.db.Exec("CREATE SCHEMA IF NOT EXISTS " + schema + ";").Error
}

func (uow *UnitOfWork) DropSchema(schema string, cascade bool) error {
	stmt := "DROP SCHEMA IF EXISTS " + schema
	if cascade {
		stmt += " CASCADE"
	}
	stmt += ";"
	return uow.db.Exec(stmt).Error
}
