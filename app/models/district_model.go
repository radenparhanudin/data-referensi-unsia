package models

import (
	"data-referensi/config"
	"data-referensi/helpers"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type MstDistrict struct {
	ID        uuid.UUID        `json:"id"`
	CityId    string           `json:"city_id"`
	City      *MstCityRelation `json:"city"`
	Name      string           `json:"name"`
	Code      string           `json:"code"`
	CreatedAt int64            `json:"created_at"`
	UpdatedAt int64            `json:"updated_at"`
}

type MstDistrictExport struct {
	ID     uuid.UUID `json:"id"`
	CityId string    `json:"city_id"`
	Name   string    `json:"name"`
	Code   string    `json:"code"`
}

type MstDistrictSearch struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

type MstDistrictRelation struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}

/* Action */
func GetDistricts(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstDistrict, error) {
	return QueryGetDistricts("sp_mst_districts_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportDistricts(c *fiber.Ctx, fileSaveAs string) error {
	districts, err := QueryExportDistricts()
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

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchDistricts(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstDistrictSearch, error) {
	return QuerySearchDistricts("sp_mst_districts_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetDistrict(id string) (MstDistrict, error) {
	return QueryGetDistrict(id)
}

func CreateDistrict(id string, city_id string, name string, code string) error {
	return QueryInsertDistrict(id, city_id, name, code)
}

func ImportDistricts(filePath string) error {
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
		var city_id string = ""
		var name string = ""
		var code string = ""

		if len(row) > 0 && row[0] != "" {
			id = row[0]
		}
		if len(row) > 1 && row[1] != "" {
			city_id = row[1]
		}
		if len(row) > 2 && row[2] != "" {
			name = row[2]
		}
		if len(row) > 3 && row[3] != "" {
			code = row[3]
		}
		if id != "" {
			exist, err := helpers.CheckModelIDExist(id, &MstDistrict{})
			if err != nil {
				return err
			}
			if exist {
				log.Print("Update")
				if err := QueryUpdateDistrict(id, city_id, name, code); err != nil {
					return err
				}
			} else {
				log.Print("Create Ada")
				if err := QueryInsertDistrict(id, city_id, name, code); err != nil {
					return err
				}
			}
		} else {
			id, err := helpers.EnsureUUID(&MstDistrict{})
			if err != nil {
				return err
			}
			if err := QueryInsertDistrict(id, city_id, name, code); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateDistrict(id string, city_id string, name string, code string) error {
	return QueryUpdateDistrict(id, city_id, name, code)
}

func DeleteDistrict(id string) error {
	if err := QueryDeleteDistrict(id); err != nil {
		return err
	}

	return nil
}

func GetTrashDistricts(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstDistrict, error) {
	return QueryGetDistricts("sp_mst_districts_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreDistrict(id string) error {
	if err := QueryRestoreDistrict(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountDistricts() int64 {
	return helpers.CountModelSize(&MstDistrict{}, true)
}

func CountTrashDistricts() int64 {
	return helpers.CountModelSize(&MstDistrict{}, false)
}

/* Query */
func QueryGetDistricts(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstDistrict, error) {
	db := config.DB
	var districts []MstDistrict

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&districts).Error
	if err != nil {
		return nil, err
	}

	for i := range districts {
		city, err := QueryGetCityRelation(districts[i].CityId)
		if err != nil {
			return []MstDistrict{}, err
		}

		districts[i].City = &city
	}

	return districts, nil
}

func QueryExportDistricts() ([]MstDistrictExport, error) {
	db := config.DB
	var districts []MstDistrictExport

	query := `
        EXEC sp_mst_districts_get 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `

	err := db.Raw(query, "", "name", "asc", 1, CountDistricts()).Scan(&districts).Error
	if err != nil {
		return nil, err
	}

	return districts, nil
}

func QuerySearchDistricts(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstDistrictSearch, error) {
	db := config.DB
	var districts []MstDistrictSearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&districts).Error
	if err != nil {
		return nil, err
	}

	return districts, nil
}

func QueryGetDistrict(id string) (MstDistrict, error) {
	db := config.DB
	var district MstDistrict

	query := `
		EXEC sp_mst_districts_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&district).Error
	if err != nil {
		return MstDistrict{}, err
	}

	city, err := QueryGetCityRelation(district.CityId)
	if err != nil {
		return MstDistrict{}, err
	}

	district.City = &city

	return district, nil
}

func QueryGetDistrictRelation(id string) (MstDistrictRelation, error) {
	db := config.DB
	var city MstDistrictRelation

	query := `
		EXEC sp_mst_districts_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&city)
	if result.Error != nil {
		return MstDistrictRelation{}, result.Error
	}

	return city, nil
}

func QueryInsertDistrict(id string, city_id string, name string, code string) error {
	db := config.DB
	now := time.Now()
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

func QueryUpdateDistrict(id string, city_id string, name string, code string) error {
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

func QueryDeleteDistrict(id string) error {
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

func QueryRestoreDistrict(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_districts_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
