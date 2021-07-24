package validators

import "github.com/asaskevich/govalidator"

func ValidateRequired(field, str string) error {
	var mapTemplate = map[string]interface{}{
		field: "required",
	}

	var inputMap = map[string]interface{}{
		field: str,
	}

	_, err := govalidator.ValidateMap(inputMap, mapTemplate)
	return err
}
