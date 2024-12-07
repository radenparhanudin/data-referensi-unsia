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

type MstVillage struct {
	ID         uuid.UUID           `json:"id"`
	DistrictId string              `json:"district_id"`
	District   MstDistrictRelation `json:"district"`
	Name       string              `json:"name"`
	Code       string              `json:"code"`
	CreatedAt  int64               `json:"created_at"`
	UpdatedAt  int64               `json:"updated_at"`
}

type MstVillageRelation struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}

func CountVillages() int64 {
	db := config.DB
	var count int64
	db.Model(&MstVillage{}).Where("deleted_at IS NULL").Count(&count)
	return count
}

func AllVillages(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstVillage, error) {
	db := config.DB
	var villages []MstVillage

	query := `
		EXEC sp_mst_villages_get
		@Filter = ?, 
		@SortBy = ?, 
		@SortDirection = ?, 
		@Page = ?, 
		@PageSize = ?
	`
	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&villages).Error
	if err != nil {
		return nil, err
	}

	return villages, nil
}

func ExportVillages(c *fiber.Ctx, outputFile string) error {
	villages, err := AllVillages("", "name", "asc", 1, CountVillages())
	if err != nil {
		return fmt.Errorf("failed to get villages: %v", err)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.NewSheet(sheetName)

	columns := []string{"A", "B", "C", "D"}
	headers := []string{"ID", "District ID", "Name", "Code"}

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", col)
		file.SetCellValue(sheetName, cell, headers[i])
	}

	for i, city := range villages {
		row := i + 2

		values := []interface{}{
			city.ID,
			city.DistrictId,
			city.Name,
			city.Code,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(villages))
	}

	if err := file.SaveAs(outputFile); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func VillageById(id string) (MstVillage, error) {
	db := config.DB
	var city MstVillage

	query := `
		EXEC sp_mst_villages_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&city)
	if result.Error != nil {
		return MstVillage{}, result.Error
	}

	if result.RowsAffected == 0 {
		return MstVillage{}, fmt.Errorf("data with id %s not found", id)
	}

	var district MstDistrictRelation
	districtQuery := `
		EXEC sp_mst_districts_get_by_id
		@id = ?
	`
	districtResult := db.Raw(districtQuery, city.DistrictId).Scan(&district)
	if districtResult.Error != nil {
		return MstVillage{}, districtResult.Error
	}

	city.District = district

	return city, nil
}

func CreateVillage(district_id string, name string, code string) error {
	db := config.DB

	now := time.Now()
	id := uuid.New().String()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
			EXEC sp_mst_villages_insert
			@id = ?,
			@district_id = ?,
			@name = ?,
			@code = ?,
			@created_at = ?,
			@created_by = ?,
			@updated_at = ?,
			@updated_by = ?
	`

	err := db.Exec(query, id, district_id, name, code, created_at, created_by, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func ImportVillages(filePath string) error {
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
		var id, district_id, name, code *string
		if len(row) > 0 && row[0] != "" {
			id = &row[0]
		}
		if len(row) > 1 && row[1] != "" {
			district_id = &row[1]
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
				EXEC sp_mst_villages_update
				@id = ?,
				@district_id = ?,
				@name = ?,
				@code = ?,
				@updated_at = ?,
				@updated_by = ?
			`

			err := db.Exec(query, id, district_id, name, code, updated_at, updated_by).Error
			if err != nil {
				return err
			}
		} else {
			/* Insert */
			id := uuid.New().String()
			query := `
					EXEC sp_mst_villages_insert
					@id = ?,
					@district_id = ?,
					@name = ?,
					@code = ?,
					@created_at = ?,
					@created_by = ?,
					@updated_at = ?,
					@updated_by = ?
				`

			err := db.Exec(query, id, district_id, name, code, created_at, created_by, updated_at, updated_by).Error
			if err != nil {
				print(err.Error())
				return err
			}
		}
	}

	return nil
}

func UpdateVillage(id string, district_id string, name string, code string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
			EXEC sp_mst_villages_update
				@id = ?,
				@district_id = ?,
				@name = ?,
				@code = ?,
				@updated_at = ?,
				@updated_by = ?
	`

	err := db.Exec(query, id, district_id, name, code, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteVillage(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_villages_delete
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

func TrashCountVillages() int64 {
	db := config.DB
	var count int64
	db.Model(&MstVillage{}).Where("deleted_at IS NOT NULL").Count(&count)
	return count
}

func TrashAllVillages(filter string, sortBy string, sortDirection string, page int, pageSize int) ([]MstVillage, error) {
	db := config.DB
	var villages []MstVillage

	query := `
		EXEC sp_mst_villages_has_deleted
		@Filter = ?, 
		@SortBy = ?, 
		@SortDirection = ?, 
		@Page = ?, 
		@PageSize = ?
	`
	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&villages).Error
	if err != nil {
		return nil, err
	}

	return villages, nil
}
