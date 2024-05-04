package uow

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"go.opentelemetry.io/otel/sdk/trace"
)

func startDatabase(name string) (*pgxpool.Pool, testcontainers.Container) {
	DB_USER := "user"
	DB_PASS := "pass"
	DB_NAME := "db"
	DB_PORT := "30000"

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image: "postgres:16",

		ExposedPorts: []string{DB_PORT + ":5432"},
		Env: map[string]string{
			"POSTGRES_DB":       DB_NAME,
			"POSTGRES_USER":     DB_USER,
			"POSTGRES_PASSWORD": DB_PASS,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
		Name:       name,
		//  "TestUnitOfWork_Option_UsingRealDB_WithConnection",
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	host, err := container.Endpoint(ctx, "")
	if err != nil {
		panic(err)
	}

	DB_URI := "postgres://" + DB_USER + ":" + DB_PASS + "@" + host + "/" + DB_NAME
	db, err := pgxpool.New(ctx, DB_URI)
	if err != nil {
		panic(err)
	}

	time.Sleep(2 * time.Second)
	ping := db.Ping(ctx)

	if ping != nil {
		log.Panic("--> error on ping: ", ping)
	}
	return db, container
}

func TestUnitOfWork_NewUnitOfWork(t *testing.T) {
	uow := NewUnitOfWorkWithOptions(&pgxpool.Pool{})
	assert.NotNil(t, uow)
}

func TestUnitOfWork_Option_WithConnection(t *testing.T) {
	uow := NewUnitOfWorkWithOptions(&pgxpool.Pool{})
	assert.NotNil(t, uow.connection)
}

func TestUnitOfWork_Option_UsingRealDB(t *testing.T) {
	db, container := startDatabase(t.Name())
	defer container.Terminate(context.Background())

	uow := NewUnitOfWorkWithOptions(db, WithSchema("asuka"), WithContext(context.Background()))
	assert.Equal(t, "asuka", *uow.schema)
	assert.NotNil(t, uow)
	assert.NotNil(t, uow.Ctx)
	assert.NotNil(t, uow.GetConnection())
}

func TestUnitOfWork_SetSchema(t *testing.T) {
	db, container := startDatabase(t.Name())
	defer container.Terminate(context.Background())
	uow := NewUnitOfWorkWithOptions(db, WithSchema("potato"), WithContext(context.Background()))
	assert.NotNil(t, uow)
	assert.NotNil(t, uow.Ctx)
	assert.Equal(t, "potato", *uow.GetSchema())
}

func TestUnitOfWork_SetInvalidSchema(t *testing.T) {
	db, container := startDatabase(t.Name())
	defer container.Terminate(context.Background())

	uow := NewUnitOfWorkWithOptions(db)
	err := uow.SetSchema("")
	assert.NotNil(t, err)
	assert.Equal(t, ErrorSchemaCannotBeEmpty, err)
}

func TestUnitOfWork_SetSchemaTwice(t *testing.T) {
	db, container := startDatabase(t.Name())
	defer container.Terminate(context.Background())
	uow := NewUnitOfWorkWithOptions(db, WithSchema("potato"))
	assert.Equal(t, "potato", *uow.GetSchema())

	err := uow.SetSchema("tomato")
	assert.NotNil(t, err)
	assert.Equal(t, ErrorSchemaAlreadySet, err)
}

func TestUnitOfWork_GetSchema(t *testing.T) {
	db, container := startDatabase(t.Name())
	defer container.Terminate(context.Background())
	uow := NewUnitOfWorkWithOptions(db, WithSchema("potato"), WithContext(context.Background()))
	assert.Equal(t, "potato", *uow.GetSchema())
}

func TestUnitOfWork_GetConnection(t *testing.T) {
	db := &pgxpool.Pool{}
	uow := NewUnitOfWorkWithOptions(db)
	assert.NotNil(t, uow.GetConnection())
}

func TestUnitOfWork_Begin(t *testing.T) {
	db, container := startDatabase(t.Name())
	defer container.Terminate(context.Background())

	uow := NewUnitOfWorkWithOptions(db, WithSchema("public"), WithContext(context.Background()))
	err := uow.Begin()
	assert.Nil(t, err)
	assert.NotNil(t, uow.transaction)
}

func TestUnitOfWork_Commit(t *testing.T) {
	db, container := startDatabase(t.Name())
	defer container.Terminate(context.Background())

	uow := NewUnitOfWorkWithOptions(db, WithSchema("public"), WithContext(context.Background()))
	err := uow.Begin()
	assert.Nil(t, err)
	assert.NotNil(t, uow.transaction)
	assert.Nil(t, uow.Commit())
}

func TestUnitOfWork_Rollback(t *testing.T) {
	db, container := startDatabase(t.Name())
	defer container.Terminate(context.Background())

	uow := NewUnitOfWorkWithOptions(db, WithSchema("public"), WithContext(context.Background()))
	err := uow.Begin()
	assert.Nil(t, err)
	assert.NotNil(t, uow.transaction)
	assert.Nil(t, uow.Rollback())
}

func TestUnitOfWork_DoubleBegin(t *testing.T) {
	db, container := startDatabase(t.Name())
	defer container.Terminate(context.Background())

	uow := NewUnitOfWorkWithOptions(db, WithSchema("public"), WithContext(context.Background()))
	err := uow.Begin()
	assert.Nil(t, err)
	assert.NotNil(t, uow.transaction)

	err = uow.Begin()
	assert.NotNil(t, err)
	assert.NotNil(t, uow.transaction)
	assert.Equal(t, ErrorTransactionNotSet, err)
}

func TestUnitOfWork_DoubleCommit(t *testing.T) {
	db, container := startDatabase(t.Name())
	defer container.Terminate(context.Background())

	uow := NewUnitOfWorkWithOptions(db, WithSchema("public"), WithContext(context.Background()))
	err := uow.Begin()
	assert.Nil(t, err)
	assert.NotNil(t, uow.transaction)

	err = uow.Commit()
	assert.Nil(t, err)
	assert.Nil(t, uow.transaction)

	err = uow.Commit()
	assert.NotNil(t, err)
	assert.Nil(t, uow.transaction)
	assert.Equal(t, ErrorTransactionNotSet, err)
}

func TestUnitOfWork_DoubleRollback(t *testing.T) {
	db, container := startDatabase(t.Name())
	defer container.Terminate(context.Background())

	uow := NewUnitOfWorkWithOptions(db, WithSchema("public"), WithContext(context.Background()))
	err := uow.Begin()
	assert.Nil(t, err)
	assert.NotNil(t, uow.transaction)

	err = uow.Rollback()
	assert.Nil(t, err)
	assert.Nil(t, uow.transaction)

	err = uow.Rollback()
	assert.NotNil(t, err)
	assert.Nil(t, uow.transaction)
	assert.Equal(t, ErrorTransactionNotSet, err)
}

func TestUnitOfWork_BeginNoSchema(t *testing.T) {
	db, container := startDatabase(t.Name())
	defer container.Terminate(context.Background())

	uow := NewUnitOfWorkWithOptions(db, WithSchema("public"), WithContext(context.Background()))
	uow.schema = nil
	err := uow.Begin()
	assert.NotNil(t, err)
	assert.Equal(t, ErrorSchemaNotSetUEI01, err)

}

func TestUnitOfWork_Option_WithTracer(t *testing.T) {
	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
	)
	tracer := tracerProvider.Tracer("X")
	uow := NewUnitOfWorkWithOptions(&pgxpool.Pool{}, WithTracer(&tracer))
	assert.NotNil(t, uow.tracer)
	assert.NotNil(t, uow.GetTracer())
	assert.Equal(t, tracer, *uow.GetTracer())
}

func BenchmarkUnitOfWork_Option_WithConnection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewUnitOfWorkWithOptions(&pgxpool.Pool{})
	}
}
