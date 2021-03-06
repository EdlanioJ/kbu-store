package validators

import "github.com/asaskevich/govalidator"

func ValidateUUIDV4(field, str string) error {
	var mapTemplate = map[string]interface{}{
		field: "uuidv4,,optional",
	}

	var inputMap = map[string]interface{}{
		field: str,
	}

	_, err := govalidator.ValidateMap(inputMap, mapTemplate)

	return err
}
