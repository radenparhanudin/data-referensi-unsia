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

func GetCountries(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	countries, err := models.GetCountries(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.GenerateRM("get", false))
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

	return handlers.SendSuccess(c, fiber.StatusOK, results, helpers.GenerateRM("get", true))
}

func ExportCountries(c *fiber.Ctx) error {
	fileName := "Countries.xlsx"
	fileSaveAs := fmt.Sprintf("tmp/exports/%s", fileName)

	if err := models.ExportCountries(c, fileSaveAs); err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.GenerateRM("export", false))
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	return c.SendFile(fileSaveAs, false)
}

func SearchCountries(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	countries, err := models.SearchCountries(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.GenerateRM("get", false))
	}

	return handlers.SendSuccess(c, fiber.StatusOK, countries, helpers.GenerateRM("get", true))
}

func GetCountry(c *fiber.Ctx) error {
	id := c.Params("id")
	country, err := models.GetCountry(id)
	if err != nil {
		return handlers.SendSuccess(c, fiber.StatusBadRequest, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusOK, country, helpers.GenerateRM("get", true))
}

func CreateCountry(c *fiber.Ctx) error {
	var req requests.CountryRequest

	if err := c.BodyParser(&req); err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	/* Check Existing ID */
	id, err := helpers.EnsureUUID(&models.MstCountry{})
	if err != nil {
		return err
	}

	err = models.CreateCountry(id, req.Name, req.PhoneCode, req.IconFlagPath)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key row") {
			return handlers.SendFailed(c, fiber.StatusBadRequest, nil, helpers.GenerateRM("exist"))
		}
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	country, err := models.GetCountry(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, country, helpers.GenerateRM("insert", true))
}

func ImportCountries(c *fiber.Ctx) error {
	file, err := c.FormFile("file_import")
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	filePath := fmt.Sprintf("./tmp/uploads/%s", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, helpers.GenerateRM("save", false))
	}

	if err := models.ImportCountries(filePath); err != nil {
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

func UpdateCountry(c *fiber.Ctx) error {
	id := c.Params("id")

	var req requests.CountryRequest

	if err := c.BodyParser(&req); err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	err := models.UpdateCountry(id, req.Name, req.PhoneCode, req.IconFlagPath)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key row") {
			return handlers.SendFailed(c, fiber.StatusBadRequest, nil, helpers.GenerateRM("exist"))
		}
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	country, err := models.GetCountry(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, country, helpers.GenerateRM("update", true))
}

func DeleteCountry(c *fiber.Ctx) error {
	id := c.Params("id")

	err := models.DeleteCountry(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.GenerateRM("delete", true))
}

func GetTrashCountries(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	countries, err := models.GetTrashCountries(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.GenerateRM("get", false))
	}

	results := map[string]interface{}{
		"data": countries,
		"metadata": map[string]interface{}{
			"page":      page,
			"per_page":  pageSize,
			"sub_total": len(countries),
			"total":     models.CountTrashCountries(),
		},
	}

	return handlers.SendSuccess(c, fiber.StatusOK, results, helpers.GenerateRM("get", true))
}

func RestoreCountry(c *fiber.Ctx) error {
	id := c.Params("id")

	err := models.RestoreCountry(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.GenerateRM("restore", true))
}
