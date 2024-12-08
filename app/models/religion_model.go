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

type MstReligion struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

type MstReligionExport struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

type MstReligionSearch struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type MstReligionRelation struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

/* Action */
func GetReligions(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstReligion, error) {
	return QueryGetReligions("sp_mst_religions_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportReligions(c *fiber.Ctx, fileSaveAs string) error {
	religions, err := QueryExportReligions()
	if err != nil {
		return fmt.Errorf("failed to get religions: %v", err)
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

	for i, religion := range religions {
		row := i + 2

		values := []interface{}{
			religion.ID,
			religion.Code,
			religion.Name,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(religions))
	}

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchReligions(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstReligionSearch, error) {
	return QuerySearchReligions("sp_mst_religions_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetReligion(id string) (MstReligion, error) {
	return QueryGetReligion(id)
}

func CreateReligion(id string, code string, name string) error {
	return QueryInsertReligion(id, code, name)
}

func ImportReligions(filePath string) error {
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
			exist, err := helpers.CheckModelIDExist(id, &MstReligion{})
			if err != nil {
				return err
			}
			if exist {
				if err := QueryUpdateReligion(id, code, name); err != nil {
					return err
				}
			} else {
				if err := QueryInsertReligion(id, code, name); err != nil {
					return err
				}
			}
		} else {
			id, err := helpers.EnsureUUID(&MstReligion{})
			if err != nil {
				return err
			}
			if err := QueryInsertReligion(id, code, name); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateReligion(id string, code string, name string) error {
	return QueryUpdateReligion(id, code, name)
}

func DeleteReligion(id string) error {
	if err := QueryDeleteReligion(id); err != nil {
		return err
	}

	return nil
}

func GetTrashReligions(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstReligion, error) {
	return QueryGetReligions("sp_mst_religions_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreReligion(id string) error {
	if err := QueryRestoreReligion(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountReligions() int64 {
	return helpers.CountModelSize(&MstReligion{}, true)
}

func CountTrashReligions() int64 {
	return helpers.CountModelSize(&MstReligion{}, false)
}

/* Query */
func QueryGetReligions(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstReligion, error) {
	db := config.DB
	var religions []MstReligion

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&religions).Error
	if err != nil {
		return nil, err
	}

	return religions, nil
}

func QueryExportReligions() ([]MstReligionExport, error) {
	db := config.DB
	var religions []MstReligionExport

	query := `
        EXEC sp_mst_religions_get 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `

	err := db.Raw(query, "", "name", "asc", 1, CountReligions()).Scan(&religions).Error
	if err != nil {
		return nil, err
	}

	return religions, nil
}

func QuerySearchReligions(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstReligionSearch, error) {
	db := config.DB
	var religions []MstReligionSearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&religions).Error
	if err != nil {
		return nil, err
	}

	return religions, nil
}

func QueryGetReligion(id string) (MstReligion, error) {
	db := config.DB
	var religion MstReligion

	query := `
		EXEC sp_mst_religions_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&religion).Error
	if err != nil {
		return MstReligion{}, err
	}

	return religion, nil
}

func QueryGetReligionRelation(id string) (MstReligionRelation, error) {
	db := config.DB
	var religion MstReligionRelation

	query := `
		EXEC sp_mst_religions_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&religion)
	if result.Error != nil {
		return MstReligionRelation{}, result.Error
	}

	return religion, nil
}

func QueryInsertReligion(id string, code string, name string) error {
	db := config.DB
	now := time.Now()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
		EXEC sp_mst_religions_insert
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

func QueryUpdateReligion(id string, code string, name string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
		EXEC sp_mst_religions_update
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

func QueryDeleteReligion(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_religions_delete
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

func QueryRestoreReligion(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_religions_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
