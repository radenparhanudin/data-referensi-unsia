package helpers

func ResponseMessage(method string, typeMessage ...bool) string {
	var messageType bool
	if len(typeMessage) > 0 {
		messageType = typeMessage[0]
	}

	switch method {
	case "get":
		if messageType {
			return "Get data successfully"
		}
		return "Get data failed"
	case "export":
		if messageType {
			return "Data export was successful"
		}
		return "Data export failed"
	case "insert":
		if messageType {
			return "Insert data successful"
		}
		return "Insert data failed"
	case "import":
		if messageType {
			return "Data import was successful"
		}
		return "Data import failed"
	case "update":
		if messageType {
			return "Data update successful"
		}
		return "Data update failed"
	case "delete":
		if messageType {
			return "Data deletion successful"
		}
		return "Delete data failed"
	case "exist":
		return "Data already exists"
	case "save":
		if messageType {
			return "Save data successfully"
		}
		return "Save data failed"
	default:
		return "Invalid method"
	}
}
