package prodrivers

import (
	"fmt"
	"os"
	"strings"

	"github.com/sol-eng/wbi/internal/system"
)

func CheckExistingProDrivers() (bool, error) {
	// check if /etc/opt/rstudio/odbcinst.ini exists
	if _, err := os.Stat("/etc/odbcinst.ini"); err == nil {
		odbciniOutput, err := system.RunCommandAndCaptureOutput("cat /etc/odbcinst.ini")
		if err != nil {
			return false, fmt.Errorf("issue checking for /etc/odbcinst.ini: %w", err)
		}

		if strings.Contains(odbciniOutput, "Installer = RStudio Pro Drivers") {
			fmt.Println("\nExisting installation of Posit Pro Drivers detected in /etc/odbcinst.ini.\nSkipping installation of Pro Drivers.")
			return true, nil
		}
	} else {
		return false, nil
	}

	return false, nil
}
