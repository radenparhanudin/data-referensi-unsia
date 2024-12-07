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

type MstCity struct {
	ID         uuid.UUID           `json:"id"`
	ProvinceId string              `json:"province_id"`
	Province   MstProvinceRelation `json:"province"`
	Name       string              `json:"name"`
	Code       string              `json:"code"`
	CreatedAt  int64               `json:"created_at"`
	UpdatedAt  int64               `json:"updated_at"`
}

type MstCityRelation struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}

func CountCities() int64 {
	db := config.DB
	var count int64
	db.Model(&MstCity{}).Where("deleted_at IS NULL").Count(&count)
	return count
}

func AllCities(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstCity, error) {
	db := config.DB
	var cities []MstCity

	query := `
		EXEC sp_mst_cities_get
		@Filter = ?, 
		@SortBy = ?, 
		@SortDirection = ?, 
		@Page = ?, 
		@PageSize = ?
	`
	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&cities).Error
	if err != nil {
		return nil, err
	}

	return cities, nil
}

func ExportCities(c *fiber.Ctx, outputFile string) error {
	cities, err := AllCities("", "name", "asc", 1, CountCities())
	if err != nil {
		return fmt.Errorf("failed to get cities: %v", err)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.NewSheet(sheetName)

	columns := []string{"A", "B", "C", "D"}
	headers := []string{"ID", "Province ID", "Name", "Code"}

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", col)
		file.SetCellValue(sheetName, cell, headers[i])
	}

	for i, city := range cities {
		row := i + 2

		values := []interface{}{
			city.ID,
			city.ProvinceId,
			city.Name,
			city.Code,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(cities))
	}

	if err := file.SaveAs(outputFile); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func CityById(id string) (MstCity, error) {
	db := config.DB
	var city MstCity

	query := `
		EXEC sp_mst_cities_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&city)
	if result.Error != nil {
		return MstCity{}, result.Error
	}

	if result.RowsAffected == 0 {
		return MstCity{}, fmt.Errorf("data with id %s not found", id)
	}

	var province MstProvinceRelation
	provinceQuery := `
		EXEC sp_mst_provinces_get_by_id
		@id = ?
	`
	provinceResult := db.Raw(provinceQuery, city.ProvinceId).Scan(&province)
	if provinceResult.Error != nil {
		return MstCity{}, provinceResult.Error
	}

	city.Province = province

	return city, nil
}

func CreateCity(province_id string, name string, code string) error {
	db := config.DB

	now := time.Now()
	id := uuid.New().String()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
			EXEC sp_mst_cities_insert
			@id = ?,
			@province_id = ?,
			@name = ?,
			@code = ?,
			@created_at = ?,
			@created_by = ?,
			@updated_at = ?,
			@updated_by = ?
	`

	err := db.Exec(query, id, province_id, name, code, created_at, created_by, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func ImportCities(filePath string) error {
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
		var id, province_id, name, code *string
		if len(row) > 0 && row[0] != "" {
			id = &row[0]
		}
		if len(row) > 1 && row[1] != "" {
			province_id = &row[1]
		}
		if len(row) > 2 && row[2] != "" {
			name = &row[2]
		}
		if len(row) > 3 && row[3] != "" {
			code = &row[3]
		}

		if id != nil {
			/* Update */
			query := `
				EXEC sp_mst_cities_update
				@id = ?,
				@province_id = ?,
				@name = ?,
				@code = ?,
				@updated_at = ?,
				@updated_by = ?
			`

			err := db.Exec(query, id, province_id, name, code, updated_at, updated_by).Error
			if err != nil {
				return err
			}
		} else {
			/* Insert */
			id := uuid.New().String()
			query := `
					EXEC sp_mst_cities_insert
					@id = ?,
					@province_id = ?,
					@name = ?,
					@code = ?,
					@created_at = ?,
					@created_by = ?,
					@updated_at = ?,
					@updated_by = ?
				`

			err := db.Exec(query, id, province_id, name, code, created_at, created_by, updated_at, updated_by).Error
			if err != nil {
				print(err.Error())
				return err
			}
		}
	}

	return nil
}

func UpdateCity(id string, province_id string, name string, code string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
			EXEC sp_mst_cities_update
				@id = ?,
				@province_id = ?,
				@name = ?,
				@code = ?,
				@updated_at = ?,
				@updated_by = ?
	`

	err := db.Exec(query, id, province_id, name, code, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteCity(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_cities_delete
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

func TrashCountCities() int64 {
	db := config.DB
	var count int64
	db.Model(&MstCity{}).Where("deleted_at IS NOT NULL").Count(&count)
	return count
}

func TrashAllCities(filter string, sortBy string, sortDirection string, page int, pageSize int) ([]MstCity, error) {
	db := config.DB
	var cities []MstCity

	query := `
		EXEC sp_mst_cities_has_deleted
		@Filter = ?, 
		@SortBy = ?, 
		@SortDirection = ?, 
		@Page = ?, 
		@PageSize = ?
	`
	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&cities).Error
	if err != nil {
		return nil, err
	}

	return cities, nil
}
