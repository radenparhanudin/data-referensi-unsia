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

type MstStudyProgram struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type MstStudyProgramExport struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type MstStudyProgramSearch struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type MstStudyProgramRelation struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

/* Action */
func GetStudyPrograms(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstStudyProgram, error) {
	return QueryGetStudyPrograms("sp_mst_study_programs_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportStudyPrograms(c *fiber.Ctx, fileSaveAs string) error {
	studyPrograms, err := QueryExportStudyPrograms()
	if err != nil {
		return fmt.Errorf("failed to get study programs: %v", err)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.NewSheet(sheetName)

	columns := []string{"A", "B"}
	headers := []string{"ID", "Name"}

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", col)
		file.SetCellValue(sheetName, cell, headers[i])
	}

	for i, studyProgram := range studyPrograms {
		row := i + 2

		values := []interface{}{
			studyProgram.ID,
			studyProgram.Name,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(studyPrograms))
	}

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchStudyPrograms(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstStudyProgramSearch, error) {
	return QuerySearchStudyPrograms("sp_mst_study_programs_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetStudyProgram(id string) (MstStudyProgram, error) {
	return QueryGetStudyProgram(id)
}

func CreateStudyProgram(id string, name string) error {
	return QueryInsertStudyProgram(id, name)
}

func ImportStudyPrograms(filePath string) error {
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
		var name string = ""

		if len(row) > 0 && row[0] != "" {
			id = row[0]
		}
		if len(row) > 1 && row[1] != "" {
			name = row[1]
		}

		if id != "" {
			exist, err := helpers.CheckModelIDExist(id, &MstStudyProgram{})
			if err != nil {
				return err
			}
			if exist {
				if err := QueryUpdateStudyProgram(id, name); err != nil {
					return err
				}
			} else {
				if err := QueryInsertStudyProgram(id, name); err != nil {
					return err
				}
			}
		} else {
			id, err := helpers.EnsureUUID(&MstStudyProgram{})
			if err != nil {
				return err
			}
			if err := QueryInsertStudyProgram(id, name); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateStudyProgram(id string, name string) error {
	return QueryUpdateStudyProgram(id, name)
}

func DeleteStudyProgram(id string) error {
	if err := QueryDeleteStudyProgram(id); err != nil {
		return err
	}

	return nil
}

func GetTrashStudyPrograms(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstStudyProgram, error) {
	return QueryGetStudyPrograms("sp_mst_study_programs_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreStudyProgram(id string) error {
	if err := QueryRestoreStudyProgram(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountStudyPrograms() int64 {
	return helpers.CountModelSize(&MstStudyProgram{}, true)
}

func CountTrashStudyPrograms() int64 {
	return helpers.CountModelSize(&MstStudyProgram{}, false)
}

/* Query */
func QueryGetStudyPrograms(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstStudyProgram, error) {
	db := config.DB
	var studyPrograms []MstStudyProgram

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&studyPrograms).Error
	if err != nil {
		return nil, err
	}

	return studyPrograms, nil
}

func QueryExportStudyPrograms() ([]MstStudyProgramExport, error) {
	db := config.DB
	var studyPrograms []MstStudyProgramExport

	query := `
        EXEC sp_mst_study_programs_get 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `

	err := db.Raw(query, "", "name", "asc", 1, CountStudyPrograms()).Scan(&studyPrograms).Error
	if err != nil {
		return nil, err
	}

	return studyPrograms, nil
}

func QuerySearchStudyPrograms(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstStudyProgramSearch, error) {
	db := config.DB
	var studyPrograms []MstStudyProgramSearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&studyPrograms).Error
	if err != nil {
		return nil, err
	}

	return studyPrograms, nil
}

func QueryGetStudyProgram(id string) (MstStudyProgram, error) {
	db := config.DB
	var studyProgram MstStudyProgram

	query := `
		EXEC sp_mst_study_programs_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&studyProgram).Error
	if err != nil {
		return MstStudyProgram{}, err
	}

	return studyProgram, nil
}

func QueryGetStudyProgramRelation(id string) (MstStudyProgramRelation, error) {
	db := config.DB
	var studyProgram MstStudyProgramRelation

	query := `
		EXEC sp_mst_study_programs_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&studyProgram)
	if result.Error != nil {
		return MstStudyProgramRelation{}, result.Error
	}

	return studyProgram, nil
}

func QueryInsertStudyProgram(id string, name string) error {
	db := config.DB
	now := time.Now()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
		EXEC sp_mst_study_programs_insert
		@id = ?,
		@name = ?,
		@created_at = ?,
		@created_by = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, name, created_at, created_by, updated_at, updated_by).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryUpdateStudyProgram(id string, name string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
		EXEC sp_mst_study_programs_update
		@id = ?,
		@name = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, name, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func QueryDeleteStudyProgram(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_study_programs_delete
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

func QueryRestoreStudyProgram(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_study_programs_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
