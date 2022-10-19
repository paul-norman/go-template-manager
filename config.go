package templateManager

import (
	"time"
)

var testsShowDetails = false
var testsShowSuccessful = false
var logErrors = true
var logWarnings = true
var dateLocalTimezone *time.Location = time.FixedZone("UTC", 0)
var dateDefaultDateFormat string = "d/m/Y"
var dateDefaultDatetimeFormat string = "d/m/Y H:i"
var dateDefaultTimeFormat string = "H:i"

func SetTimezoneLocation(location time.Location) {
	dateLocalTimezone = &location
}

func SetTimezoneLocationString(location string) error {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return err
	}

	dateLocalTimezone = loc

	return nil
}

func SetTimezoneFixed(name string, offset int) {
	dateLocalTimezone = time.FixedZone(name, offset)
}