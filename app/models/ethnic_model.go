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

type MstEthnic struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	RegionOfOrigin string    `json:"region_of_origin"`
	CreatedAt      int64     `json:"created_at"`
	UpdatedAt      int64     `json:"updated_at"`
}

type MstEthnicExport struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	RegionOfOrigin string    `json:"region_of_origin"`
}

type MstEthnicSearch struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type MstEthnicRelation struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

/* Action */
func GetEthnics(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEthnic, error) {
	return QueryGetEthnics("sp_mst_ethnics_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportEthnics(c *fiber.Ctx, fileSaveAs string) error {
	ethnics, err := QueryExportEthnics()
	if err != nil {
		return fmt.Errorf("failed to get ethnics: %v", err)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.NewSheet(sheetName)

	columns := []string{"A", "B", "C"}
	headers := []string{"ID", "Name", "RegionOfOrigin"}

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", col)
		file.SetCellValue(sheetName, cell, headers[i])
	}

	for i, ethnic := range ethnics {
		row := i + 2

		values := []interface{}{
			ethnic.ID,
			ethnic.Name,
			ethnic.RegionOfOrigin,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(ethnics))
	}

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchEthnics(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEthnicSearch, error) {
	return QuerySearchEthnics("sp_mst_ethnics_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetEthnic(id string) (MstEthnic, error) {
	return QueryGetEthnic(id)
}

func CreateEthnic(id string, name string, region_of_origin string) error {
	return QueryInsertEthnic(id, name, region_of_origin)
}

func ImportEthnics(filePath string) error {
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
		var region_of_origin string = ""

		if len(row) > 0 && row[0] != "" {
			id = row[0]
		}
		if len(row) > 1 && row[1] != "" {
			name = row[1]
		}
		if len(row) > 2 && row[2] != "" {
			region_of_origin = row[2]
		}

		if id != "" {
			exist, err := helpers.CheckModelIDExist(id, &MstEthnic{})
			if err != nil {
				return err
			}
			if exist {
				if err := QueryUpdateEthnic(id, name, region_of_origin); err != nil {
					return err
				}
			} else {
				if err := QueryInsertEthnic(id, name, region_of_origin); err != nil {
					return err
				}
			}
		} else {
			id, err := helpers.EnsureUUID(&MstEthnic{})
			if err != nil {
				return err
			}
			if err := QueryInsertEthnic(id, name, region_of_origin); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateEthnic(id string, name string, region_of_origin string) error {
	return QueryUpdateEthnic(id, name, region_of_origin)
}

func DeleteEthnic(id string) error {
	if err := QueryDeleteEthnic(id); err != nil {
		return err
	}

	return nil
}

func GetTrashEthnics(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEthnic, error) {
	return QueryGetEthnics("sp_mst_ethnics_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreEthnic(id string) error {
	if err := QueryRestoreEthnic(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountEthnics() int64 {
	return helpers.CountModelSize(&MstEthnic{}, true)
}

func CountTrashEthnics() int64 {
	return helpers.CountModelSize(&MstEthnic{}, false)
}

/* Query */
func QueryGetEthnics(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEthnic, error) {
	db := config.DB
	var ethnics []MstEthnic

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&ethnics).Error
	if err != nil {
		return nil, err
	}

	return ethnics, nil
}

func QueryExportEthnics() ([]MstEthnicExport, error) {
	db := config.DB
	var ethnics []MstEthnicExport

	query := `
        EXEC sp_mst_ethnics_get 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `

	err := db.Raw(query, "", "name", "asc", 1, CountEthnics()).Scan(&ethnics).Error
	if err != nil {
		return nil, err
	}

	return ethnics, nil
}

func QuerySearchEthnics(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstEthnicSearch, error) {
	db := config.DB
	var ethnics []MstEthnicSearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&ethnics).Error
	if err != nil {
		return nil, err
	}

	return ethnics, nil
}

func QueryGetEthnic(id string) (MstEthnic, error) {
	db := config.DB
	var ethnic MstEthnic

	query := `
		EXEC sp_mst_ethnics_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&ethnic).Error
	if err != nil {
		return MstEthnic{}, err
	}

	return ethnic, nil
}

func QueryGetEthnicRelation(id string) (MstEthnicRelation, error) {
	db := config.DB
	var ethnic MstEthnicRelation

	query := `
		EXEC sp_mst_ethnics_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&ethnic)
	if result.Error != nil {
		return MstEthnicRelation{}, result.Error
	}

	return ethnic, nil
}

func QueryInsertEthnic(id string, name string, region_of_origin string) error {
	db := config.DB
	now := time.Now()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
		EXEC sp_mst_ethnics_insert
		@id = ?,
		@name = ?,
		@region_of_origin = ?,
		@created_at = ?,
		@created_by = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, name, region_of_origin, created_at, created_by, updated_at, updated_by).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryUpdateEthnic(id string, name string, region_of_origin string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
		EXEC sp_mst_ethnics_update
		@id = ?,
		@name = ?,
		@region_of_origin = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, name, region_of_origin, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func QueryDeleteEthnic(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_ethnics_delete
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

func QueryRestoreEthnic(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_ethnics_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
