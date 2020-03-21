package controller

import (
	"strconv"

	validator "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.SetTagName("validate")
}

type BaseController struct{}

func (b *BaseController) getInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}

	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (b *BaseController) GetInt(s string) int {
	num, _ := b.getInt(s)
	return num
}
