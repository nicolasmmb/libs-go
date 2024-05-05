package uow

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

type SchemaId struct{}

type UowOption func(c *UnitOfWork)

type UnitOfWork struct {
	Ctx        context.Context
	schema     *string
	connection *pgxpool.Pool
	tracer     trace.Tracer

	// transaction
	transaction pgx.Tx
}
