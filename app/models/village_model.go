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
	ID         uuid.UUID            `json:"id"`
	DistrictId string               `json:"district_id"`
	District   *MstDistrictRelation `json:"district"`
	Name       string               `json:"name"`
	Code       string               `json:"code"`
	CreatedAt  int64                `json:"created_at"`
	UpdatedAt  int64                `json:"updated_at"`
}

type MstVillageExport struct {
	ID         uuid.UUID `json:"id"`
	DistrictId string    `json:"district_id"`
	Name       string    `json:"name"`
	Code       string    `json:"code"`
}
type MstVillageSearch struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

type MstVillageRelation struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}

/* Action */
func GetVillages(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstVillage, error) {
	return QueryGetVillages("sp_mst_villages_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportVillages(c *fiber.Ctx, fileSaveAs string) error {
	villages, err := QueryExportVillages()
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

	for i, district := range villages {
		row := i + 2

		values := []interface{}{
			district.ID,
			district.DistrictId,
			district.Name,
			district.Code,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(villages))
	}

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchVillages(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstVillageSearch, error) {
	return QuerySearchVillages("sp_mst_villages_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetVillageByDistrictId(district_id string) ([]MstVillageSearch, error) {
	return QueryGetVillageByDistrictId(district_id)
}

func GetVillage(id string) (MstVillage, error) {
	return QueryGetVillage(id)
}

func CreateVillage(id string, district_id string, name string, code string) error {
	return QueryInsertVillage(id, district_id, name, code)
}

func ImportVillages(filePath string) error {
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
		var district_id string = ""
		var name string = ""
		var code string = ""

		if len(row) > 0 && row[0] != "" {
			id = row[0]
		}
		if len(row) > 1 && row[1] != "" {
			district_id = row[1]
		}
		if len(row) > 2 && row[2] != "" {
			name = row[2]
		}
		if len(row) > 3 && row[3] != "" {
			code = row[3]
		}

		if id != "" {
			exist, err := helpers.CheckModelIDExist(id, &MstVillage{})
			if err != nil {
				return err
			}
			if exist {
				if err := QueryUpdateVillage(id, district_id, name, code); err != nil {
					return err
				}
			} else {
				if err := QueryInsertVillage(id, district_id, name, code); err != nil {
					return err
				}
			}
		} else {
			id, err := helpers.EnsureUUID(&MstVillage{})
			if err != nil {
				return err
			}
			if err := QueryInsertVillage(id, district_id, name, code); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateVillage(id string, district_id string, name string, code string) error {
	return QueryUpdateVillage(id, district_id, name, code)
}

func DeleteVillage(id string) error {
	if err := QueryDeleteVillage(id); err != nil {
		return err
	}

	return nil
}

func GetTrashVillages(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstVillage, error) {
	return QueryGetVillages("sp_mst_villages_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreVillage(id string) error {
	if err := QueryRestoreVillage(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountVillages() int64 {
	return helpers.CountModelSize(&MstVillage{}, true)
}

func CountTrashVillages() int64 {
	return helpers.CountModelSize(&MstVillage{}, false)
}

/* Query */
func QueryGetVillages(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstVillage, error) {
	db := config.DB
	var villages []MstVillage

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&villages).Error
	if err != nil {
		return nil, err
	}

	for i := range villages {
		district, err := QueryGetDistrictRelation(villages[i].DistrictId)
		if err != nil {
			return []MstVillage{}, err
		}

		villages[i].District = &district
	}

	return villages, nil
}

func QueryExportVillages() ([]MstVillageExport, error) {
	db := config.DB
	var villages []MstVillageExport

	query := `
        EXEC sp_mst_villages_get
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `

	err := db.Raw(query, "", "name", "asc", 1, CountVillages()).Scan(&villages).Error
	if err != nil {
		return nil, err
	}

	return villages, nil
}

func QuerySearchVillages(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstVillageSearch, error) {
	db := config.DB
	var villages []MstVillageSearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&villages).Error
	if err != nil {
		return nil, err
	}

	return villages, nil
}

func QueryGetVillageByDistrictId(district_id string) ([]MstVillageSearch, error) {
	db := config.DB
	var village []MstVillageSearch

	query := `
		EXEC sp_mst_villages_get_by_district_id
		@district_id = ?
	`
	err := db.Raw(query, district_id).Scan(&village).Error
	if err != nil {
		return []MstVillageSearch{}, err
	}

	return village, nil
}

func QueryGetVillage(id string) (MstVillage, error) {
	db := config.DB
	var village MstVillage

	query := `
		EXEC sp_mst_villages_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&village).Error
	if err != nil {
		return MstVillage{}, err
	}

	district, err := QueryGetDistrictRelation(village.DistrictId)
	if err != nil {
		return MstVillage{}, err
	}

	village.District = &district

	return village, nil
}

func QueryGetVillageRelation(id string) (MstVillageRelation, error) {
	db := config.DB
	var district MstVillageRelation

	query := `
		EXEC sp_mst_villages_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&district)
	if result.Error != nil {
		return MstVillageRelation{}, result.Error
	}

	return district, nil
}

func QueryInsertVillage(id string, district_id string, name string, code string) error {
	db := config.DB
	now := time.Now()
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

func QueryUpdateVillage(id string, district_id string, name string, code string) error {
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

func QueryDeleteVillage(id string) error {
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

func QueryRestoreVillage(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_villages_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
