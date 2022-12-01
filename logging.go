package templateManager

/*
Functions dedicated to logging
*/

import (
	"fmt"
)

// Logs error messages
func logError(format string, a ...any) error {
	if consoleErrors {
		fmt.Println("\033[31m" + fmt.Sprintf(format, a...) + "\033[0m")
	}

	if haltOnErrors {
		return fmt.Errorf(format, a...)
	}

	return nil
}

// Logs warning messages
func logWarning(format string, a ...any) error {
	if consoleWarnings {
		fmt.Println("\033[33m" + fmt.Sprintf(format, a...) + "\033[0m")
	}

	if haltOnWarnings {
		return fmt.Errorf(format, a...)
	}

	return nil
}

// Logs informational messages
func logInformation(information string) {
	fmt.Println("\033[36m" + information + "\033[0m")
}

// Logs success messages
func logSuccess(success string) {
	fmt.Println("\033[32m" + success + "\033[0m")
}