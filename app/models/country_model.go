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
	// CreatedAtString string    `json:"created_at_string"`
	// UpdatedAtString string    `json:"updated_at_string"`
}

type MstCountryRelation struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	PhoneCode string    `json:"phone_code"`
}

func CountCountries() int64 {
	db := config.DB
	var count int64
	db.Model(&MstCountry{}).Where("deleted_at IS NULL").Count(&count)
	return count
}

func AllCountries(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstCountry, error) {
	db := config.DB
	var countries []MstCountry

	query := `
		EXEC sp_mst_countries_get
		@Filter = ?, 
		@SortBy = ?, 
		@SortDirection = ?, 
		@Page = ?, 
		@PageSize = ?
	`
	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&countries).Error
	if err != nil {
		return nil, err
	}

	// for i := range countries {
	// 	countries[i].CreatedAtString = time.UnixMilli(countries[i].CreatedAt).Format("2006-01-02 15:04:05")
	// 	countries[i].UpdatedAtString = time.UnixMilli(countries[i].UpdatedAt).Format("2006-01-02 15:04:05")
	// }

	return countries, nil
}

func ExportCountries(c *fiber.Ctx, outputFile string) error {
	countries, err := AllCountries("", "name", "asc", 1, CountCountries())
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

	if err := file.SaveAs(outputFile); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func CountryById(id string) (MstCountry, error) {
	db := config.DB
	var country MstCountry

	query := `
		EXEC sp_mst_countries_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&country)
	if result.Error != nil {
		return MstCountry{}, result.Error
	}

	if result.RowsAffected == 0 {
		return MstCountry{}, fmt.Errorf("data with id %s not found", id)
	}
	return country, nil
}

func CreateCountry(name string, phone_code string, icon_flag_path string) error {
	db := config.DB

	now := time.Now()
	id := uuid.New().String()
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

func ImportCountries(filePath string) error {
	db := config.DB
	now := time.Now()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

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
		var id, name, phone_code, icon_flag_path *string
		if len(row) > 0 && row[0] != "" {
			id = &row[0]
		}
		if len(row) > 1 && row[1] != "" {
			name = &row[1]
		}
		if len(row) > 2 && row[2] != "" {
			phone_code = &row[2]
		}
		if len(row) > 3 && row[3] != "" {
			icon_flag_path = &row[3]
		}

		if id != nil {
			/* Update */
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
		} else {
			/* Insert */
			id := uuid.New().String()
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
				print(err.Error())
				return err
			}
		}
	}

	return nil
}

func UpdateCountry(id string, name string, phone_code string, icon_flag_path string) error {
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

func DeleteCountry(id string) error {
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

func TrashCountCountries() int64 {
	db := config.DB
	var count int64
	db.Model(&MstCountry{}).Where("deleted_at IS NOT NULL").Count(&count)
	return count
}

func TrashAllCountries(filter string, sortBy string, sortDirection string, page int, pageSize int) ([]MstCountry, error) {
	db := config.DB
	var countries []MstCountry

	query := `
		EXEC sp_mst_countries_has_deleted
		@Filter = ?, 
		@SortBy = ?, 
		@SortDirection = ?, 
		@Page = ?, 
		@PageSize = ?
	`
	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&countries).Error
	if err != nil {
		return nil, err
	}

	return countries, nil
}
