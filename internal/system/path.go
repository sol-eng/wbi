package system

import (
	"fmt"
	"os"
)

// AddToPATH adds a path to the PATH environment variable in a profile.d script
func AddToPATH(path string, filename string) error {
	fullFileName := "/etc/profile.d/wbi_" + filename + ".sh"
	err := WriteStrings([]string{"PATH=" + path + ":$PATH"}, fullFileName, 0644)
	if err != nil {
		return fmt.Errorf("failed to add to PATH: %w", err)
	}

	fmt.Println("Added " + path + " to PATH in " + fullFileName)
	return nil
}

func VerifyFileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return true
}
