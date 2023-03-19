package cmd

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/sol-eng/wbi/internal/authentication"
	"github.com/sol-eng/wbi/internal/connect"
	"github.com/sol-eng/wbi/internal/jupyter"
	"github.com/sol-eng/wbi/internal/languages"
	"github.com/sol-eng/wbi/internal/license"
	"github.com/sol-eng/wbi/internal/operatingsystem"
	"github.com/sol-eng/wbi/internal/packagemanager"
	"github.com/sol-eng/wbi/internal/prodrivers"
	"github.com/sol-eng/wbi/internal/ssl"
	"github.com/sol-eng/wbi/internal/workbench"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type setupCmd struct {
	cmd  *cobra.Command
	opts setupOpts
}

type setupOpts struct {
}

func newSetup(setupOpts setupOpts) error {

	fmt.Println("Welcome to the Workbench Installer!\n")

	// Check if running as root
	err := operatingsystem.CheckIfRunningAsRoot()
	if err != nil {
		return err
	}

	// Determine OS and install pre-requisites
	osType, err := operatingsystem.DetectOS()
	if err != nil {
		return err
	}
	ConfirmInstall, err := operatingsystem.PromptInstallPrereqs()
	if err != nil {
		return err
	}

	if ConfirmInstall {
		err = operatingsystem.InstallPrereqs(osType)
	} else if !ConfirmInstall {
		log.Fatal("Exited Workbench Installer")
	}
	if err != nil {
		return err
	}
	// Determine if we should disable the local firewall, then disable it
	// TODO: Add support for Ubuntu ufw
	firewalldEnabled, err := operatingsystem.CheckFirewallStatus(osType)
	if err != nil {
		return err
	}

	if firewalldEnabled {
		disableFirewall, err := operatingsystem.FirewallPrompt()
		if err != nil {
			return err
		}

		if disableFirewall {
			err = operatingsystem.DisableFirewall(osType)
			if err != nil {
				return err
			}
		}
	}

	// Determine Linux security status for the OS, then disable it
	// TODO: Add support for Ubuntu AppArmor
	selinuxEnabled, err := operatingsystem.CheckLinuxSecurityStatus(osType)
	if err != nil {
		return err
	}

	if selinuxEnabled {
		disableSELinux, err := operatingsystem.LinuxSecurityPrompt(osType)
		if err != nil {
			return err
		}

		if disableSELinux {
			err = operatingsystem.DisableLinuxSecurity()
			if err != nil {
				return err
			}
		}
	}
	// Languages
	selectedLanguages, err := languages.PromptAndRespond()
	if err != nil {
		return fmt.Errorf("issue selecting languages: %w", err)
	}

	// R
	err = languages.ScanAndHandleRVersions(osType)
	if err != nil {
		return fmt.Errorf("issue scanning, prompting or installing R: %w", err)
	}
	// Python
	if lo.Contains(selectedLanguages, "python") {
		err := languages.ScanAndHandlePythonVersions(osType)
		if err != nil {
			return fmt.Errorf("issue scanning, prompting or installing Python: %w", err)
		}
	}

	// Workbench
	err = workbench.CheckPromptDownloadAndInstallWorkbench(osType)
	if err != nil {
		return fmt.Errorf("issue checking, prompting, downloading or installing Workbench: %w", err)
	}

	// Licensing
	err = license.CheckPromptAndActivateLicense()
	if err != nil {
		return fmt.Errorf("issue checking, prompting or activating license: %w", err)
	}

	// Jupyter
	err = jupyter.ScanPromptInstallAndConfigJupyter()
	if err != nil {
		return fmt.Errorf("issue scanning, prompting, installing or configuring Jupyter: %w", err)
	}

	// Pro Drivers
	err = prodrivers.CheckPromptDownloadAndInstallProDrivers(osType)
	if err != nil {
		return fmt.Errorf("issue checking, prompting, downloading or installing Pro Drivers: %w", err)
	}

	// SSL
	sslChoice, err := ssl.PromptSSL()
	if err != nil {
		return fmt.Errorf("issue selecting if SSL is to be used: %w", err)
	}
	if sslChoice {
		err = ssl.PromptVerifyAndConfigSSL()
		if err != nil {
			return fmt.Errorf("issue verifying and configuring SSL: %w", err)
		}
	}

	// Authentication
	authChoice, err := authentication.PromptAuth()
	if err != nil {
		return fmt.Errorf("issue selecting if Authentication is to be setup: %w", err)
	}
	if authChoice {
		err = authentication.PromptAndConfigAuth(osType)
		if err != nil {
			return fmt.Errorf("issue prompting and configuring Authentication: %w", err)
		}
	}

	// Package Manager URL
	packageManagerChoice, err := packagemanager.PromptPackageManagerChoice()
	if err != nil {
		return fmt.Errorf("issue in prompt for Posit Package Manager choice: %w", err)
	}
	if packageManagerChoice {
		err = packagemanager.InteractivePackageManagerPrompts(osType)
		if err != nil {
			return fmt.Errorf("issue in interactive Posit Package Manager repo verification steps: %w", err)
		}
	} else {
		publicPackageManagerChoice, err := packagemanager.PromptPublicPackageManagerChoice()
		if err != nil {
			return fmt.Errorf("issue in prompt for Posit Public Package Manager choice: %w", err)
		}
		if publicPackageManagerChoice {
			err = packagemanager.VerifyAndBuildPublicPackageManager(osType)
			if err != nil {
				return fmt.Errorf("issue in verifying and building Public Posit Package Manager URL and repo: %w", err)
			}
		}
	}

	// Connect URL
	connectChoice, err := connect.PromptConnectChoice()
	if err != nil {
		return fmt.Errorf("issue in prompt for Connect URL choice: %w", err)
	}
	if connectChoice {
		err = connect.PromptVerifyAndConfigConnect()
		if err != nil {
			return fmt.Errorf("issue in prompting, verifying and saving Connect URL: %w", err)
		}
	}

	fmt.Println("\n Restarting RStudio Server and Launcher...")

	err = workbench.RestartRStudioServerAndLauncher()
	if err != nil {
		return fmt.Errorf("issue restarting RStudio Server and Launcher: %w", err)
	}

	fmt.Println("\n Printing the status of RStudio Server and Launcher...")

	err = workbench.StatusRStudioServerAndLauncher()
	if err != nil {
		return fmt.Errorf("issue running status for RStudio Server and Launcher: %w", err)
	}

	verifyChoice, err := workbench.PromptInstallVerify()
	if err != nil {
		return fmt.Errorf("issue selecting if verification is to be run: %w", err)
	}
	if verifyChoice {
		err = workbench.VerifyInstallation()
		if err != nil {
			return fmt.Errorf("issue running verification: %w", err)
		}
	}

	fmt.Println("\nThanks for using wbi! Please remember to make any needed manual configuration changes and restart RStudio Server and Launcher.")
	return nil
}

func setSetupOpts(setupOpts *setupOpts) {

}

func (opts *setupOpts) Validate() error {
	return nil
}

func newSetupCmd() *setupCmd {
	root := &setupCmd{opts: setupOpts{}}

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "setup",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setSetupOpts(&root.opts)
			if err := root.opts.Validate(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("setup-opts")
			if err := newSetup(root.opts); err != nil {
				return err
			}
			return nil
		},
	}
	root.cmd = cmd
	return root
}
