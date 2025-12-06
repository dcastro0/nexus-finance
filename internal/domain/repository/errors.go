package repository

import "errors"

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidCPF        = errors.New("invalid cpf format")
	ErrWeakPassword      = errors.New("password is too weak: must have 8+ chars, 1 number, 1 uppercase and 1 special char")
)
