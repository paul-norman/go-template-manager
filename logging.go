package templateManager

/*
Functions dedicated to logging
*/

import (
	"fmt"
)

// Logs error messages
func logError(err string) {
	if logErrors {
		fmt.Println("\033[31m" + err + "\033[0m")
	}
}

// Logs warning messages
func logWarning(warning string) {
	if logWarnings {
		fmt.Println("\033[33m" + warning + "\033[0m")
	}
}

// Logs informational messages
func logInformation(information string) {
	fmt.Println("\033[36m" + information + "\033[0m")
}

// Logs success messages
func logSuccess(success string) {
	fmt.Println("\033[32m" + success + "\033[0m")
}