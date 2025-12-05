package funcs

type ValidationSetting struct {
	Title     string
	Key       string
	ErrorType string
}

func CreateValidationTestDataName(setting *ValidationSetting) string {
	if setting.Key != "name" {
		return "Test User"
	}

	if setting.ErrorType == "required" {
		return ""
	}

	return "Test User"
}

func CreateValidationTestDataEmail(setting *ValidationSetting) string {
	if setting.Key != "email" {
		return "testuser@example.com"
	}

	if setting.ErrorType == "required" {
		return ""
	}

	if setting.ErrorType == "email" {
		return "invalid-email"
	}

	return "testuser@example.com"
}

func CreateValidationTestDataPassword(setting *ValidationSetting) string {
	if setting.Key != "password" {
		return "securepassword"
	}

	if setting.ErrorType == "required" {
		return ""
	}

	if setting.ErrorType == "min" {
		return "short"
	}

	return "securepassword"
}
