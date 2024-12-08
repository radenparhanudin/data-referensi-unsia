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

type MstAlmamaterSize struct {
	ID         uuid.UUID `json:"id"`
	Code       string    `json:"code"`
	Size       string    `json:"size"`
	ChestSize  string    `json:"chest_size"`
	ArmLength  string    `json:"arm_length"`
	BodyLength string    `json:"body_length"`
	CreatedAt  int64     `json:"created_at"`
	UpdatedAt  int64     `json:"updated_at"`
}

type MstAlmamaterSizeExport struct {
	ID         uuid.UUID `json:"id"`
	Code       string    `json:"code"`
	Size       string    `json:"size"`
	ChestSize  string    `json:"chest_size"`
	ArmLength  string    `json:"arm_length"`
	BodyLength string    `json:"body_length"`
}

type MstAlmamaterSizeSearch struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Size string    `json:"size"`
}

type MstAlmamaterSizeRelation struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Size string    `json:"size"`
}

/* Action */
func GetAlmamaterSizes(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstAlmamaterSize, error) {
	return QueryGetAlmamaterSizes("sp_mst_almamater_sizes_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportAlmamaterSizes(c *fiber.Ctx, fileSaveAs string) error {
	almamaterSizes, err := QueryExportAlmamaterSizes()
	if err != nil {
		return fmt.Errorf("failed to get alamater sizes: %v", err)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.NewSheet(sheetName)

	columns := []string{"A", "B", "C", "D", "E", "F"}
	headers := []string{"ID", "Code", "Size", "Chest Size", "Arm Length", "Body Length"}

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", col)
		file.SetCellValue(sheetName, cell, headers[i])
	}

	for i, alamaterSize := range almamaterSizes {
		row := i + 2

		values := []interface{}{
			alamaterSize.ID,
			alamaterSize.Code,
			alamaterSize.Size,
			alamaterSize.ChestSize,
			alamaterSize.ArmLength,
			alamaterSize.BodyLength,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(almamaterSizes))
	}

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchAlmamaterSizes(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstAlmamaterSizeSearch, error) {
	return QuerySearchAlmamaterSizes("sp_mst_almamater_sizes_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetAlmamaterSize(id string) (MstAlmamaterSize, error) {
	return QueryGetAlmamaterSize(id)
}

func CreateAlmamaterSize(id string, code string, size string, chest_size string, arm_length string, body_length string) error {
	return QueryInsertAlmamaterSize(id, code, size, chest_size, arm_length, body_length)
}

func ImportAlmamaterSizes(filePath string) error {
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
		var size string = ""
		var chest_size string = ""
		var arm_length string = ""
		var body_length string = ""

		if len(row) > 0 && row[0] != "" {
			id = row[0]
		}
		if len(row) > 1 && row[1] != "" {
			size = row[1]
		}
		if len(row) > 2 && row[2] != "" {
			chest_size = row[2]
		}
		if len(row) > 3 && row[3] != "" {
			arm_length = row[3]
		}
		if len(row) > 4 && row[4] != "" {
			body_length = row[4]
		}

		if id != "" {
			exist, err := helpers.CheckModelIDExist(id, &MstAlmamaterSize{})
			if err != nil {
				return err
			}
			if exist {
				if err := QueryUpdateAlmamaterSize(id, code, size, chest_size, arm_length, body_length); err != nil {
					return err
				}
			} else {
				if err := QueryInsertAlmamaterSize(id, code, size, chest_size, arm_length, body_length); err != nil {
					return err
				}
			}
		} else {
			id, err := helpers.EnsureUUID(&MstAlmamaterSize{})
			if err != nil {
				return err
			}
			if err := QueryInsertAlmamaterSize(id, code, size, chest_size, arm_length, body_length); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateAlmamaterSize(id string, code string, size string, chest_size string, arm_length string, body_length string) error {
	return QueryUpdateAlmamaterSize(id, code, size, chest_size, arm_length, body_length)
}

func DeleteAlmamaterSize(id string) error {
	if err := QueryDeleteAlmamaterSize(id); err != nil {
		return err
	}

	return nil
}

func GetTrashAlmamaterSizes(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstAlmamaterSize, error) {
	return QueryGetAlmamaterSizes("sp_mst_almamater_sizes_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreAlmamaterSize(id string) error {
	if err := QueryRestoreAlmamaterSize(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountAlmamaterSizes() int64 {
	return helpers.CountModelSize(&MstAlmamaterSize{}, true)
}

func CountTrashAlmamaterSizes() int64 {
	return helpers.CountModelSize(&MstAlmamaterSize{}, false)
}

/* Query */
func QueryGetAlmamaterSizes(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstAlmamaterSize, error) {
	db := config.DB
	var almamater_sizes []MstAlmamaterSize

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&almamater_sizes).Error
	if err != nil {
		return nil, err
	}

	return almamater_sizes, nil
}

func QueryExportAlmamaterSizes() ([]MstAlmamaterSizeExport, error) {
	db := config.DB
	var almamater_sizes []MstAlmamaterSizeExport

	query := `
        EXEC sp_mst_almamater_sizes_get 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `

	err := db.Raw(query, "", "name", "asc", 1, CountAlmamaterSizes()).Scan(&almamater_sizes).Error
	if err != nil {
		return nil, err
	}

	return almamater_sizes, nil
}

func QuerySearchAlmamaterSizes(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstAlmamaterSizeSearch, error) {
	db := config.DB
	var almamater_sizes []MstAlmamaterSizeSearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&almamater_sizes).Error
	if err != nil {
		return nil, err
	}

	return almamater_sizes, nil
}

func QueryGetAlmamaterSize(id string) (MstAlmamaterSize, error) {
	db := config.DB
	var alamater_size MstAlmamaterSize

	query := `
		EXEC sp_mst_almamater_sizes_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&alamater_size).Error
	if err != nil {
		return MstAlmamaterSize{}, err
	}

	return alamater_size, nil
}

func QueryGetAlmamaterSizeRelation(id string) (MstAlmamaterSizeRelation, error) {
	db := config.DB
	var alamater_size MstAlmamaterSizeRelation

	query := `
		EXEC sp_mst_almamater_sizes_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&alamater_size)
	if result.Error != nil {
		return MstAlmamaterSizeRelation{}, result.Error
	}

	return alamater_size, nil
}

func QueryInsertAlmamaterSize(id string, code string, size string, chest_size string, arm_length string, body_length string) error {
	db := config.DB
	now := time.Now()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
		EXEC sp_mst_almamater_sizes_insert
		@id = ?,
		@code = ?,
		@size = ?,
		@chest_size = ?,
		@arm_length = ?,
		@body_length = ?,
		@created_at = ?,
		@created_by = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, code, size, chest_size, arm_length, body_length, created_at, created_by, updated_at, updated_by).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryUpdateAlmamaterSize(id string, code string, size string, chest_size string, arm_length string, body_length string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
		EXEC sp_mst_almamater_sizes_update
		@id = ?,
		@code = ?,
		@size = ?,
		@chest_size = ?,
		@arm_length = ?,
		@body_length = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, code, size, chest_size, arm_length, body_length, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func QueryDeleteAlmamaterSize(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_almamater_sizes_delete
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

func QueryRestoreAlmamaterSize(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_almamater_sizes_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
