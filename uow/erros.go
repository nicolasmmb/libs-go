package uow

import "errors"

var (
	ErrorSchemaAlreadySet = errors.New("--> schema is already set")
	ErrorSchemaNotSet     = errors.New("--> schema is not set")
	ErrorSchemaNotMatch   = errors.New("--> schema does not match")
	ErrorSchemaNotValid   = errors.New("--> schema is not valid")
	ErrorOnSetSchema      = errors.New("--> error on set schema")	
)
