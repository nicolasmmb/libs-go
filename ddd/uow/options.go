package uow

import "gorm.io/gorm"

func WithDB(db *gorm.DB) func(*UnitOfWork) {
	return func(uow *UnitOfWork) {
		uow.db = db
	}
}

func WithSchema(schema string) func(*UnitOfWork) {
	return func(uow *UnitOfWork) {
		uow.schema = schema
	}
}
