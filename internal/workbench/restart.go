package workbench

import (
	"fmt"

	"github.com/sol-eng/wbi/internal/system"
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
	err := system.RunCommand("rstudio-server restart")
	if err != nil {
		return fmt.Errorf("issue restarting rstudio-server: %w", err)
	}
	return nil
}

func RestartRStudioLauncher() error {
	err := system.RunCommand("rstudio-launcher restart")
	if err != nil {
		return fmt.Errorf("issue restarting rstudio-launcher: %w", err)
	}
	return nil
}

func StopRStudioServer() error {
	err := system.RunCommand("rstudio-server stop")
	if err != nil {
		return fmt.Errorf("issue stopping rstudio-server: %w", err)
	}
	return nil
}

func StartRStudioServer() error {
	err := system.RunCommand("rstudio-server start")
	if err != nil {
		return fmt.Errorf("issue starting rstudio-server: %w", err)
	}
	return nil
}
