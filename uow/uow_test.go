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
	db := &pgxpool.Pool{}
	uow := NewUnitOfWork(db)
	assert.NotNil(t, uow)
}

func TestUnitOfWork_Option_WithConnection(t *testing.T) {
	uow := NewUnitOfWorkWithOptions(WithConnection(&pgxpool.Pool{}))
	assert.NotNil(t, uow.Connection)
}

func TestUnitOfWork_Option_UsingRealDB(t *testing.T) {
	db, container := startDatabase(t.Name())
	defer container.Terminate(context.Background())

	uow := NewUnitOfWorkWithOptions(WithConnection(db), WithSchema("public"), WithContext(context.Background()))
	assert.Equal(t, "public", *uow.Schema)
	assert.NotNil(t, uow)
	assert.NotNil(t, uow.Ctx)
	assert.NotNil(t, uow.Connection)
}

func BenchmarkUnitOfWork_NewUnitOfWork(b *testing.B) {

	for i := 0; i < b.N; i++ {
		NewUnitOfWork(&pgxpool.Pool{})
	}
}

func BenchmarkUnitOfWork_Option_WithConnection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewUnitOfWorkWithOptions(WithConnection(&pgxpool.Pool{}))
	}
}
