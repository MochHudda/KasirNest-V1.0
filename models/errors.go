package models

import "errors"

// Common errors used across models
var (
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidQuantity   = errors.New("invalid quantity")
	ErrInvalidPrice      = errors.New("invalid price")
	ErrProductNotFound   = errors.New("product not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidUser       = errors.New("invalid user")
	ErrUnauthorized      = errors.New("unauthorized access")
)
