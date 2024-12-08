package helpers

import (
	"data-referensi/config"
)

/* Check ID Model Is Exist */
func CheckModelIDExist(id string, model interface{}) (bool, error) {
	db := config.DB
	var count int64

	err := db.Model(model).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

/* Check Model Is Null Deleted At  */
func CheckModelIsNullDeleted(id string, model interface{}) (bool, error) {
	db := config.DB
	var count int64

	err := db.Model(model).Where("deleted_at IS NULL").Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

/* Check Model Is Not Null Deleted At  */
func CheckModelIsNotNullDeleted(id string, model interface{}) (bool, error) {
	db := config.DB
	var count int64

	err := db.Model(model).Where("deleted_at IS NOT NULL").Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

/* Check Model Is Not Found */
func CheckModelIsNotFound(id string, model interface{}) error {
	exist, err := CheckModelIsNullDeleted(id, model)
	if err != nil {
		return err
	}

	if !exist {
		return GenerateEM(id)
	}

	return nil
}
