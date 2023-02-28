package system

import "fmt"

// AddToPATH adds a path to the PATH environment variable in a profile.d script
func AddToPATH(path string) error {
	err := WriteStrings([]string{"PATH=" + path + ":$PATH"}, "/etc/profile.d/wbi_python.sh", 0644)
	if err != nil {
		return fmt.Errorf("failed to add to PATH: %w", err)
	}
	return nil
}
