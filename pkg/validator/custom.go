package validator

import (
	"fmt"
	"regexp"
	"tm-backend-trainee-impl-clean-template/pkg/logger"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var checkCost = regexp.MustCompile(`^\d+(\.\d{1,2})?$`)

func validateCost(fl validator.FieldLevel) bool {
	return checkCost.MatchString(fl.Field().String())
}

func RegisterCustomValidators(l logger.Interface) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("ValidateCost", validateCost); err != nil {
			l.Fatal(fmt.Errorf("failed to register custom validators: %w", err))
		}
	}
}
