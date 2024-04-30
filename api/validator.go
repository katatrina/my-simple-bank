package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/katatrina/my-simple-bank/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	currency := fieldLevel.Field().String()
	if currency != "" {
		return util.IsSupportedCurrency(currency)
	}

	return false
}
