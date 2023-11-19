package types

import (
	"errors"
	"strings"
)

type CPF string

func (cpf CPF) IsValid() error {
	chars := []string{"-", ".", "/", " "}
	for _, char := range chars {
		cpf = CPF(strings.ReplaceAll(string(cpf), char, ""))
	}

	if len(cpf) != 11 {
		return errors.New("--> CPF must have 11 digits")
	}

	// if cpf has only one number repeated
	countRepeated := 0
	for i := 0; i < len(cpf); i++ {
		if cpf[0] == cpf[i] {
			countRepeated++
		}
	}
	if countRepeated == len(cpf) {
		return errors.New("--> CPF has only one number repeated")
	}

	// if cpf is a sequence of numbers
	var sum int
	for i := 0; i < 9; i++ {
		sum += int(cpf[i]) * (10 - i)
	}
	rest := sum % 11
	if rest < 2 {
		rest = 0
	} else {
		rest = 11 - rest
	}

	// '0' is used to convert the byte to int
	if rest != int(cpf[9]-'0') {
		return errors.New("--> The First Verification Digit is invalid")
	}
	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(cpf[i]-'0') * (11 - i)
	}
	rest = sum % 11
	if rest < 2 {
		rest = 0
	} else {
		rest = 11 - rest
	}

	// '0' is used to convert the byte to int
	if rest != int(cpf[10]-'0') {
		return errors.New("--> The Second Verification Digit is invalid")
	}
	return nil
}
