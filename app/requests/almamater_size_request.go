package requests

import (
	"data-referensi/handlers"
	"data-referensi/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AlmamaterSizeRequest struct {
	Code       string `json:"code" validate:"required,max=50"`
	Size       string `json:"size" validate:"required,max=255"`
	ChestSize  string `json:"chest_size" validate:"required,max=255"`
	ArmLength  string `json:"arm_length" validate:"required,max=255"`
	BodyLength string `json:"body_length" validate:"required,max=255"`
}

func ValidateAlmamaterSize(c *fiber.Ctx) error {
	validate := validator.New()
	var req AlmamaterSizeRequest
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
