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

type MstDistrict struct {
	ID        uuid.UUID       `json:"id"`
	CityId    string          `json:"city_id"`
	City      MstCityRelation `json:"city"`
	Name      string          `json:"name"`
	Code      string          `json:"code"`
	CreatedAt int64           `json:"created_at"`
	UpdatedAt int64           `json:"updated_at"`
}

type MstDistrictRelation struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func CountDistricts() int64 {
	db := config.DB
	var count int64
	db.Model(&MstDistrict{}).Where("deleted_at IS NULL").Count(&count)
	return count
}

func AllDistricts(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstDistrict, error) {
	db := config.DB
	var districts []MstDistrict

	query := `
		EXEC sp_mst_districts_get
		@Filter = ?, 
		@SortBy = ?, 
		@SortDirection = ?, 
		@Page = ?, 
		@PageSize = ?
	`

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&districts).Error
	if err != nil {
		return nil, err
	}

	return districts, nil
}

func ExportDistricts(c *fiber.Ctx, outputFile string) error {
	districts, err := AllDistricts("", "name", "asc", 1, CountDistricts())
	if err != nil {
		return fmt.Errorf("failed to get districts: %v", err)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.NewSheet(sheetName)

	columns := []string{"A", "B", "C", "D"}
	headers := []string{"ID", "City ID", "Name", "Code"}

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", col)
		file.SetCellValue(sheetName, cell, headers[i])
	}

	for i, city := range districts {
		row := i + 2

		values := []interface{}{
			city.ID,
			city.CityId,
			city.Name,
			city.Code,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}

	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(districts))
	}

	if err := file.SaveAs(outputFile); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func DistrictById(id string) (MstDistrict, error) {
	db := config.DB
	var district MstDistrict

	query := `
		EXEC sp_mst_districts_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&district)
	if result.Error != nil {
		return MstDistrict{}, result.Error
	}

	if result.RowsAffected == 0 {
		return MstDistrict{}, fmt.Errorf("data with id %s not found", id)
	}

	var city MstCityRelation
	cityQuery := `
		EXEC sp_mst_cities_get_by_id
		@id = ?
	`
	cityResult := db.Raw(cityQuery, district.CityId).Scan(&city)
	if cityResult.Error != nil {
		return MstDistrict{}, cityResult.Error
	}

	district.City = city

	return district, nil
}

func CreateDistrict(city_id string, name string, code string) error {
	db := config.DB

	now := time.Now()
	id := uuid.New().String()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
			EXEC sp_mst_districts_insert
			@id = ?,
			@city_id = ?,
			@name = ?,
			@code = ?,
			@created_at = ?,
			@created_by = ?,
			@updated_at = ?,
			@updated_by = ?
	`

	err := db.Exec(query, id, city_id, name, code, created_at, created_by, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func ImportDistricts(filePath string) error {
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
		var id, city_id, name, code *string
		if len(row) > 0 && row[0] != "" {
			id = &row[0]
		}
		if len(row) > 1 && row[1] != "" {
			city_id = &row[1]
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
				EXEC sp_mst_districts_update
				@id = ?,
				@city_id = ?,
				@name = ?,
				@code = ?,
				@updated_at = ?,
				@updated_by = ?
			`

			err := db.Exec(query, id, city_id, name, code, updated_at, updated_by).Error
			if err != nil {
				return err
			}
		} else {
			/* Insert */
			id := uuid.New().String()
			query := `
					EXEC sp_mst_districts_insert
					@id = ?,
					@city_id = ?,
					@name = ?,
					@code = ?,
					@created_at = ?,
					@created_by = ?,
					@updated_at = ?,
					@updated_by = ?
				`

			err := db.Exec(query, id, city_id, name, code, created_at, created_by, updated_at, updated_by).Error
			if err != nil {
				print(err.Error())
				return err
			}
		}
	}

	return nil
}

func UpdateDistrict(id string, city_id string, name string, code string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
			EXEC sp_mst_districts_update
				@id = ?,
				@city_id = ?,
				@name = ?,
				@code = ?,
				@updated_at = ?,
				@updated_by = ?
	`

	err := db.Exec(query, id, city_id, name, code, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteDistrict(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_districts_delete
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

func TrashCountDistricts() int64 {
	db := config.DB
	var count int64
	db.Model(&MstDistrict{}).Where("deleted_at IS NOT NULL").Count(&count)
	return count
}

func TrashAllDistricts(filter string, sortBy string, sortDirection string, page int, pageSize int) ([]MstDistrict, error) {
	db := config.DB
	var districts []MstDistrict

	query := `
		EXEC sp_mst_districts_has_deleted
		@Filter = ?, 
		@SortBy = ?, 
		@SortDirection = ?, 
		@Page = ?, 
		@PageSize = ?
	`
	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&districts).Error
	if err != nil {
		return nil, err
	}

	return districts, nil
}
