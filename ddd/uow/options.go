package uow

import "gorm.io/gorm"

func WithDB(db *gorm.DB) func(*UOW) {
	return func(uow *UOW) {
		uow.db = db
	}
}

func WithSchema(schema string) func(*UOW) {
	return func(uow *UOW) {
		uow.schema = schema
	}
}
