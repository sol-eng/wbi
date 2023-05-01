package system

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"

	cmdlog "github.com/sol-eng/wbi/internal/logging"
)

// WriteStrings appends a slice of strings to a file and creates the file if it doesn't exist
func WriteStrings(lines []string, filepath string, perm fs.FileMode, save bool) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range lines {
		_, err := datawriter.WriteString(data + "\n")
		if err != nil {
			return fmt.Errorf("failed to write line: %w", err)
		}
		if save {
			cmdlog.Info("echo \"" + data + "\" " + ">> " + filepath + "\n")
		}
	}

	datawriter.Flush()
	file.Close()

	return nil
}
