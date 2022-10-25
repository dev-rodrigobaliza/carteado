package services

import (
	"strings"

	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/go-playground/validator/v10"
	"github.com/gobeam/stringy"
)

// getFieldName is a utility that returns the field name from the struct namespace
func getFieldName(original string) string {
	s := strings.Split(original, ".")
	t := stringy.New(s[len(s)-1])
	r := t.SnakeCase("?", "")
	return r.ToLower()
}

// validate is a utility that validates the input data against the provided validation rules defined in the struct tags
func validate(s interface{}) []*response.ErrorValidation {
	var errors []*response.ErrorValidation
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := getFieldName(err.StructNamespace())

			var element response.ErrorValidation
			element.FailedField = fieldName
			element.Rule = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	return errors
}