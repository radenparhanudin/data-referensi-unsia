package helpers

/* Ensure UUID */
func EnsureUUID(model interface{}) (string, error) {
	for {
		id := GenerateUUID()
		exists, err := CheckModelIDExist(id, model)
		if err != nil {
			return "", err
		}
		if !exists {
			return id, nil
		}
	}
}
