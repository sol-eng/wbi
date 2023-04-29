package cmd

import (
	"fmt"
	"strings"

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
	"github.com/sol-eng/wbi/internal/system"
	"github.com/sol-eng/wbi/internal/workbench"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type setupCmd struct {
	cmd  *cobra.Command
	opts setupOpts
}

type setupOpts struct {
	step string
}

func newSetup(setupOpts setupOpts) error {

	// define step either "" if no flag or a step if flag is set
	step := setupOpts.step
	if step == "" {
		step = "start"
	}

	if step == "start" {
		system.PrintAndLogInfo("Welcome to the Workbench Installer!")
		step = "prereqs"
	}

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

	if step == "prereqs" {
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
		step = "firewall"
	}

	if step == "firewall" {
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
		step = "security"
	}

	if step == "security" {
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
		step = "languages"
	}

	selectedLanguages := []string{"r", "python"}
	if step == "languages" {
		// Languages
		selectedLanguages, err = languages.PromptAndRespond()
		if err != nil {
			return fmt.Errorf("issue selecting languages: %w", err)
		}
		step = "r"
	}

	if step == "r" {
		// R
		err = languages.ScanAndHandleRVersions(osType)
		if err != nil {
			return fmt.Errorf("issue scanning, prompting or installing R: %w", err)
		}
		step = "python"
	}

	if step == "python" {
		// Python
		if lo.Contains(selectedLanguages, "python") {
			err := languages.ScanAndHandlePythonVersions(osType)
			if err != nil {
				return fmt.Errorf("issue scanning, prompting or installing Python: %w", err)
			}
		}
		step = "workbench"
	}

	if step == "workbench" {
		// Workbench
		err = workbench.CheckPromptDownloadAndInstallWorkbench(osType)
		if err != nil {
			return fmt.Errorf("issue checking, prompting, downloading or installing Workbench: %w", err)
		}
		step = "license"
	}

	if step == "license" {
		// Licensing
		err = license.CheckPromptAndActivateLicense()
		if err != nil {
			return fmt.Errorf("issue checking, prompting or activating license: %w", err)
		}
		step = "jupyter"
	}

	if step == "jupyter" {
		// Jupyter
		err = jupyter.ScanPromptInstallAndConfigJupyter()
		if err != nil {
			return fmt.Errorf("issue scanning, prompting, installing or configuring Jupyter: %w", err)
		}
		step = "prodrivers"
	}

	if step == "prodrivers" {
		// Pro Drivers
		err = prodrivers.CheckPromptDownloadAndInstallProDrivers(osType)
		if err != nil {
			return fmt.Errorf("issue checking, prompting, downloading or installing Pro Drivers: %w", err)
		}
		step = "ssl"
	}

	if step == "ssl" {
		// SSL
		sslChoice, err := ssl.PromptSSL()
		if err != nil {
			return fmt.Errorf("issue selecting if SSL is to be used: %w", err)
		}
		if sslChoice {
			certPath, keyPath, err := ssl.PromptAndVerifySSL()
			if err != nil {
				return fmt.Errorf("issue verifying and configuring SSL: %w", err)
			}
			serverURL, err := ssl.PromptServerURL()
			if err != nil {
				return fmt.Errorf("issue prompting for server URL: %w", err)
			}
			workbench.WriteSSLConfig(certPath, keyPath, serverURL)
			if err != nil {
				return fmt.Errorf("error writing ssl configuration to file rserver.conf: %w", err)
			}
		}
		step = "auth"
	}

	if step == "auth" {
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
		step = "packagemanager"
	}

	if step == "packagemanager" {
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
		step = "connect"
	}

	if step == "connect" {
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
		step = "restart"
	}

	if step == "restart" {
		system.PrintAndLogInfo("\nRestarting RStudio Server and Launcher...")

		err = workbench.RestartRStudioServerAndLauncher()
		if err != nil {
			return fmt.Errorf("issue restarting RStudio Server and Launcher: %w", err)
		}
		step = "status"
	}

	if step == "status" {
		system.PrintAndLogInfo("\nPrinting the status of RStudio Server and Launcher...")

		err = workbench.StatusRStudioServerAndLauncher()
		if err != nil {
			return fmt.Errorf("issue running status for RStudio Server and Launcher: %w", err)
		}
		step = "verify"
	}

	if step == "verify" {
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
		step = "done"
	}

	system.PrintAndLogInfo("\nThanks for using wbi!")
	return nil
}

func setSetupOpts(setupOpts *setupOpts) {
	setupOpts.step = viper.GetString("step")
}

func (opts *setupOpts) Validate(args []string) error {
	// ensure no args are passed
	if len(args) > 0 {
		return fmt.Errorf("no arguments are supported for this command")
	}
	// ensure step is valid
	validSteps := []string{"start", "prereqs", "firewall", "security", "languages", "r", "python", "workbench", "license", "jupyter", "prodrivers", "ssl", "auth", "packagemanager", "connect", "restart", "status", "verify"}
	if opts.step != "" && !lo.Contains(validSteps, opts.step) {
		return fmt.Errorf("invalid step: %s", opts.step)
	}

	return nil
}

func newSetupCmd() *setupCmd {
	root := &setupCmd{opts: setupOpts{}}

	// adding two spaces to have consistent formatting
	exampleText := []string{
		"To start an interactive setup process for Workbench:",
		"  wbi setup",
		"",
		"To start an interactive setup process for Workbench at a certain step:",
		"  wbi setup --step [STEP]",
	}

	cmd := &cobra.Command{
		Use:     "setup",
		Short:   "Launch an interactive setup process for Workbench",
		Example: strings.Join(exampleText, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setSetupOpts(&root.opts)
			if err := root.opts.Validate(args); err != nil {
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
		SilenceUsage: true,
	}

	stepHelp := `The step to start at. Valid steps are: start, prereqs, firewall, security, languages, r, python, workbench, license, jupyter, prodrivers, ssl, auth, packagemanager, connect, restart, status, verify.`

	cmd.Flags().StringP("step", "s", "", stepHelp)
	viper.BindPFlag("step", cmd.Flags().Lookup("step"))

	root.cmd = cmd
	return root
}
