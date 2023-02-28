package system

import "fmt"

// AddToPATH adds a path to the PATH environment variable in a profile.d script
func AddToPATH(path string, filename string) error {
	fullFileName := "/etc/profile.d/wbi_" + filename + ".sh"
	err := WriteStrings([]string{"PATH=" + path + ":$PATH"}, fullFileName, 0644)
	if err != nil {
		return fmt.Errorf("failed to add to PATH: %w", err)
	}
	return nil
}
