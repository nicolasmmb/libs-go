package uow

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

type UnitOfWorkInterface interface {
	GetSchema() *string
	SetSchema(schema string) error
	GetTracer() *trace.Tracer
	GetConnection() *pgxpool.Pool
	Begin() *pgxpool.Tx
	Commit() error
	Rollback() error
}

func (u *UnitOfWork) GetTracer() trace.Tracer {
	return u.tracer
}

func NewUnitOfWorkWithOptions(cnx *pgxpool.Pool, opts ...UowOption) *UnitOfWork {
	u := &UnitOfWork{
		connection: cnx,
		Ctx:        context.Background(),
	}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

func (u *UnitOfWork) SetSchema(schema string) error {
	if schema == "" {
		return ErrorSchemaCannotBeEmpty
	}
	if u.schema != nil {
		return ErrorSchemaAlreadySet
	}

	u.schema = &schema
	err := u.setSearchPath()
	return err
}

func (u *UnitOfWork) GetSchema() *string {
	return u.schema
}

func (u *UnitOfWork) GetConnection() *pgxpool.Pool {
	return u.connection
}

func (u *UnitOfWork) Commit() error {
	if !u.hasTransaction() {
		return ErrorTransactionNotSet
	}
	err := u.transaction.Commit(u.Ctx)
	u.transaction = nil
	return err
}

func (u *UnitOfWork) Rollback() error {
	if !u.hasTransaction() {
		return ErrorTransactionNotSet
	}
	err := u.transaction.Rollback(u.Ctx)
	if err != nil {
		u.transaction = nil
		return err
	}
	u.transaction = nil
	return nil
}

func (u *UnitOfWork) Begin() error {
	if u.hasTransaction() {
		return ErrorTransactionNotSet
	}
	if !u.hasSchema() {
		return ErrorSchemaNotSetUEI01
	}

	tx, err := u.connection.Begin(u.Ctx)
	if err != nil {
		return err
	}

	SQL_SET_SCHEMA := "SET search_path TO '" + *u.schema + "';"

	_, err = tx.Exec(u.Ctx, SQL_SET_SCHEMA)
	if err != nil {
		return err
	}
	u.transaction = tx
	return nil
}

func (u *UnitOfWork) hasTransaction() bool {
	return u.transaction != nil
}

func (u *UnitOfWork) hasSchema() bool {
	return u.schema != nil
}

func (u *UnitOfWork) setSearchPath() error {
	// ctx from ctx
	ctx := context.WithValue(u.Ctx, "schema", *u.schema)
	SQL_SET_SCHEMA := "SET search_path TO '" + *u.schema + "';"
	_, err := u.connection.Exec(ctx, SQL_SET_SCHEMA)
	return err
}
