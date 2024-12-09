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
	ID         uuid.UUID            `json:"id"`
	ProvinceId string               `json:"province_id"`
	Province   *MstProvinceRelation `json:"province"`
	Name       string               `json:"name"`
	Code       string               `json:"code"`
	CreatedAt  int64                `json:"created_at"`
	UpdatedAt  int64                `json:"updated_at"`
}

type MstCityExport struct {
	ID         uuid.UUID `json:"id"`
	ProvinceId string    `json:"province_id"`
	Name       string    `json:"name"`
	Code       string    `json:"code"`
}

type MstCitySearch struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

type MstCityRelation struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}

/* Action */
func GetCities(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstCity, error) {
	return QueryGetCities("sp_mst_cities_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportCities(c *fiber.Ctx, fileSaveAs string) error {
	cities, err := QueryExportCities()
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

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchCities(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstCitySearch, error) {
	return QuerySearchCities("sp_mst_cities_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetCityByProvinceId(province_id string) ([]MstCitySearch, error) {
	return QueryGetCityByProvinceId(province_id)
}

func GetCity(id string) (MstCity, error) {
	return QueryGetCity(id)
}

func CreateCity(id string, province_id string, name string, code string) error {
	return QueryInsertCity(id, province_id, name, code)
}

func ImportCities(filePath string) error {
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
		var province_id string = ""
		var name string = ""
		var code string = ""

		if len(row) > 0 && row[0] != "" {
			id = row[0]
		}
		if len(row) > 1 && row[1] != "" {
			province_id = row[1]
		}
		if len(row) > 2 && row[2] != "" {
			name = row[2]
		}
		if len(row) > 3 && row[3] != "" {
			code = row[3]
		}

		if id != "" {
			exist, err := helpers.CheckModelIDExist(id, &MstCity{})
			if err != nil {
				return err
			}
			if exist {
				if err := QueryUpdateCity(id, province_id, name, code); err != nil {
					return err
				}
			} else {
				if err := QueryInsertCity(id, province_id, name, code); err != nil {
					return err
				}
			}
		} else {
			id, err := helpers.EnsureUUID(&MstCity{})
			if err != nil {
				return err
			}
			if err := QueryInsertCity(id, province_id, name, code); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateCity(id string, province_id string, name string, code string) error {
	return QueryUpdateCity(id, province_id, name, code)
}

func DeleteCity(id string) error {
	if err := QueryDeleteCity(id); err != nil {
		return err
	}

	return nil
}

func GetTrashCities(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstCity, error) {
	return QueryGetCities("sp_mst_cities_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreCity(id string) error {
	if err := QueryRestoreCity(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountCities() int64 {
	return helpers.CountModelSize(&MstCity{}, true)
}

func CountTrashCities() int64 {
	return helpers.CountModelSize(&MstCity{}, false)
}

/* Query */
func QueryGetCities(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstCity, error) {
	db := config.DB
	var cities []MstCity

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&cities).Error
	if err != nil {
		return nil, err
	}

	for i := range cities {
		province, err := QueryGetProvinceRelation(cities[i].ProvinceId)
		if err != nil {
			return []MstCity{}, err
		}

		cities[i].Province = &province
	}

	return cities, nil
}

func QueryExportCities() ([]MstCityExport, error) {
	db := config.DB
	var cities []MstCityExport

	query := `
        EXEC sp_mst_cities_get 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `

	err := db.Raw(query, "", "name", "asc", 1, CountCities()).Scan(&cities).Error
	if err != nil {
		return nil, err
	}

	return cities, nil
}

func QuerySearchCities(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstCitySearch, error) {
	db := config.DB
	var cities []MstCitySearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&cities).Error
	if err != nil {
		return nil, err
	}

	return cities, nil
}

func QueryGetCityByProvinceId(province_id string) ([]MstCitySearch, error) {
	db := config.DB
	var cities []MstCitySearch

	query := `
		EXEC sp_mst_cities_get_by_province_id
		@province_id = ?
	`
	err := db.Raw(query, province_id).Scan(&cities).Error
	if err != nil {
		return []MstCitySearch{}, err
	}

	return cities, nil
}

func QueryGetCity(id string) (MstCity, error) {
	db := config.DB
	var city MstCity

	query := `
		EXEC sp_mst_cities_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&city).Error
	if err != nil {
		return MstCity{}, err
	}

	province, err := QueryGetProvinceRelation(city.ProvinceId)
	if err != nil {
		return MstCity{}, err
	}

	city.Province = &province

	return city, nil
}

func QueryGetCityRelation(id string) (MstCityRelation, error) {
	db := config.DB
	var city MstCityRelation

	query := `
		EXEC sp_mst_cities_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&city)
	if result.Error != nil {
		return MstCityRelation{}, result.Error
	}

	return city, nil
}

func QueryInsertCity(id string, province_id string, name string, code string) error {
	db := config.DB
	now := time.Now()
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

func QueryUpdateCity(id string, province_id string, name string, code string) error {
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

func QueryDeleteCity(id string) error {
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

func QueryRestoreCity(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_cities_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
