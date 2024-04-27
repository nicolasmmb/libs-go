package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	uow "github.com/niko-labs/libs-go/uow"
)

type RepositoryBase interface {
	RepositoryName() string
	SetDB(db *pgxpool.Pool)
	GetDB() *pgxpool.Pool
	IsOnlyRead() bool
}

type Repository[R RepositoryBase] struct {
	Queries R
}

func NewRepository[R RepositoryBase](session *pgxpool.Pool, repo R) *Repository[R] {
	repo.SetDB(session)
	rp := &Repository[R]{Queries: repo}
	return rp
}
func NewRepositoryFromUoW[R RepositoryBase](uow *uow.UnitOfWork, repo R) *Repository[R] {
	repo.SetDB(uow.Connection)
	rp := &Repository[R]{Queries: repo}
	return rp
}
