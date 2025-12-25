package validation

import "github.com/go-playground/validator/v10"

// NewValidator creates and configures a new validator instance
func NewValidator() *validator.Validate {
	v := validator.New()
	// Register custom validations here if needed
	// Example: v.RegisterValidation("custom_tag", customValidationFunc)
	return v
}
