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

type MstJob struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   int64     `json:"created_at"`
	UpdatedAt   int64     `json:"updated_at"`
}

type MstJobExport struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type MstJobSearch struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

type MstJobRelation struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

/* Action */
func GetJobs(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstJob, error) {
	return QueryGetJobs("sp_mst_jobs_get", filter, sortBy, sortDirection, page, pageSize)
}

func ExportJobs(c *fiber.Ctx, fileSaveAs string) error {
	jobs, err := QueryExportJobs()
	if err != nil {
		return fmt.Errorf("failed to get jobs: %v", err)
	}

	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.NewSheet(sheetName)

	columns := []string{"A", "B", "C", "D"}
	headers := []string{"ID", "Code", "Name", "Description"}

	for i, col := range columns {
		cell := fmt.Sprintf("%s1", col)
		file.SetCellValue(sheetName, cell, headers[i])
	}

	for i, job := range jobs {
		row := i + 2

		values := []interface{}{
			job.ID,
			job.Code,
			job.Name,
			job.Description,
		}

		for i, col := range columns {
			cell := fmt.Sprintf("%s%d", col, row)
			file.SetCellValue(sheetName, cell, values[i])
		}
	}

	for _, col := range columns {
		helpers.ExcelAutoSizeColumn(file, sheetName, col, len(jobs))
	}

	if err := file.SaveAs(fileSaveAs); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	return nil
}

func SearchJobs(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstJobSearch, error) {
	return QuerySearchJobs("sp_mst_jobs_get", filter, sortBy, sortDirection, page, pageSize)
}

func GetJob(id string) (MstJob, error) {
	return QueryGetJob(id)
}

func CreateJob(id string, code string, name string, description string) error {
	return QueryInsertJob(id, code, name, description)
}

func ImportJobs(filePath string) error {
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
		var description string = ""

		if len(row) > 0 && row[0] != "" {
			id = row[0]
		}
		if len(row) > 1 && row[1] != "" {
			code = row[1]
		}
		if len(row) > 2 && row[2] != "" {
			name = row[2]
		}
		if len(row) > 3 && row[3] != "" {
			description = row[3]
		}

		if id != "" {
			exist, err := helpers.CheckModelIDExist(id, &MstJob{})
			if err != nil {
				return err
			}
			if exist {
				if err := QueryUpdateJob(id, code, name, description); err != nil {
					return err
				}
			} else {
				if err := QueryInsertJob(id, code, name, description); err != nil {
					return err
				}
			}
		} else {
			id, err := helpers.EnsureUUID(&MstJob{})
			if err != nil {
				return err
			}
			if err := QueryInsertJob(id, code, name, description); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateJob(id string, code string, name string, description string) error {
	return QueryUpdateJob(id, code, name, description)
}

func DeleteJob(id string) error {
	if err := QueryDeleteJob(id); err != nil {
		return err
	}

	return nil
}

func GetTrashJobs(filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstJob, error) {
	return QueryGetJobs("sp_mst_jobs_has_deleted", filter, sortBy, sortDirection, page, pageSize)
}

func RestoreJob(id string) error {
	if err := QueryRestoreJob(id); err != nil {
		return err
	}

	return nil
}

/* Count */
func CountJobs() int64 {
	return helpers.CountModelSize(&MstJob{}, true)
}

func CountTrashJobs() int64 {
	return helpers.CountModelSize(&MstJob{}, false)
}

/* Query */
func QueryGetJobs(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstJob, error) {
	db := config.DB
	var jobs []MstJob

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&jobs).Error
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func QueryExportJobs() ([]MstJobExport, error) {
	db := config.DB
	var jobs []MstJobExport

	query := `
        EXEC sp_mst_jobs_get 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `

	err := db.Raw(query, "", "name", "asc", 1, CountJobs()).Scan(&jobs).Error
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func QuerySearchJobs(sp string, filter string, sortBy string, sortDirection string, page int, pageSize int64) ([]MstJobSearch, error) {
	db := config.DB
	var jobs []MstJobSearch

	query := fmt.Sprintf(`
        EXEC %s 
        @Filter = ?, 
        @SortBy = ?, 
        @SortDirection = ?, 
        @Page = ?, 
        @PageSize = ?
    `, sp)

	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&jobs).Error
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func QueryGetJob(id string) (MstJob, error) {
	db := config.DB
	var job MstJob

	query := `
		EXEC sp_mst_jobs_get_by_id
		@id = ?
	`
	err := db.Raw(query, id).Scan(&job).Error
	if err != nil {
		return MstJob{}, err
	}

	return job, nil
}

func QueryGetJobRelation(id string) (MstJobRelation, error) {
	db := config.DB
	var job MstJobRelation

	query := `
		EXEC sp_mst_jobs_get_by_id
		@id = ?
	`
	result := db.Raw(query, id).Scan(&job)
	if result.Error != nil {
		return MstJobRelation{}, result.Error
	}

	return job, nil
}

func QueryInsertJob(id string, code string, name string, description string) error {
	db := config.DB
	now := time.Now()
	created_at := now.UnixMilli()
	updated_at := now.UnixMilli()
	var created_by, updated_by *string = nil, nil

	query := `
		EXEC sp_mst_jobs_insert
		@id = ?,
		@code = ?,
		@name = ?,
		@description = ?,
		@created_at = ?,
		@created_by = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, code, name, description, created_at, created_by, updated_at, updated_by).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryUpdateJob(id string, code string, name string, description string) error {
	db := config.DB

	now := time.Now()
	updated_at := now.UnixMilli()
	var updated_by *string = nil

	query := `
		EXEC sp_mst_jobs_update
		@id = ?,
		@code = ?,
		@name = ?,
		@description = ?,
		@updated_at = ?,
		@updated_by = ?
	`

	err := db.Exec(query, id, code, name, description, updated_at, updated_by).Error
	if err != nil {
		return err
	}

	return nil
}

func QueryDeleteJob(id string) error {
	db := config.DB

	now := time.Now()
	deleted_at := now.UnixMilli()
	var deleted_by *string = nil

	query := `
		EXEC sp_mst_jobs_delete
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

func QueryRestoreJob(id string) error {
	db := config.DB
	query := `
		EXEC sp_mst_jobs_restore
		@id = ?
	`

	err := db.Exec(query, id).Error
	if err != nil {
		return err
	}

	return nil
}
