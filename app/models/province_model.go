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
	ID         uuid.UUID          `json:"id"`
	CountryId  string             `json:"country_id"`
	Country    MstCountryRelation `json:"country"`
	Name       string             `json:"name"`
	Code       string             `json:"code"`
	RegionCode string             `json:"region_code"`
	CreatedAt  int64              `json:"created_at"`
	UpdatedAt  int64              `json:"updated_at"`
}

type MstProvinceRelation struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}

func CountProvinces() int64 {
	db := config.DB
	var count int64
	db.Model(&MstProvince{}).Where("deleted_at IS NULL").Count(&count)
	return count
}

func AllProvinces(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstProvince, error) {
	db := config.DB
	var provinces []MstProvince

	query := `
		EXEC sp_mst_provinces_get
		@Filter = ?, 
		@SortBy = ?, 
		@SortDirection = ?, 
		@Page = ?, 
		@PageSize = ?
	`
	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&provinces).Error
	if err != nil {
		return nil, err
	}

	return provinces, nil
}

func ExportProvinces(c *fiber.Ctx, outputFile string) error {
	provinces, err := AllProvinces("", "name", "asc", 1, CountProvinces())
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

	if err := file.SaveAs(outputFile); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func ProvinceById(id string) (MstProvince, error) {
	db := config.DB
	var province MstProvince

	query := `
		EXEC sp_mst_provinces_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&province)
	if result.Error != nil {
		return MstProvince{}, result.Error
	}

	if result.RowsAffected == 0 {
		return MstProvince{}, fmt.Errorf("data with id %s not found", id)
	}

	var country MstCountryRelation
	countryQuery := `
		EXEC sp_mst_countries_get_by_id
		@id = ?
	`
	countryResult := db.Raw(countryQuery, province.CountryId).Scan(&country)
	if countryResult.Error != nil {
		return MstProvince{}, countryResult.Error
	}

	province.Country = country

	return province, nil
}

func CreateProvince(country_id string, name string, code string, region_code string) error {
	db := config.DB

	now := time.Now()
	id := uuid.New().String()
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

func ImportProvinces(filePath string) error {
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
		var id, country_id, name, code, region_code *string
		if len(row) > 0 && row[0] != "" {
			id = &row[0]
		}
		if len(row) > 1 && row[1] != "" {
			country_id = &row[1]
		}
		if len(row) > 2 && row[2] != "" {
			name = &row[2]
		}
		if len(row) > 3 && row[3] != "" {
			code = &row[3]
		}
		if len(row) > 4 && row[4] != "" {
			region_code = &row[4]
		}

		if id != nil {
			/* Update */
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
		} else {
			/* Insert */
			id := uuid.New().String()
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
				print(err.Error())
				return err
			}
		}
	}

	return nil
}

func UpdateProvince(id string, country_id string, name string, code string, region_code string) error {
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

func DeleteProvince(id string) error {
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

func TrashCountProvinces() int64 {
	db := config.DB
	var count int64
	db.Model(&MstProvince{}).Where("deleted_at IS NULL").Count(&count)
	return count
}
func TrashAllProvinces(filter string, sortBy string, sortDirection string, page int, pageSize int) ([]MstProvince, error) {
	db := config.DB
	var provinces []MstProvince

	query := `
		EXEC sp_mst_provinces_has_deleted
		@Filter = ?, 
		@SortBy = ?, 
		@SortDirection = ?, 
		@Page = ?, 
		@PageSize = ?
	`
	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&provinces).Error
	if err != nil {
		return nil, err
	}

	return provinces, nil
}
