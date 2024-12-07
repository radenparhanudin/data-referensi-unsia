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

func AllCities(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	countries, err := models.AllCities(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.ResponseMessage("get", false))
	}

	results := map[string]interface{}{
		"data": countries,
		"metadata": map[string]interface{}{
			"page":      page,
			"per_page":  pageSize,
			"sub_total": len(countries),
			"total":     models.CountCities(),
		},
	}

	return handlers.SendSuccess(c, fiber.StatusOK, results, helpers.ResponseMessage("get", true))
}

func ExportCities(c *fiber.Ctx) error {
	outputFile := "Cities.xlsx"

	if err := models.ExportCities(c, outputFile); err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.ResponseMessage("export", false))
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", outputFile))
	return c.SendFile(outputFile, false)
}

func SearchCities(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	countries, err := models.AllCities(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.ResponseMessage("get", false))
	}

	return handlers.SendSuccess(c, fiber.StatusOK, countries, helpers.ResponseMessage("get", true))
}

func CityById(c *fiber.Ctx) error {
	id := c.Params("id")
	province, err := models.CityById(id)
	if err != nil {
		return handlers.SendSuccess(c, fiber.StatusBadRequest, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusOK, province, helpers.ResponseMessage("get", true))
}

func CreateCity(c *fiber.Ctx) error {
	var req requests.CityRequest

	if err := c.BodyParser(&req); err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	err := models.CreateCity(req.ProvinceId, req.Name, req.Code)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key row") {
			return handlers.SendFailed(c, fiber.StatusBadRequest, nil, helpers.ResponseMessage("insert", false))
		}
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.ResponseMessage("insert", true))
}

func ImportCities(c *fiber.Ctx) error {
	file, err := c.FormFile("file_import")
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	filePath := fmt.Sprintf("./uploads/temp/%s", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, helpers.ResponseMessage("save", false))
	}

	if err := models.ImportCities(filePath); err != nil {
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

func UpdateCity(c *fiber.Ctx) error {
	id := c.Params("id")

	var req requests.CityRequest

	if err := c.BodyParser(&req); err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	err := models.UpdateCity(id, req.ProvinceId, req.Name, req.Code)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key row") {
			return handlers.SendFailed(c, fiber.StatusBadRequest, nil, helpers.ResponseMessage("exist"))
		}
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.ResponseMessage("update", true))
}

func DeleteCity(c *fiber.Ctx) error {
	id := c.Params("id")

	err := models.DeleteCity(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.ResponseMessage("delete", true))
}

func TrashAllCities(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)

	countries, err := models.TrashAllCities(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.ResponseMessage("get", false))
	}

	results := map[string]interface{}{
		"data": countries,
		"metadata": map[string]interface{}{
			"page":      page,
			"per_page":  pageSize,
			"sub_total": len(countries),
			"total":     models.TrashCountCities(),
		},
	}

	return handlers.SendSuccess(c, fiber.StatusOK, results, helpers.ResponseMessage("get", true))
}

// func RestoreCity(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	return c.JSON(fiber.Map{"message": "RestoreCity", "id": id})
// }
