package validators

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

func ValidateLongitude(lng float64) error {
	var mapTemplate = map[string]interface{}{
		"longitude": "longitude",
	}

	var inputMap = map[string]interface{}{
		"longitude": fmt.Sprintf("%f", lng),
	}

	_, err := govalidator.ValidateMap(inputMap, mapTemplate)
	return err
}
