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

type MstProvince struct {
	ID         uuid.UUID           `json:"id"`
	CountryId  string              `json:"country_id"`
	Country    *MstCountryRelation `json:"country"`
	Name       string              `json:"name"`
	Code       string              `json:"code"`
	RegionCode string              `json:"region_code"`
	CreatedAt  int64               `json:"created_at"`
	UpdatedAt  int64               `json:"updated_at"`
}

type MstProvinceExport struct {
	ID         uuid.UUID `json:"id"`
	CountryId  string    `json:"country_id"`
	Name       string    `json:"name"`
	Code       string    `json:"code"`
	RegionCode string    `json:"region_code"`
}

type MstProvinceSearch struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}
type MstProvinceRelation struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

/* Action */
func GetProvinces(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstProvince, error) {
	return QueryGetProvinces("sp_mst_provinces_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportProvinces(c *fiber.Ctx, fileSaveAs string) error {
	provinces, err := QueryExportProvinces()
	if err != nil {
		return fmt.Errorf("failed to get provinces: %v", err)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.NewSheet(sheetName)

	columns := []string{"A", "B", "C", "D", "E"}
	headers := []string{"ID", "Country ID", "Name", "Code", "Region Code"}

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", col)
		file.SetCellValue(sheetName, cell, headers[i])
	}

	for i, province := range provinces {
		row := i + 2

		values := []interface{}{
			province.ID,
			province.CountryId,
			province.Name,
			province.Code,
			province.RegionCode,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(provinces))
	}

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchProvinces(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstProvinceSearch, error) {
	return QuerySearchProvinces("sp_mst_provinces_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetProvince(id string) (MstProvince, error) {
	return QueryGetProvince(id)
}

func CreateProvince(id string, country_id string, name string, code string, region_code string) error {
	return QueryInsertProvince(id, country_id, name, code, region_code)
}

func ImportProvinces(filePath string) error {
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
		var country_id string = ""
		var name string = ""
		var code string = ""
		var region_code string = ""

		if len(row) > 0 && row[0] != "" {
			id = row[0]
		}
		if len(row) > 1 && row[1] != "" {
			country_id = row[1]
		}
		if len(row) > 2 && row[2] != "" {
			name = row[2]
		}
		if len(row) > 3 && row[3] != "" {
			code = row[3]
		}
		if len(row) > 4 && row[4] != "" {
			region_code = row[4]
		}

		if id != "" {
			exist, err := helpers.CheckModelIDExist(id, &MstProvince{})
			if err != nil {
				return err
			}
			if exist {
				if err := QueryUpdateProvince(id, country_id, name, code, region_code); err != nil {
					return err
				}
			} else {
				id, err := helpers.EnsureUUID(&MstProvince{})
				if err != nil {
					return err
				}
				if err := QueryInsertProvince(id, country_id, name, code, region_code); err != nil {
					return err
				}
			}
		} else {
			if err := QueryInsertProvince(id, country_id, name, code, region_code); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateProvince(id string, country_id string, name string, code string, region_code string) error {
	return QueryUpdateProvince(id, country_id, name, code, region_code)
}

func DeleteProvince(id string) error {
	if err := QueryDeleteProvince(id); err != nil {
		return err
	}

	return nil
}

func GetTrashProvinces(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstProvince, error) {
	return QueryGetProvinces("sp_mst_provinces_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreProvince(id string) error {
	if err := QueryRestoreProvince(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountProvinces() int64 {
	return helpers.CountModelSize(&MstProvince{}, true)
}

func CountTrashProvinces() int64 {
	return helpers.CountModelSize(&MstProvince{}, false)
}

/* Query */
func QueryGetProvinces(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstProvince, error) {
	db := config.DB
	var provinces []MstProvince

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&provinces).Error
	if err != nil {
		return nil, err
	}

	for i := range provinces {
		country, err := QueryGetCountryRelation(provinces[i].CountryId)
		if err != nil {
			return []MstProvince{}, err
		}

		provinces[i].Country = &country
	}

	return provinces, nil
}

func QueryExportProvinces() ([]MstProvinceExport, error) {
	db := config.DB
	var provinces []MstProvinceExport

	query := `
        EXEC sp_mst_provinces_get 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `

	err := db.Raw(query, "", "name", "asc", 1, CountProvinces()).Scan(&provinces).Error
	if err != nil {
		return nil, err
	}

	return provinces, nil
}

func QuerySearchProvinces(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstProvinceSearch, error) {
	db := config.DB
	var provinces []MstProvinceSearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&provinces).Error
	if err != nil {
		return nil, err
	}

	return provinces, nil
}

func QueryGetProvince(id string) (MstProvince, error) {
	db := config.DB
	var province MstProvince

	query := `
		EXEC sp_mst_provinces_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&province).Error
	if err != nil {
		return MstProvince{}, err
	}

	country, err := QueryGetCountryRelation(province.CountryId)
	if err != nil {
		return MstProvince{}, err
	}

	province.Country = &country

	return province, nil
}

func QueryGetProvinceRelation(id string) (MstProvinceRelation, error) {
	db := config.DB
	var province MstProvinceRelation

	query := `
		EXEC sp_mst_provinces_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&province)
	if result.Error != nil {
		return MstProvinceRelation{}, result.Error
	}

	return province, nil
}

func QueryInsertProvince(id string, country_id string, name string, code string, region_code string) error {
	db := config.DB
	now := time.Now()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
		EXEC sp_mst_provinces_insert
		@id = ?,
		@country_id = ?,
		@name = ?,
		@code = ?,
		@region_code = ?,
		@created_at = ?,
		@created_by = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, country_id, name, code, region_code, created_at, created_by, updated_at, updated_by).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryUpdateProvince(id string, country_id string, name string, code string, region_code string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
		EXEC sp_mst_provinces_update
		@id = ?,
		@country_id = ?,
		@name = ?,
		@code = ?,
		@region_code = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, country_id, name, code, region_code, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func QueryDeleteProvince(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_provinces_delete
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

func QueryRestoreProvince(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_provinces_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
