package system

import (
	"bufio"
	"fmt"
	"os"
)

// WriteStrings appends a slice of strings to a file and creates the file if it doesn't exist
func WriteStrings(lines []string, filepath string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range lines {
		_, err := datawriter.WriteString(data + "\n")
		if err != nil {
			return fmt.Errorf("failed to write line: %w", err)
		}
	}

	datawriter.Flush()
	file.Close()

	return nil
}
