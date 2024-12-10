package models

import (
	"data-referensi/config"
	"data-referensi/helpers"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type MstEducation struct {
	ID                 uuid.UUID                    `json:"id"`
	EducationalLevelId string                       `json:"educational_level_id"`
	EducationalLevel   *MstEducationalLevelRelation `json:"educational_level"`
	StudyProgramId     string                       `json:"study_program_id"`
	StudyProgram       *MstStudyProgramRelation     `json:"study_program"`
	Name               string                       `json:"name"`
}

type MstEducationExport struct {
	ID                 uuid.UUID `json:"id"`
	EducationalLevelId string    `json:"educational_level_id"`
	StudyProgramId     string    `json:"study_program_id"`
	Name               string    `json:"name"`
}

type MstEducationSearch struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
type MstEducationRelation struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

/* Action */
func GetEducations(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEducation, error) {
	return QueryGetEducations("sp_mst_educations_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportEducations(c *fiber.Ctx, fileSaveAs string) error {
	educations, err := QueryExportEducations()
	if err != nil {
		return fmt.Errorf("failed to get educations: %v", err)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.NewSheet(sheetName)

	columns := []string{"A", "B", "C", "D"}
	headers := []string{"ID", "Educational Level ID", "Study Program ID", "Name"}

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", col)
		file.SetCellValue(sheetName, cell, headers[i])
	}

	for i, education := range educations {
		row := i + 2

		values := []interface{}{
			education.ID,
			education.EducationalLevelId,
			education.StudyProgramId,
			education.Name,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(educations))
	}

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchEducations(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEducationSearch, error) {
	return QuerySearchEducations("sp_mst_educations_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetEducationByEducationalLevelId(ducation_level_id string) ([]MstEducationSearch, error) {
	return QueryGetEducationByEducationalLevelId(ducation_level_id)
}

func GetEducation(id string) (MstEducation, error) {
	return QueryGetEducation(id)
}

func CreateEducation(id string, educational_level_id string, study_program_id string, name string) error {
	return QueryInsertEducation(id, educational_level_id, study_program_id, name)
}

func ImportEducations(filePath string) error {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open excel file: %v", err)
	}

	sheetName := "Sheet1"
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get rows: %v", err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}

		var id string = ""
		var educational_level_id string = ""
		var study_program_id string = ""
		var name string = ""

		if len(row) > 0 && row[0] != "" {
			id = row[0]
		}
		if len(row) > 1 && row[1] != "" {
			educational_level_id = row[1]
		}
		if len(row) > 2 && row[2] != "" {
			study_program_id = row[2]
		}
		if len(row) > 3 && row[3] != "" {
			name = row[3]
		}

		if id != "" {
			exist, err := helpers.CheckModelIDExist(id, &MstEducation{})
			if err != nil {
				return err
			}
			if exist {
				if err := QueryUpdateEducation(id, educational_level_id, study_program_id, name); err != nil {
					return err
				}
			} else {
				if err := QueryInsertEducation(id, educational_level_id, study_program_id, name); err != nil {
					return err
				}
			}
		} else {
			id, err := helpers.EnsureUUID(&MstEducation{})
			if err != nil {
				return err
			}
			if err := QueryInsertEducation(id, educational_level_id, study_program_id, name); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateEducation(id string, educational_level_id string, study_program_id string, name string) error {
	return QueryUpdateEducation(id, educational_level_id, study_program_id, name)
}

func DeleteEducation(id string) error {
	if err := QueryDeleteEducation(id); err != nil {
		return err
	}

	return nil
}

func GetTrashEducations(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEducation, error) {
	return QueryGetEducations("sp_mst_educations_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreEducation(id string) error {
	if err := QueryRestoreEducation(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountEducations() int64 {
	return helpers.CountModelSize(&MstEducation{}, true)
}

func CountTrashEducations() int64 {
	return helpers.CountModelSize(&MstEducation{}, false)
}

/* Query */
func QueryGetEducations(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEducation, error) {
	db := config.DB
	var educations []MstEducation

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&educations).Error
	if err != nil {
		return nil, err
	}

	for i := range educations {
		educationalLevel, err := QueryGetEducationalLevelRelation(string(educations[i].EducationalLevelId))
		if err != nil {
			return []MstEducation{}, err
		}

		educations[i].EducationalLevel = &educationalLevel
	}

	return educations, nil
}

func QueryExportEducations() ([]MstEducationExport, error) {
	db := config.DB
	var educations []MstEducationExport

	query := `
        EXEC sp_mst_educations_get 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `

	err := db.Raw(query, "", "name", "asc", 1, CountEducations()).Scan(&educations).Error
	if err != nil {
		return nil, err
	}

	return educations, nil
}

func QuerySearchEducations(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEducationSearch, error) {
	db := config.DB
	var educations []MstEducationSearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&educations).Error
	if err != nil {
		return nil, err
	}

	return educations, nil
}

func QueryGetEducationByEducationalLevelId(ducation_level_id string) ([]MstEducationSearch, error) {
	db := config.DB
	var educations []MstEducationSearch

	query := `
		EXEC sp_mst_educations_get_by_education_level_id
		@ducation_level_id = ?
	`
	err := db.Raw(query, ducation_level_id).Scan(&educations).Error
	if err != nil {
		return []MstEducationSearch{}, err
	}

	return educations, nil
}

func QueryGetEducation(id string) (MstEducation, error) {
	db := config.DB
	var education MstEducation

	query := `
		EXEC sp_mst_educations_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&education).Error
	if err != nil {
		return MstEducation{}, err
	}

	educationalLevel, err := QueryGetEducationalLevelRelation(education.EducationalLevelId)
	if err != nil {
		return MstEducation{}, err
	}

	education.EducationalLevel = &educationalLevel

	studyProgram, err := QueryGetStudyProgramRelation(education.StudyProgramId)
	if err != nil {
		return MstEducation{}, err
	}

	education.StudyProgram = &studyProgram

	return education, nil
}

func QueryGetEducationRelation(id string) (MstEducationRelation, error) {
	db := config.DB
	var education MstEducationRelation

	query := `
		EXEC sp_mst_educations_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&education)
	if result.Error != nil {
		return MstEducationRelation{}, result.Error
	}

	return education, nil
}

func QueryInsertEducation(id string, educational_level_id string, study_program_id string, name string) error {
	db := config.DB
	now := time.Now()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
		EXEC sp_mst_educations_insert
		@id = ?,
		@educational_level_id = ?,
		@study_program_id = ?,
		@name = ?,
		@created_at = ?,
		@created_by = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, educational_level_id, study_program_id, name, created_at, created_by, updated_at, updated_by).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryUpdateEducation(id string, educational_level_id string, study_program_id string, name string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
		EXEC sp_mst_educations_update
		@id = ?,
		@educational_level_id = ?,
		@study_program_id = ?,
		@name = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, educational_level_id, study_program_id, name, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func QueryDeleteEducation(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_educations_delete
		@id = ?,
		@deleted_at = ?,
		@deleted_by = ?
	`

	err := db.Exec(query, id, deleted_at, deleted_by).Error
	if err != nil {
		return err
	}

	return nil
}

func QueryRestoreEducation(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_educations_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
