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

type MstCountry struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	PhoneCode    string    `json:"phone_code"`
	IconFlagPath string    `json:"icon_flag_path"`
	CreatedAt    int64     `json:"created_at"`
	UpdatedAt    int64     `json:"updated_at"`
}

type MstCountrySearch struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
type MstCountryRelation struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	PhoneCode string    `json:"phone_code"`
}

/* Action */
func GetCountries(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstCountry, error) {
	return QueryGetCountries("sp_mst_countries_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportCountries(c *fiber.Ctx, fileSaveAs string) error {
	countries, err := GetCountries("", "name", "asc", 1, CountCountries())
	if err != nil {
		return fmt.Errorf("failed to get countries: %v", err)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.NewSheet(sheetName)

	columns := []string{"A", "B", "C", "D"}
	headers := []string{"ID", "Name", "Phone Code", "Icon Flag Path"}

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", col)
		file.SetCellValue(sheetName, cell, headers[i])
	}

	for i, country := range countries {
		row := i + 2

		values := []interface{}{
			country.ID,
			country.Name,
			country.PhoneCode,
			country.IconFlagPath,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(countries))
	}

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchCountries(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstCountrySearch, error) {
	return QuerySearchCountries("sp_mst_countries_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetCountry(id string) (MstCountry, error) {
	return QueryGetCountry(id)
}

func CreateCountry(id string, name string, phone_code string, icon_flag_path string) error {
	return QueryInsertCountry(id, name, phone_code, icon_flag_path)
}

func ImportCountries(filePath string) error {
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
		var phone_code string = ""
		var icon_flag_path string = ""

		if len(row) > 0 && row[0] != "" {
			id = row[0]
		}
		if len(row) > 1 && row[1] != "" {
			name = row[1]
		}
		if len(row) > 2 && row[2] != "" {
			phone_code = row[2]
		}
		if len(row) > 3 && row[3] != "" {
			icon_flag_path = row[3]
		}

		if id != "" {
			exist, err := helpers.CheckModelIDExist(id, &MstCountry{})
			if err != nil {
				return err
			}
			if exist {
				if err := QueryUpdateCountry(id, name, phone_code, icon_flag_path); err != nil {
					return err
				}
			} else {
				if err := QueryInsertCountry(id, name, phone_code, icon_flag_path); err != nil {
					return err
				}
			}
		} else {
			id, err := helpers.EnsureUUID(&MstCountry{})
			if err != nil {
				return err
			}
			if err := QueryInsertCountry(id, name, phone_code, icon_flag_path); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateCountry(id string, name string, phone_code string, icon_flag_path string) error {
	return QueryUpdateCountry(id, name, phone_code, icon_flag_path)
}

func DeleteCountry(id string) error {
	if err := QueryDeleteCountry(id); err != nil {
		return err
	}

	return nil
}

func GetTrashCountries(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstCountry, error) {
	return QueryGetCountries("sp_mst_countries_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreCountry(id string) error {
	if err := QueryRestoreCountry(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountCountries() int64 {
	return helpers.CountModelSize(&MstCountry{}, true)
}

func CountTrashCountries() int64 {
	return helpers.CountModelSize(&MstCountry{}, false)
}

/* Query */
func QueryGetCountries(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstCountry, error) {
	db := config.DB
	var countries []MstCountry

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&countries).Error
	if err != nil {
		return nil, err
	}

	return countries, nil
}

func QuerySearchCountries(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstCountrySearch, error) {
	db := config.DB
	var countries []MstCountrySearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&countries).Error
	if err != nil {
		return nil, err
	}

	return countries, nil
}

func QueryGetCountry(id string) (MstCountry, error) {
	db := config.DB
	var country MstCountry

	query := `
		EXEC sp_mst_countries_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&country).Error
	if err != nil {
		return MstCountry{}, err
	}

	return country, nil
}

func QueryGetCountryRelation(id string) (MstCountryRelation, error) {
	db := config.DB
	var country MstCountryRelation

	query := `
		EXEC sp_mst_countries_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&country)
	if result.Error != nil {
		return MstCountryRelation{}, result.Error
	}

	return country, nil
}

func QueryInsertCountry(id string, name string, phone_code string, icon_flag_path string) error {
	db := config.DB
	now := time.Now()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
		EXEC sp_mst_countries_insert
		@id = ?,
		@name = ?,
		@phone_code = ?,
		@icon_flag_path = ?,
		@created_at = ?,
		@created_by = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, name, phone_code, icon_flag_path, created_at, created_by, updated_at, updated_by).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryUpdateCountry(id string, name string, phone_code string, icon_flag_path string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
		EXEC sp_mst_countries_update
		@id = ?,
		@name = ?,
		@phone_code = ?,
		@icon_flag_path = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, name, phone_code, icon_flag_path, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func QueryDeleteCountry(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_countries_delete
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

func QueryRestoreCountry(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_countries_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
