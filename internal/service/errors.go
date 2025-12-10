/*
Copyright © 2025 Darko Luketic <info@icod.de>
*/
package service

import "errors"

var (
	ErrInvalidPlan  = errors.New("invalid plan ID")
	ErrInvalidAddon = errors.New("invalid addon ID")
	ErrEmptyCart    = errors.New("cart is empty")
)
