package system

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// DeleteStrings deletes a slice of strings from a file
func DeleteStrings(lines []string, filepath string, perm fs.FileMode) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var bs []byte
	buf := bytes.NewBuffer(bs)

	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		for _, line := range lines {
			if strings.Contains(fileScanner.Text(), line) {
				continue
			}
			_, err = buf.WriteString(fileScanner.Text() + "\n")
			if err != nil {
				return fmt.Errorf("failed to write line: %w", err)
			}
		}
	}

	err = os.WriteFile(filepath, buf.Bytes(), perm)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
