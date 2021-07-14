package validators

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

func ValidateLatitude(lat float64) error {
	var mapTemplate = map[string]interface{}{
		"latitude": "latitude",
	}

	var inputMap = map[string]interface{}{
		"latitude": fmt.Sprintf("%f", lat),
	}
	_, err := govalidator.ValidateMap(inputMap, mapTemplate)
	return err
}
