package models

import (
	"data-referensi/config"

	"github.com/google/uuid"
)

type MstCountry struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	PhoneCode    string    `json:"phone_code"`
	IconFlagPath string    `json:"icon_flag_path"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

type CreateOrUpdateCountry struct {
	Name         string `json:"name"  validate:"required"`
	PhoneCode    string `json:"phone_code" validate:"required"`
	IconFlagPath string `json:"icon_flag_path"   validate:"required"`
}

func GetCountries(filter string, sortBy string, sortDirection string, page int, pageSize int) ([]MstCountry, error) {
	db := config.DB
	var countries []MstCountry

	query := `
		EXEC sp_mst_countries_get
		@Filter = ?, 
		@SortBy = ?, 
		@SortDirection = ?, 
		@Page = ?, 
		@PageSize = ?
	`
	err := db.Raw(query, filter, sortBy, sortDirection, page, pageSize).Scan(&countries).Error
	if err != nil {
		return nil, err
	}

	return countries, nil
}
