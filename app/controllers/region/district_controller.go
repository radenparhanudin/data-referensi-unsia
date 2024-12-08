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

func GetDistricts(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	countries, err := models.GetDistricts(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.GenerateRM("get", false))
	}

	results := map[string]interface{}{
		"data": countries,
		"metadata": map[string]interface{}{
			"page":      page,
			"per_page":  pageSize,
			"sub_total": len(countries),
			"total":     models.CountDistricts(),
		},
	}

	return handlers.SendSuccess(c, fiber.StatusOK, results, helpers.GenerateRM("get", true))
}

func ExportDistricts(c *fiber.Ctx) error {
	fileName := "Districts.xlsx"
	fileSaveAs := fmt.Sprintf("tmp/exports/%s", fileName)

	if err := models.ExportDistricts(c, fileSaveAs); err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.GenerateRM("export", false))
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	return c.SendFile(fileSaveAs, false)
}

func SearchDistricts(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	countries, err := models.SearchDistricts(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.GenerateRM("get", false))
	}

	return handlers.SendSuccess(c, fiber.StatusOK, countries, helpers.GenerateRM("get", true))
}

func GetDistrict(c *fiber.Ctx) error {
	id := c.Params("id")
	country, err := models.GetDistrict(id)
	if err != nil {
		return handlers.SendSuccess(c, fiber.StatusBadRequest, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusOK, country, helpers.GenerateRM("get", true))
}

func CreateDistrict(c *fiber.Ctx) error {
	var req requests.DistrictRequest

	if err := c.BodyParser(&req); err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	/* Check Existing ID */
	id, err := helpers.EnsureUUID(&models.MstDistrict{})
	if err != nil {
		return err
	}

	err = models.CreateDistrict(id, req.CityId, req.Name, req.Code)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key row") {
			return handlers.SendFailed(c, fiber.StatusBadRequest, nil, helpers.GenerateRM("exist"))
		}
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	country, err := models.GetDistrict(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, country, helpers.GenerateRM("insert", true))
}

func ImportDistricts(c *fiber.Ctx) error {
	file, err := c.FormFile("file_import")
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	filePath := fmt.Sprintf("./tmp/uploads/%s", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, helpers.GenerateRM("save", false))
	}

	if err := models.ImportDistricts(filePath); err != nil {
		if strings.Contains(err.Error(), "duplicate key row") {
			return handlers.SendFailed(c, fiber.StatusBadRequest, nil, helpers.GenerateRM("exist"))
		}
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	if err := os.Remove(filePath); err != nil {
		log.Println("Error removing uploaded file:", err)
	}

	return handlers.SendSuccess(c, fiber.StatusOK, nil, helpers.GenerateRM("import", true))
}

func UpdateDistrict(c *fiber.Ctx) error {
	id := c.Params("id")

	var req requests.DistrictRequest

	if err := c.BodyParser(&req); err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	err := models.UpdateDistrict(id, req.CityId, req.Name, req.Code)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key row") {
			return handlers.SendFailed(c, fiber.StatusBadRequest, nil, helpers.GenerateRM("exist"))
		}
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	country, err := models.GetDistrict(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, country, helpers.GenerateRM("update", true))
}

func DeleteDistrict(c *fiber.Ctx) error {
	id := c.Params("id")

	err := models.DeleteDistrict(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.GenerateRM("delete", true))
}

func GetTrashDistricts(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	countries, err := models.GetTrashDistricts(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.GenerateRM("get", false))
	}

	results := map[string]interface{}{
		"data": countries,
		"metadata": map[string]interface{}{
			"page":      page,
			"per_page":  pageSize,
			"sub_total": len(countries),
			"total":     models.CountTrashDistricts(),
		},
	}

	return handlers.SendSuccess(c, fiber.StatusOK, results, helpers.GenerateRM("get", true))
}

func RestoreDistrict(c *fiber.Ctx) error {
	id := c.Params("id")

	err := models.RestoreDistrict(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.GenerateRM("restore", true))
}
