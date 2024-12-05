package controllers

import (
	"data-referensi/app/models"
	"data-referensi/handlers"
	"data-referensi/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func AllCountries(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)

	countries, err := models.GetCountries(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get countries",
		})
	}

	results := map[string]interface{}{
		"data":      countries,
		"page":      page,
		"page_size": pageSize,
	}

	return handlers.SendSuccess(c, fiber.StatusOK, results, "Get data countries")
}

func ExportCountries(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "ExportCountries"})
}

func CountryById(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "CountryById", "id": id})
}

func CreateCountry(c *fiber.Ctx) error {
	var country models.CreateOrUpdateCountry
	if err := helpers.BodyParser(c, &country); err != nil {
		return err
	}

	validate := validator.New()
	err := validate.Struct(country)
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

	return c.JSON(fiber.Map{"message": "CreateCountry"})
}

func ImportCountries(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "ImportCountries"})
}

func UpdateCountry(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "UpdateCountry"})
}

func DeleteCountry(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "DeleteCountry", "id": id})
}

func TrashAllCountries(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "TrashAllCountries"})
}

func RestoreCountry(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "RestoreCountry", "id": id})
}

func DeletePermanentCountry(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"message": "DeletePermanentCountry", "id": id})
}
