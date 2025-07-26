package commons

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

func ValidateStruct(obj interface{}) error {
	err := validate.Struct(obj)
	if err == nil {
		return nil
	}

	var sb strings.Builder
	sb.WriteString("error(s): ")

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, fieldErr := range validationErrors {
			fieldName := strings.ToLower(fieldErr.Field())
			switch fieldErr.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf(`"%s" is required field`, fieldName))
			case "email":
				messages = append(messages, fmt.Sprintf(`"%s" must be a valid email`, fieldName))
			case "url":
				messages = append(messages, fmt.Sprintf(`"%s" must be a valid URL`, fieldName))
			case "oneof":
				messages = append(messages, fmt.Sprintf(`"%s" must be one of [%s]`, fieldName, fieldErr.Param()))
			case "min":
				messages = append(messages, fmt.Sprintf(`"%s" must have at least %s item(s)`, fieldName, fieldErr.Param()))
			default:
				messages = append(messages, fmt.Sprintf(`"%s" is invalid`, fieldName))
			}
		}
		sb.WriteString(strings.Join(messages, ", "))
		return fmt.Errorf(sb.String())
	}

	// If it's not a validation error, just return it
	return err
}

func GetQueryInt(c echo.Context, name string, defaultVal int) int {
	val := c.QueryParam(name)
	if val == "" {
		return defaultVal
	}
	if i, err := strconv.Atoi(val); err == nil {
		return i
	}
	return defaultVal
}
