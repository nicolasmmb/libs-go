package uow

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SchemaId struct{}

type UnitOfWorkInterface interface {
	SetSchema(schema string)
}

type UowOption func(c *UnitOfWork)

type UnitOfWork struct {
	Ctx        context.Context
	Schema     *string
	Connection *pgxpool.Pool
}

func WithSchema(schema string) UowOption {
	return func(c *UnitOfWork) {
		c.SetSchema(schema)
	}
}
func WithContext(ctx context.Context) UowOption {
	return func(c *UnitOfWork) {
		c.Ctx = ctx
	}
}
func WithConnection(cnx *pgxpool.Pool) UowOption {
	return func(c *UnitOfWork) {
		c.Connection = cnx
	}
}

func NewUnitOfWorkWithOptions(opts ...UowOption) *UnitOfWork {
	u := &UnitOfWork{}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

func NewUnitOfWork(cnx *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{
		Ctx:        context.Background(),
		Connection: cnx,
	}
}

func (u *UnitOfWork) SetSchema(schema string) {
	if schema == "" {
		log.Panicln(ErrorOnSetSchema, "Schema cannot be empty: ", schema)
	}
	if u.Connection == nil {
		log.Panicln(ErrorOnSetSchema, "Connection is nil")
	}
	if u.Ctx == nil {
		u.Ctx = context.Background()
	}
	if u.Schema != nil {
		log.Panicln(ErrorOnSetSchema, "Schema already set: ", *u.Schema)
	}

	u.Schema = &schema
	_SQL_SET_SCHEMA := "SET search_path TO '" + schema + "';"
	_, err := u.Connection.Exec(u.Ctx, _SQL_SET_SCHEMA)
	if err != nil {
		log.Panicln(ErrorOnSetSchema, err)
	}
}
