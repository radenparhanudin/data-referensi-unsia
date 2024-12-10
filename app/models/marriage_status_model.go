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

type MstMarriageStatus struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type MstMarriageStatusExport struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type MstMarriageStatusSearch struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type MstMarriageStatusRelation struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

/* Action */
func GetMarriageStatuses(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstMarriageStatus, error) {
	return QueryGetMarriageStatuses("sp_mst_marriage_statuses_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportMarriageStatuses(c *fiber.Ctx, fileSaveAs string) error {
	marriage_statues, err := QueryExportMarriageStatuses()
	if err != nil {
		return fmt.Errorf("failed to get marriage statues: %v", err)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.NewSheet(sheetName)

	columns := []string{"A", "B"}
	headers := []string{"ID", "Name"}

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", col)
		file.SetCellValue(sheetName, cell, headers[i])
	}

	for i, religion := range marriage_statues {
		row := i + 2

		values := []interface{}{
			religion.ID,
			religion.Name,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(marriage_statues))
	}

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchMarriageStatuses(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstMarriageStatusSearch, error) {
	return QuerySearchMarriageStatuses("sp_mst_marriage_statuses_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetMarriageStatus(id string) (MstMarriageStatus, error) {
	return QueryGetMarriageStatus(id)
}

func CreateMarriageStatus(id string, name string) error {
	return QueryInsertMarriageStatus(id, name)
}

func ImportMarriageStatuses(filePath string) error {
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

		if len(row) > 0 && row[0] != "" {
			id = row[0]
		}
		if len(row) > 1 && row[1] != "" {
			name = row[1]
		}

		if id != "" {
			exist, err := helpers.CheckModelIDExist(id, &MstMarriageStatus{})
			if err != nil {
				return err
			}
			if exist {
				if err := QueryUpdateMarriageStatus(id, name); err != nil {
					return err
				}
			} else {
				if err := QueryInsertMarriageStatus(id, name); err != nil {
					return err
				}
			}
		} else {
			id, err := helpers.EnsureUUID(&MstMarriageStatus{})
			if err != nil {
				return err
			}
			if err := QueryInsertMarriageStatus(id, name); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateMarriageStatus(id string, name string) error {
	return QueryUpdateMarriageStatus(id, name)
}

func DeleteMarriageStatus(id string) error {
	if err := QueryDeleteMarriageStatus(id); err != nil {
		return err
	}

	return nil
}

func GetTrashMarriageStatuses(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstMarriageStatus, error) {
	return QueryGetMarriageStatuses("sp_mst_marriage_statuses_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreMarriageStatus(id string) error {
	if err := QueryRestoreMarriageStatus(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountMarriageStatuses() int64 {
	return helpers.CountModelSize(&MstMarriageStatus{}, true)
}

func CountTrashMarriageStatuses() int64 {
	return helpers.CountModelSize(&MstMarriageStatus{}, false)
}

/* Query */
func QueryGetMarriageStatuses(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstMarriageStatus, error) {
	db := config.DB
	var marriage_statues []MstMarriageStatus

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&marriage_statues).Error
	if err != nil {
		return nil, err
	}

	return marriage_statues, nil
}

func QueryExportMarriageStatuses() ([]MstMarriageStatusExport, error) {
	db := config.DB
	var marriage_statues []MstMarriageStatusExport

	query := `
        EXEC sp_mst_marriage_statuses_get 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `

	err := db.Raw(query, "", "name", "asc", 1, CountMarriageStatuses()).Scan(&marriage_statues).Error
	if err != nil {
		return nil, err
	}

	return marriage_statues, nil
}

func QuerySearchMarriageStatuses(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstMarriageStatusSearch, error) {
	db := config.DB
	var marriage_statues []MstMarriageStatusSearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&marriage_statues).Error
	if err != nil {
		return nil, err
	}

	return marriage_statues, nil
}

func QueryGetMarriageStatus(id string) (MstMarriageStatus, error) {
	db := config.DB
	var religion MstMarriageStatus

	query := `
		EXEC sp_mst_marriage_statuses_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&religion).Error
	if err != nil {
		return MstMarriageStatus{}, err
	}

	return religion, nil
}

func QueryGetMarriageStatusRelation(id string) (MstMarriageStatusRelation, error) {
	db := config.DB
	var religion MstMarriageStatusRelation

	query := `
		EXEC sp_mst_marriage_statuses_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&religion)
	if result.Error != nil {
		return MstMarriageStatusRelation{}, result.Error
	}

	return religion, nil
}

func QueryInsertMarriageStatus(id string, name string) error {
	db := config.DB
	now := time.Now()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
		EXEC sp_mst_marriage_statuses_insert
		@id = ?,
		@name = ?,
		@created_at = ?,
		@created_by = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, name, created_at, created_by, updated_at, updated_by).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryUpdateMarriageStatus(id string, name string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
		EXEC sp_mst_marriage_statuses_update
		@id = ?,
		@name = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, name, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func QueryDeleteMarriageStatus(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_marriage_statuses_delete
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

func QueryRestoreMarriageStatus(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_marriage_statuses_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
