package system

import (
	"errors"
	"os"
	"strings"
)

// CheckStringExists checks if a string exists in a file
func CheckStringExists(text string, filepath string) (bool, error) {
	if _, err := os.Stat(filepath); err == nil {
		file, err := os.ReadFile(filepath)
		if err != nil {
			return false, errors.New("failed to read file")
		}
		if strings.Contains(string(file), text) {
			return true, nil
		}
	} else {
		return false, errors.New("failed to check if file exists")
	}
	return false, nil
}
