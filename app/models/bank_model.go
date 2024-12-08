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

type MstBank struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

type MstBankExport struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

type MstBankSearch struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

type MstBankRelation struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

/* Action */
func GetBanks(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstBank, error) {
	return QueryGetBanks("sp_mst_banks_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportBanks(c *fiber.Ctx, fileSaveAs string) error {
	banks, err := QueryExportBanks()
	if err != nil {
		return fmt.Errorf("failed to get banks: %v", err)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.NewSheet(sheetName)

	columns := []string{"A", "B", "C"}
	headers := []string{"ID", "Code", "Name"}

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", col)
		file.SetCellValue(sheetName, cell, headers[i])
	}

	for i, bank := range banks {
		row := i + 2

		values := []interface{}{
			bank.ID,
			bank.Code,
			bank.Name,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(banks))
	}

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchBanks(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstBankSearch, error) {
	return QuerySearchBanks("sp_mst_banks_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetBank(id string) (MstBank, error) {
	return QueryGetBank(id)
}

func CreateBank(id string, code string, name string) error {
	return QueryInsertBank(id, code, name)
}

func ImportBanks(filePath string) error {
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
		var code string = ""
		var name string = ""

		if len(row) > 0 && row[0] != "" {
			id = row[0]
		}
		if len(row) > 1 && row[1] != "" {
			code = row[1]
		}
		if len(row) > 2 && row[2] != "" {
			name = row[2]
		}

		if id != "" {
			exist, err := helpers.CheckModelIDExist(id, &MstBank{})
			if err != nil {
				return err
			}
			if exist {
				if err := QueryUpdateBank(id, code, name); err != nil {
					return err
				}
			} else {
				if err := QueryInsertBank(id, code, name); err != nil {
					return err
				}
			}
		} else {
			id, err := helpers.EnsureUUID(&MstBank{})
			if err != nil {
				return err
			}
			if err := QueryInsertBank(id, code, name); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateBank(id string, code string, name string) error {
	return QueryUpdateBank(id, code, name)
}

func DeleteBank(id string) error {
	if err := QueryDeleteBank(id); err != nil {
		return err
	}

	return nil
}

func GetTrashBanks(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstBank, error) {
	return QueryGetBanks("sp_mst_banks_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreBank(id string) error {
	if err := QueryRestoreBank(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountBanks() int64 {
	return helpers.CountModelSize(&MstBank{}, true)
}

func CountTrashBanks() int64 {
	return helpers.CountModelSize(&MstBank{}, false)
}

/* Query */
func QueryGetBanks(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstBank, error) {
	db := config.DB
	var banks []MstBank

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&banks).Error
	if err != nil {
		return nil, err
	}

	return banks, nil
}

func QueryExportBanks() ([]MstBankExport, error) {
	db := config.DB
	var banks []MstBankExport

	query := `
        EXEC sp_mst_banks_get 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `

	err := db.Raw(query, "", "name", "asc", 1, CountBanks()).Scan(&banks).Error
	if err != nil {
		return nil, err
	}

	return banks, nil
}

func QuerySearchBanks(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstBankSearch, error) {
	db := config.DB
	var banks []MstBankSearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&banks).Error
	if err != nil {
		return nil, err
	}

	return banks, nil
}

func QueryGetBank(id string) (MstBank, error) {
	db := config.DB
	var bank MstBank

	query := `
		EXEC sp_mst_banks_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&bank).Error
	if err != nil {
		return MstBank{}, err
	}

	return bank, nil
}

func QueryGetBankRelation(id string) (MstBankRelation, error) {
	db := config.DB
	var bank MstBankRelation

	query := `
		EXEC sp_mst_banks_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&bank)
	if result.Error != nil {
		return MstBankRelation{}, result.Error
	}

	return bank, nil
}

func QueryInsertBank(id string, code string, name string) error {
	db := config.DB
	now := time.Now()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
		EXEC sp_mst_banks_insert
		@id = ?,
		@code = ?,
		@name = ?,
		@created_at = ?,
		@created_by = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, code, name, created_at, created_by, updated_at, updated_by).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryUpdateBank(id string, code string, name string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
		EXEC sp_mst_banks_update
		@id = ?,
		@code = ?,
		@name = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, code, name, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func QueryDeleteBank(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_banks_delete
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

func QueryRestoreBank(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_banks_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
