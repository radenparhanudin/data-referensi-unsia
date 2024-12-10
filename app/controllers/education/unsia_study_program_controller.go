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

func GetUnsiaStudyPrograms(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	studyPrograms, err := models.GetUnsiaStudyPrograms(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.GenerateRM("get", false))
	}

	results := map[string]interface{}{
		"data": studyPrograms,
		"metadata": map[string]interface{}{
			"page":      page,
			"per_page":  pageSize,
			"sub_total": len(studyPrograms),
			"total":     models.CountUnsiaStudyPrograms(),
		},
	}

	return handlers.SendSuccess(c, fiber.StatusOK, results, helpers.GenerateRM("get", true))
}

func ExportUnsiaStudyPrograms(c *fiber.Ctx) error {
	fileName := "UnsiaStudyPrograms.xlsx"
	fileSaveAs := fmt.Sprintf("tmp/exports/%s", fileName)

	if err := models.ExportUnsiaStudyPrograms(c, fileSaveAs); err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.GenerateRM("export", false))
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	return c.SendFile(fileSaveAs, false)
}

func SearchUnsiaStudyPrograms(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	studyPrograms, err := models.SearchUnsiaStudyPrograms(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.GenerateRM("get", false))
	}

	return handlers.SendSuccess(c, fiber.StatusOK, studyPrograms, helpers.GenerateRM("get", true))
}

func GetUnsiaStudyProgram(c *fiber.Ctx) error {
	id := c.Params("id")
	studyProgram, err := models.GetUnsiaStudyProgram(id)
	if err != nil {
		return handlers.SendSuccess(c, fiber.StatusBadRequest, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusOK, studyProgram, helpers.GenerateRM("get", true))
}

func CreateUnsiaStudyProgram(c *fiber.Ctx) error {
	var req requests.UnsiaStudyProgramRequest

	if err := c.BodyParser(&req); err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	/* Check Existing ID */
	id, err := helpers.EnsureUUID(&models.MstUnsiaStudyProgram{})
	if err != nil {
		return err
	}

	err = models.CreateUnsiaStudyProgram(id, req.Code, req.Name)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key row") {
			return handlers.SendFailed(c, fiber.StatusBadRequest, nil, helpers.GenerateRM("exist"))
		}
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	studyProgram, err := models.GetUnsiaStudyProgram(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, studyProgram, helpers.GenerateRM("insert", true))
}

func ImportUnsiaStudyPrograms(c *fiber.Ctx) error {
	file, err := c.FormFile("file_import")
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	filePath := fmt.Sprintf("./tmp/uploads/%s", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, helpers.GenerateRM("save", false))
	}

	if err := models.ImportUnsiaStudyPrograms(filePath); err != nil {
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

func UpdateUnsiaStudyProgram(c *fiber.Ctx) error {
	id := c.Params("id")

	var req requests.UnsiaStudyProgramRequest

	if err := c.BodyParser(&req); err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	err := models.UpdateUnsiaStudyProgram(id, req.Code, req.Name)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key row") {
			return handlers.SendFailed(c, fiber.StatusBadRequest, nil, helpers.GenerateRM("exist"))
		}
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	studyProgram, err := models.GetUnsiaStudyProgram(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusBadRequest, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, studyProgram, helpers.GenerateRM("update", true))
}

func DeleteUnsiaStudyProgram(c *fiber.Ctx) error {
	id := c.Params("id")

	err := models.DeleteUnsiaStudyProgram(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.GenerateRM("delete", true))
}

func GetTrashUnsiaStudyPrograms(c *fiber.Ctx) error {
	filter := c.Query("filter", "")
	sortBy := c.Query("sort_by", "name")
	sortDirection := c.Query("sort_direction", "asc")
	page := c.QueryInt("page", 1)
	pageSize := int64(c.QueryInt("page_size", 10))

	studyPrograms, err := models.GetTrashUnsiaStudyPrograms(filter, sortBy, sortDirection, page, pageSize)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusOK, nil, helpers.GenerateRM("get", false))
	}

	results := map[string]interface{}{
		"data": studyPrograms,
		"metadata": map[string]interface{}{
			"page":      page,
			"per_page":  pageSize,
			"sub_total": len(studyPrograms),
			"total":     models.CountTrashUnsiaStudyPrograms(),
		},
	}

	return handlers.SendSuccess(c, fiber.StatusOK, results, helpers.GenerateRM("get", true))
}

func RestoreUnsiaStudyProgram(c *fiber.Ctx) error {
	id := c.Params("id")

	err := models.RestoreUnsiaStudyProgram(id)
	if err != nil {
		return handlers.SendFailed(c, fiber.StatusInternalServerError, nil, err.Error())
	}

	return handlers.SendSuccess(c, fiber.StatusCreated, nil, helpers.GenerateRM("restore", true))
}
