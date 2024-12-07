package controllers

import (
	"data-referensi/app/models"
	"data-referensi/app/requests"
	"data-referensi/handlers"
	"data-referensi/helpers"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AllCountries(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	countries, err := models.AllCountries(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.ResponseMessage("get", false))
	}

	results := map[string]interface{}{
		"data": countries,
		"metadata": map[string]interface{}{
			"page":      page,
			"per_page":  pageSize,
			"sub_total": len(countries),
			"total":     models.CountCountries(),
		},
	}

	return handlers.SendSuccess(c, fiber.StatusOK, results, helpers.ResponseMessage("get", true))
}

func ExportCountries(c *fiber.Ctx) error {
	outputFile := "Countries.xlsx"

	if err := models.ExportCountries(c, outputFile); err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.ResponseMessage("export", false))
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", outputFile))
	return c.SendFile(outputFile, false)
}

func SearchCountries(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	countries, err := models.AllCountries(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.ResponseMessage("get", false))
	}

	return handlers.SendSuccess(c, fiber.StatusOK, countries, helpers.ResponseMessage("get", true))
}

func CountryById(c *fiber.Ctx) error {
	id := c.Params("id")
	country, err := models.CountryById(id)
	if err != nil {
		return handlers.SendSuccess(c, fiber.StatusBadRequest, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusOK, country, "Successfully retrieved country data")
}

func CreateCountry(c *fiber.Ctx) error {
	var req requests.CountryRequest

	if err := c.BodyParser(&req); err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	err := models.CreateCountry(req.Name, req.PhoneCode, req.IconFlagPath)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key row") {
			return handlers.SendFailed(c, fiber.StatusBadRequest, nil, helpers.ResponseMessage("insert", false))
		}
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.ResponseMessage("insert", true))
}

func ImportCountries(c *fiber.Ctx) error {
	file, err := c.FormFile("file_import")
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	filePath := fmt.Sprintf("./uploads/temp/%s", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, helpers.ResponseMessage("save", false))
	}

	if err := models.ImportCountries(filePath); err != nil {
		if strings.Contains(err.Error(), "duplicate key row") {
			return handlers.SendFailed(c, fiber.StatusBadRequest, nil, helpers.ResponseMessage("exist"))
		}
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	if err := os.Remove(filePath); err != nil {
		log.Println("Error removing uploaded file:", err)
	}

	return handlers.SendSuccess(c, fiber.StatusOK, nil, helpers.ResponseMessage("import", true))
}

func UpdateCountry(c *fiber.Ctx) error {
	id := c.Params("id")

	var req requests.CountryRequest

	if err := c.BodyParser(&req); err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	err := models.UpdateCountry(id, req.Name, req.PhoneCode, req.IconFlagPath)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key row") {
			return handlers.SendFailed(c, fiber.StatusBadRequest, nil, helpers.ResponseMessage("exist"))
		}
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.ResponseMessage("update", true))
}

func DeleteCountry(c *fiber.Ctx) error {
	id := c.Params("id")

	err := models.DeleteCountry(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.ResponseMessage("delete", true))
}

func TrashAllCountries(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)

	countries, err := models.TrashAllCountries(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.ResponseMessage("get", false))
	}

	results := map[string]interface{}{
		"data": countries,
		"metadata": map[string]interface{}{
			"page":      page,
			"per_page":  pageSize,
			"sub_total": len(countries),
			"total":     models.CountCountries(),
		},
	}

	return handlers.SendSuccess(c, fiber.StatusOK, results, helpers.ResponseMessage("get", true))
}

// func RestoreCountry(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	return c.JSON(fiber.Map{"message": "RestoreCountry", "id": id})
// }
