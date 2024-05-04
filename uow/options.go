package uow

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

func WithTracer(tracer *trace.Tracer) UowOption {
	return func(c *UnitOfWork) {
		c.tracer = tracer
	}
}

func WithSchema(schema string) UowOption {
	return func(c *UnitOfWork) {
		c.SetSchema(schema)
	}
}
func WithContext(ctx context.Context) UowOption {
	return func(c *UnitOfWork) {
		c.Ctx = context.WithValue(ctx, SchemaId{}, c.GetSchema())
	}
}
func WithConnection(cnx *pgxpool.Pool) UowOption {
	return func(c *UnitOfWork) {
		c.connection = cnx
	}
}
