package requests

import (
	"data-referensi/helpers"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PaginationRequest struct {
	PerPage string `query:"per_page" validate:"omitempty,numeric"`
	Page    string `query:"page" validate:"omitempty,numeric"`
}

func ValidatePagination(c *fiber.Ctx) error {
	var request PaginationRequest

	if err := c.QueryParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	validate := validator.New()

	err := validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)

		for _, fieldError := range validationErrors {
			fieldName := helpers.CamelCaseToSnakeCase(fieldError.Field())
			tag := fieldError.Tag()
			errorMessages[fieldName] = helpers.GenerateValidationErrorMessage(fieldName, tag)
		}

		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Validation error",
			"errors":  errorMessages,
		})
	}

	perPage, err := strconv.Atoi(request.PerPage)
	if err != nil || perPage <= 0 {
		perPage = 10
	}

	page, err := strconv.Atoi(request.Page)
	if err != nil || page <= 0 {
		page = 1
	}

	c.Locals("per_page", perPage)
	c.Locals("page", page)

	return c.Next()
}
