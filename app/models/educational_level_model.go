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

type MstEducationalLevel struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   int64     `json:"created_at"`
	UpdatedAt   int64     `json:"updated_at"`
}

type MstEducationalLevelExport struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type MstEducationalLevelSearch struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

type MstEducationalLevelRelation struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

/* Action */
func GetEducationalLevels(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEducationalLevel, error) {
	return QueryGetEducationalLevels("sp_mst_educational_levels_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportEducationalLevels(c *fiber.Ctx, fileSaveAs string) error {
	educational_levels, err := QueryExportEducationalLevels()
	if err != nil {
		return fmt.Errorf("failed to get educational_levels: %v", err)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.NewSheet(sheetName)

	columns := []string{"A", "B", "C", "D"}
	headers := []string{"ID", "Code", "Name", "Description"}

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", col)
		file.SetCellValue(sheetName, cell, headers[i])
	}

	for i, job := range educational_levels {
		row := i + 2

		values := []interface{}{
			job.ID,
			job.Code,
			job.Name,
			job.Description,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(educational_levels))
	}

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchEducationalLevels(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEducationalLevelSearch, error) {
	return QuerySearchEducationalLevels("sp_mst_educational_levels_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetEducationalLevel(id string) (MstEducationalLevel, error) {
	return QueryGetEducationalLevel(id)
}

func CreateEducationalLevel(id string, code string, name string, description string) error {
	return QueryInsertEducationalLevel(id, code, name, description)
}

func ImportEducationalLevels(filePath string) error {
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
		var description string = ""

		if len(row) > 0 && row[0] != "" {
			id = row[0]
		}
		if len(row) > 1 && row[1] != "" {
			code = row[1]
		}
		if len(row) > 2 && row[2] != "" {
			name = row[2]
		}
		if len(row) > 3 && row[3] != "" {
			description = row[3]
		}

		if id != "" {
			exist, err := helpers.CheckModelIDExist(id, &MstEducationalLevel{})
			if err != nil {
				return err
			}
			if exist {
				if err := QueryUpdateEducationalLevel(id, code, name, description); err != nil {
					return err
				}
			} else {
				if err := QueryInsertEducationalLevel(id, code, name, description); err != nil {
					return err
				}
			}
		} else {
			id, err := helpers.EnsureUUID(&MstEducationalLevel{})
			if err != nil {
				return err
			}
			if err := QueryInsertEducationalLevel(id, code, name, description); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateEducationalLevel(id string, code string, name string, description string) error {
	return QueryUpdateEducationalLevel(id, code, name, description)
}

func DeleteEducationalLevel(id string) error {
	if err := QueryDeleteEducationalLevel(id); err != nil {
		return err
	}

	return nil
}

func GetTrashEducationalLevels(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEducationalLevel, error) {
	return QueryGetEducationalLevels("sp_mst_educational_levels_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreEducationalLevel(id string) error {
	if err := QueryRestoreEducationalLevel(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountEducationalLevels() int64 {
	return helpers.CountModelSize(&MstEducationalLevel{}, true)
}

func CountTrashEducationalLevels() int64 {
	return helpers.CountModelSize(&MstEducationalLevel{}, false)
}

/* Query */
func QueryGetEducationalLevels(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEducationalLevel, error) {
	db := config.DB
	var educational_levels []MstEducationalLevel

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&educational_levels).Error
	if err != nil {
		return nil, err
	}

	return educational_levels, nil
}

func QueryExportEducationalLevels() ([]MstEducationalLevelExport, error) {
	db := config.DB
	var educational_levels []MstEducationalLevelExport

	query := `
        EXEC sp_mst_educational_levels_get 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `

	err := db.Raw(query, "", "name", "asc", 1, CountEducationalLevels()).Scan(&educational_levels).Error
	if err != nil {
		return nil, err
	}

	return educational_levels, nil
}

func QuerySearchEducationalLevels(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEducationalLevelSearch, error) {
	db := config.DB
	var educational_levels []MstEducationalLevelSearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&educational_levels).Error
	if err != nil {
		return nil, err
	}

	return educational_levels, nil
}

func QueryGetEducationalLevel(id string) (MstEducationalLevel, error) {
	db := config.DB
	var job MstEducationalLevel

	query := `
		EXEC sp_mst_educational_levels_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&job).Error
	if err != nil {
		return MstEducationalLevel{}, err
	}

	return job, nil
}

func QueryGetEducationalLevelRelation(id string) (MstEducationalLevelRelation, error) {
	db := config.DB
	var job MstEducationalLevelRelation

	query := `
		EXEC sp_mst_educational_levels_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&job)
	if result.Error != nil {
		return MstEducationalLevelRelation{}, result.Error
	}

	return job, nil
}

func QueryInsertEducationalLevel(id string, code string, name string, description string) error {
	db := config.DB
	now := time.Now()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
		EXEC sp_mst_educational_levels_insert
		@id = ?,
		@code = ?,
		@name = ?,
		@description = ?,
		@created_at = ?,
		@created_by = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, code, name, description, created_at, created_by, updated_at, updated_by).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryUpdateEducationalLevel(id string, code string, name string, description string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
		EXEC sp_mst_educational_levels_update
		@id = ?,
		@code = ?,
		@name = ?,
		@description = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, code, name, description, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func QueryDeleteEducationalLevel(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_educational_levels_delete
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

func QueryRestoreEducationalLevel(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_educational_levels_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
