package repository

import (
	"libs/uow"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

type RepositoryExample struct {
	db *pgxpool.Pool
}

func (r RepositoryExample) RepositoryName() string {
	return reflect.TypeOf(r).Name()
}
func (r *RepositoryExample) SetDB(db *pgxpool.Pool) { r.db = db }
func (r *RepositoryExample) GetDB() *pgxpool.Pool   { return r.db }
func (r *RepositoryExample) IsOnlyRead() bool       { return false }

func (r *RepositoryExample) FindOne() (*RepositoryExample, error) { return nil, nil }

func TestRepository_NewRepository(t *testing.T) {
	db := &pgxpool.Pool{}
	repo := NewRepository(db, &RepositoryExample{})
	assert.NotNil(t, repo.Queries)
	assert.NotNil(t, repo.Queries.db)
	assert.NotNil(t, repo.Queries.GetDB())
}

func TestRepository_RepositoryName(t *testing.T) {
	repo := &RepositoryExample{}
	repoName := reflect.TypeOf(RepositoryExample{}).Name()
	assert.Equal(t, repoName, repo.RepositoryName())
}

func TestRepository_SetDB(t *testing.T) {
	repo := &RepositoryExample{}
	db := &pgxpool.Pool{}
	repo.SetDB(db)
	assert.Equal(t, db, repo.db)
}

func TestRepository_GetDB(t *testing.T) {
	repo := &RepositoryExample{}
	db := &pgxpool.Pool{}
	repo.SetDB(db)
	assert.Equal(t, db, repo.GetDB())
}

func TestRepository_IsOnlyRead(t *testing.T) {
	repo := &RepositoryExample{}
	assert.False(t, repo.IsOnlyRead())
}

func TestRepositoryExample_FindOne(t *testing.T) {
	repo := &RepositoryExample{}
	itens, err := repo.FindOne()
	assert.Nil(t, itens)
	assert.Nil(t, err)
}

func TestNewRepositoryFromUoW(t *testing.T) {
	db := &pgxpool.Pool{}
	unitOfWork := uow.NewUnitOfWorkWithOptions(uow.WithConnection(db))
	repo := NewRepositoryFromUoW(unitOfWork, &RepositoryExample{})
	itens, err := repo.Queries.FindOne()
	assert.Nil(t, itens)
	assert.Nil(t, err)
}
