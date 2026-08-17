package util

import "github.com/go-playground/validator/v10"

// validate 全局参数校验器（validator/v10）。
var validate = validator.New()

// ValidateStruct 校验结构体。
func ValidateStruct(v any) error {
	return validate.Struct(v)
}
