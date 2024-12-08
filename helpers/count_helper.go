package helpers

import "data-referensi/config"

func CountModelSize(model interface{}, nullableDeletedAt bool) int64 {
	db := config.DB
	var count int64
	var where string

	if nullableDeletedAt {
		where = "deleted_at IS NULL"
	} else {
		where = "deleted_at IS NOT NULL"
	}

	db.Model(model).Where(where).Count(&count)
	return count
}
