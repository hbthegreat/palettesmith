package validation

import (
	"fmt"
	"palettesmith/internal/color"
	"palettesmith/internal/plugin"
	"strings"
)

type PluginValidator struct {
	colorConverter *color.Converter
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func NewPluginValidator() *PluginValidator {
	return &PluginValidator{}
}

func (pv *PluginValidator) ValidateField(field plugin.Field, value string) []FieldError {
	var errors []FieldError

	if strings.TrimSpace(value) == "" && field.Default == "" {
		errors = append(errors, FieldError{
			Field:   field.Key,
			Message: "Field is required",
			Code:    "required",
		})
		return errors
	}

	switch field.Type {
	case "color":
		if err := pv.validateColor(field, value); err != nil {
			errors = append(errors, *err)
		}
	case "number":
		if err := pv.validateNumber(field, value); err != nil {
			errors = append(errors, *err)
		}
	case "select":
		if err := pv.validateSelect(field, value); err != nil {
			errors = append(errors, *err)
		}
	}

	return errors
}

func (pv *PluginValidator) ValidatePlugin(p plugin.Plugin) []error {
	var errors []error

	if p.Manifest.ID == "" {
		errors = append(errors, fmt.Errorf("plugin ID is required"))
	}

	if p.Manifest.Title == "" {
		errors = append(errors, fmt.Errorf("plugin title is required"))
	}

	if p.Manifest.SpecRelPath == "" {
		errors = append(errors, fmt.Errorf("spec path is required"))
	}

	if p.Spec.ID == "" {
		errors = append(errors, fmt.Errorf("spec ID is required"))
	}

	for _, field := range p.Spec.Fields {
		if field.Key == "" {
			errors = append(errors, fmt.Errorf("field key is required"))
		}
		if field.Type == "" {
			errors = append(errors, fmt.Errorf("field %s: type is required", field.Key))
		}
		if !pv.isValidFieldType(field.Type) {
			errors = append(errors, fmt.Errorf("field %s: invalid type '%s'", field.Key, field.Type))
		}
	}

	return errors
}

func (pv *PluginValidator) validateColor(field plugin.Field, value string) *FieldError {
	converter, err := color.NewConverter(value)
	if err != nil {
		return &FieldError{
			Field:   field.Key,
			Message: fmt.Sprintf("Invalid color format: %s", err.Error()),
			Code:    "invalid_color",
		}
	}

	if converter == nil {
		return &FieldError{
			Field:   field.Key,
			Message: "Color conversion failed",
			Code:    "color_conversion_failed",
		}
	}

	return nil
}

func (pv *PluginValidator) validateNumber(field plugin.Field, value string) *FieldError {
	if value == "" {
		return nil
	}

	var num float64
	if _, err := fmt.Sscanf(value, "%f", &num); err != nil {
		return &FieldError{
			Field:   field.Key,
			Message: "Must be a valid number",
			Code:    "invalid_number",
		}
	}

	if field.Min != nil && num < *field.Min {
		return &FieldError{
			Field:   field.Key,
			Message: fmt.Sprintf("Must be at least %g", *field.Min),
			Code:    "below_minimum",
		}
	}

	if field.Max != nil && num > *field.Max {
		return &FieldError{
			Field:   field.Key,
			Message: fmt.Sprintf("Must be at most %g", *field.Max),
			Code:    "above_maximum",
		}
	}

	return nil
}

func (pv *PluginValidator) validateSelect(field plugin.Field, value string) *FieldError {
	if value == "" {
		return nil
	}

	for _, option := range field.Enum {
		if value == option {
			return nil
		}
	}

	return &FieldError{
		Field:   field.Key,
		Message: fmt.Sprintf("Must be one of: %s", strings.Join(field.Enum, ", ")),
		Code:    "invalid_option",
	}
}

func (pv *PluginValidator) isValidFieldType(fieldType string) bool {
	validTypes := []string{"color", "text", "number", "select"}
	for _, validType := range validTypes {
		if fieldType == validType {
			return true
		}
	}
	return false
}