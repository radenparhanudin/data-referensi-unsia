package requests

import (
	"data-referensi/handlers"
	"data-referensi/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type EducationRequest struct {
	EducationalLevelId string `json:"educational_level_id" validate:"required,max=36"`
	StudyProgramId     string `json:"study_program_id" validate:"max=36"`
	Name               string `json:"name" validate:"required,max=255"`
}

func ValidateEducation(c *fiber.Ctx) error {
	validate := validator.New()
	var req EducationRequest
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
