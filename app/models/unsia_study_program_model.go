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

type MstUnsiaStudyProgram struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

type MstUnsiaStudyProgramExport struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

type MstUnsiaStudyProgramSearch struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

type MstUnsiaStudyProgramRelation struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

/* Action */
func GetUnsiaStudyPrograms(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstUnsiaStudyProgram, error) {
	return QueryGetUnsiaStudyPrograms("sp_mst_unsia_study_programs_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportUnsiaStudyPrograms(c *fiber.Ctx, fileSaveAs string) error {
	unsiaStudyPrograms, err := QueryExportUnsiaStudyPrograms()
	if err != nil {
		return fmt.Errorf("failed to get unsia study programs: %v", err)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.NewSheet(sheetName)

	columns := []string{"A", "B", "C"}
	headers := []string{"ID", "Code", "Name"}

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", col)
		file.SetCellValue(sheetName, cell, headers[i])
	}

	for i, unsiaStudyProgram := range unsiaStudyPrograms {
		row := i + 2

		values := []interface{}{
			unsiaStudyProgram.ID,
			unsiaStudyProgram.Code,
			unsiaStudyProgram.Name,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(unsiaStudyPrograms))
	}

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchUnsiaStudyPrograms(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstUnsiaStudyProgramSearch, error) {
	return QuerySearchUnsiaStudyPrograms("sp_mst_unsia_study_programs_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetUnsiaStudyProgram(id string) (MstUnsiaStudyProgram, error) {
	return QueryGetUnsiaStudyProgram(id)
}

func CreateUnsiaStudyProgram(id string, code string, name string) error {
	return QueryInsertUnsiaStudyProgram(id, code, name)
}

func ImportUnsiaStudyPrograms(filePath string) error {
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
		var code string = ""
		var name string = ""

		if len(row) > 0 && row[0] != "" {
			id = row[0]
		}
		if len(row) > 1 && row[1] != "" {
			code = row[1]
		}
		if len(row) > 2 && row[2] != "" {
			name = row[2]
		}

		if id != "" {
			exist, err := helpers.CheckModelIDExist(id, &MstUnsiaStudyProgram{})
			if err != nil {
				return err
			}
			if exist {
				if err := QueryUpdateUnsiaStudyProgram(id, code, name); err != nil {
					return err
				}
			} else {
				if err := QueryInsertUnsiaStudyProgram(id, code, name); err != nil {
					return err
				}
			}
		} else {
			id, err := helpers.EnsureUUID(&MstUnsiaStudyProgram{})
			if err != nil {
				return err
			}
			if err := QueryInsertUnsiaStudyProgram(id, code, name); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateUnsiaStudyProgram(id string, code string, name string) error {
	return QueryUpdateUnsiaStudyProgram(id, code, name)
}

func DeleteUnsiaStudyProgram(id string) error {
	if err := QueryDeleteUnsiaStudyProgram(id); err != nil {
		return err
	}

	return nil
}

func GetTrashUnsiaStudyPrograms(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstUnsiaStudyProgram, error) {
	return QueryGetUnsiaStudyPrograms("sp_mst_unsia_study_programs_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreUnsiaStudyProgram(id string) error {
	if err := QueryRestoreUnsiaStudyProgram(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountUnsiaStudyPrograms() int64 {
	return helpers.CountModelSize(&MstUnsiaStudyProgram{}, true)
}

func CountTrashUnsiaStudyPrograms() int64 {
	return helpers.CountModelSize(&MstUnsiaStudyProgram{}, false)
}

/* Query */
func QueryGetUnsiaStudyPrograms(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstUnsiaStudyProgram, error) {
	db := config.DB
	var unsia_study_programs []MstUnsiaStudyProgram

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&unsia_study_programs).Error
	if err != nil {
		return nil, err
	}

	return unsia_study_programs, nil
}

func QueryExportUnsiaStudyPrograms() ([]MstUnsiaStudyProgramExport, error) {
	db := config.DB
	var unsia_study_programs []MstUnsiaStudyProgramExport

	query := `
        EXEC sp_mst_unsia_study_programs_get 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `

	err := db.Raw(query, "", "name", "asc", 1, CountUnsiaStudyPrograms()).Scan(&unsia_study_programs).Error
	if err != nil {
		return nil, err
	}

	return unsia_study_programs, nil
}

func QuerySearchUnsiaStudyPrograms(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstUnsiaStudyProgramSearch, error) {
	db := config.DB
	var unsia_study_programs []MstUnsiaStudyProgramSearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&unsia_study_programs).Error
	if err != nil {
		return nil, err
	}

	return unsia_study_programs, nil
}

func QueryGetUnsiaStudyProgram(id string) (MstUnsiaStudyProgram, error) {
	db := config.DB
	var unsiaStudyProgram MstUnsiaStudyProgram

	query := `
		EXEC sp_mst_unsia_study_programs_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&unsiaStudyProgram).Error
	if err != nil {
		return MstUnsiaStudyProgram{}, err
	}

	return unsiaStudyProgram, nil
}

func QueryGetUnsiaStudyProgramRelation(id string) (MstUnsiaStudyProgramRelation, error) {
	db := config.DB
	var unsiaStudyProgram MstUnsiaStudyProgramRelation

	query := `
		EXEC sp_mst_unsia_study_programs_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&unsiaStudyProgram)
	if result.Error != nil {
		return MstUnsiaStudyProgramRelation{}, result.Error
	}

	return unsiaStudyProgram, nil
}

func QueryInsertUnsiaStudyProgram(id string, code string, name string) error {
	db := config.DB
	now := time.Now()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
		EXEC sp_mst_unsia_study_programs_insert
		@id = ?,
		@code = ?,
		@name = ?,
		@created_at = ?,
		@created_by = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, code, name, created_at, created_by, updated_at, updated_by).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryUpdateUnsiaStudyProgram(id string, code string, name string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
		EXEC sp_mst_unsia_study_programs_update
		@id = ?,
		@code = ?,
		@name = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, code, name, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func QueryDeleteUnsiaStudyProgram(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_unsia_study_programs_delete
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

func QueryRestoreUnsiaStudyProgram(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_unsia_study_programs_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
