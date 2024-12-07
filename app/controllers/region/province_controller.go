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

func AllProvinces(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	countries, err := models.AllProvinces(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.ResponseMessage("get", false))
	}

	results := map[string]interface{}{
		"data": countries,
		"metadata": map[string]interface{}{
			"page":      page,
			"per_page":  pageSize,
			"sub_total": len(countries),
			"total":     models.CountProvinces(),
		},
	}

	return handlers.SendSuccess(c, fiber.StatusOK, results, helpers.ResponseMessage("get", true))
}

func ExportProvinces(c *fiber.Ctx) error {
	outputFile := "Provinces.xlsx"

	if err := models.ExportProvinces(c, outputFile); err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.ResponseMessage("export", false))
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", outputFile))
	return c.SendFile(outputFile, false)
}

func SearchProvinces(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	countries, err := models.AllProvinces(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.ResponseMessage("get", false))
	}

	return handlers.SendSuccess(c, fiber.StatusOK, countries, helpers.ResponseMessage("get", true))
}

func ProvinceById(c *fiber.Ctx) error {
	id := c.Params("id")
	province, err := models.ProvinceById(id)
	if err != nil {
		return handlers.SendSuccess(c, fiber.StatusBadRequest, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusOK, province, "Successfully retrieved province data")
}

func CreateProvince(c *fiber.Ctx) error {
	var req requests.ProvinceRequest

	if err := c.BodyParser(&req); err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	err := models.CreateProvince(req.CountryId, req.Name, req.Code, req.RegionCode)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key row") {
			return handlers.SendFailed(c, fiber.StatusBadRequest, nil, helpers.ResponseMessage("insert", false))
		}
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.ResponseMessage("insert", true))
}

func ImportProvinces(c *fiber.Ctx) error {
	file, err := c.FormFile("file_import")
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	filePath := fmt.Sprintf("./uploads/temp/%s", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, helpers.ResponseMessage("save", false))
	}

	if err := models.ImportProvinces(filePath); err != nil {
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

func UpdateProvince(c *fiber.Ctx) error {
	id := c.Params("id")

	var req requests.ProvinceRequest

	if err := c.BodyParser(&req); err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	err := models.UpdateProvince(id, req.CountryId, req.Name, req.Code, req.RegionCode)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key row") {
			return handlers.SendFailed(c, fiber.StatusBadRequest, nil, helpers.ResponseMessage("exist"))
		}
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.ResponseMessage("update", true))
}

func DeleteProvince(c *fiber.Ctx) error {
	id := c.Params("id")

	err := models.DeleteProvince(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.ResponseMessage("delete", true))
}

func TrashAllProvinces(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)

	countries, err := models.TrashAllProvinces(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.ResponseMessage("get", false))
	}

	results := map[string]interface{}{
		"data": countries,
		"metadata": map[string]interface{}{
			"page":      page,
			"per_page":  pageSize,
			"sub_total": len(countries),
			"total":     models.CountProvinces(),
		},
	}

	return handlers.SendSuccess(c, fiber.StatusOK, results, helpers.ResponseMessage("get", true))
}

// func RestoreProvince(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	return c.JSON(fiber.Map{"message": "RestoreProvince", "id": id})
// }
