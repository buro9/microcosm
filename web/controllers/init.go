package controllers

import (
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

var decoder *form.Decoder
var validate *validator.Validate
