package uow

import "errors"

var (
	ErrorSchemaAlreadySet = errors.New("--> schema is already set")

	//
	ErrorSchemaNotSetUEI01 = errors.New("--> schema is not set - UEI01")
	ErrorSchemaNotSetUEI02 = errors.New("--> schema is not set - UEI02")
	ErrorSchemaNotSetUEI03 = errors.New("--> schema is not set - UEI03")

	//
	ErrorSchemaCannotBeEmpty = errors.New("--> schema cannot be empty")

	//
	ErrorSchemaNotMatch    = errors.New("--> schema does not match")
	ErrorSchemaNotValid    = errors.New("--> schema is not valid")
	ErrorOnSetSchema       = errors.New("--> error on set schema")
	ErrorTransactionNotSet = errors.New("--> transaction not set")
	ErrorTransactionNotEnd = errors.New("--> transaction not end")
)
