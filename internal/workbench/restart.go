package workbench

import (
	"fmt"

	"github.com/dpastoor/wbi/internal/system"
)

func RestartRStudioServerAndLauncher() error {
	err := RestartRStudioServer()
	if err != nil {
		return fmt.Errorf("issue restarting rstudio-server: %w", err)
	}
	err = RestartRStudioLauncher()
	if err != nil {
		return fmt.Errorf("issue restarting rstudio-launcher: %w", err)
	}
	return nil
}

func RestartRStudioServer() error {
	err := system.RunCommand("sudo rstudio-server restart")
	if err != nil {
		return fmt.Errorf("issue restarting rstudio-server: %w", err)
	}
	return nil
}

func RestartRStudioLauncher() error {
	err := system.RunCommand("sudo rstudio-launcher restart")
	if err != nil {
		return fmt.Errorf("issue restarting rstudio-launcher: %w", err)
	}
	return nil
}
