package requests

import (
	"data-referensi/handlers"
	"data-referensi/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ReligionRequest struct {
	Code string `json:"code" validate:"required,max=2"`
	Name string `json:"name" validate:"required,max=255"`
}

func ValidateReligion(c *fiber.Ctx) error {
	validate := validator.New()
	var req ReligionRequest
	if err := c.BodyParser(&req); err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)

		for _, fieldError := range validationErrors {
			fieldName := helpers.ConvertCCToSC(fieldError.Field())
			errorMessages[fieldName] = helpers.GenerateVEM(fieldName, fieldError.Tag())
		}

		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   true,
			"message": errorMessages,
		})
	}
	return c.Next()
}