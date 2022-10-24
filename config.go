package templateManager

import (
	"time"
)

var (
	testsShowDetails = false
	testsShowSuccessful = false
	logErrors = true
	logWarnings = true
	dateDefaultDateFormat = "d/m/Y"
	dateDefaultDatetimeFormat = "d/m/Y H:i"
	dateDefaultTimeFormat = "H:i"
	dateLocalTimezone *time.Location = time.FixedZone("UTC", 0)
)

// Sets the default format for the `date` function (default: d/m/Y)
// May be in Go, PHP or Python format
func SetDefaultDateFormat(format string) {
	dateDefaultDateFormat = format
}

// Sets the default format for the `datetime` function (default: d/m/Y H:i)
// May be in Go, PHP or Python format
func SetDefaultDatetimeFormat(format string) {
	dateDefaultDatetimeFormat = format
}

// Sets the default format for the `time` function (default: H:i)
// May be in Go, PHP or Python format
func SetDefaultTimeFormat(format string) {
	dateDefaultTimeFormat = format
}

// Control whether errors are written to the log
func SetErrors(errors bool) {
	logErrors = errors
}

// Sets the default timezone location used by date / time functions (default: UTC)
func SetTimezoneLocation(location time.Location) {
	dateLocalTimezone = &location
}

// Sets the default timezone location used by date / time functions from a string (default: UTC)
func SetTimezoneLocationString(location string) error {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return err
	}

	dateLocalTimezone = loc

	return nil
}

// Sets the default timezone location used by date / time functions to a fixed offset (default: UTC)
func SetTimezoneFixed(name string, offset int) {
	dateLocalTimezone = time.FixedZone(name, offset)
}

// Control whether warnings are written to the log
func SetWarnings(warnings bool) {
	logWarnings = warnings
}