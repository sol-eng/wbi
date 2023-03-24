package cmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sol-eng/wbi/internal/jupyter"
	"github.com/sol-eng/wbi/internal/languages"
	"github.com/sol-eng/wbi/internal/operatingsystem"

	"github.com/sol-eng/wbi/internal/prodrivers"
	"github.com/sol-eng/wbi/internal/system"
	"github.com/sol-eng/wbi/internal/workbench"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type installCmd struct {
	cmd  *cobra.Command
	opts installOpts
}

type installOpts struct {
	versions  []string
	path      string
	symlink   bool
	addToPATH bool
}

func newInstall(installOpts installOpts, program string) error {
	// Determine OS
	osType, err := operatingsystem.DetectOS()
	if err != nil {
		return err
	}

	if program == "r" {
		if len(installOpts.versions) == 0 {
			err = languages.ScanAndHandleRVersions(osType)
			if err != nil {
				return fmt.Errorf("ScanAndHandleRVersions: %w", err)
			}
		} else {
			for _, rVersion := range installOpts.versions {
				err = languages.DownloadAndInstallR(rVersion, osType)
				if err != nil {
					return fmt.Errorf("issue installing R versions: %w", err)
				}
			}
			if installOpts.symlink {
				fullRPath := "/opt/R/" + installOpts.versions[0] + "/bin/R"
				err = languages.CheckAndSetRSymlinks(fullRPath)
				if err != nil {
					return fmt.Errorf("issue setting R symlinks: %w", err)
				}
			}
		}
	} else if program == "python" {
		if len(installOpts.versions) == 0 {
			err = languages.ScanAndHandlePythonVersions(osType)
			if err != nil {
				return fmt.Errorf("ScanAndHandlePythonVersions: %w", err)
			}
		} else {
			for _, pythonVersion := range installOpts.versions {
				err = languages.DownloadAndInstallPython(pythonVersion, osType)
				if err != nil {
					return fmt.Errorf("issue installing Python versions: %w", err)
				}
			}
			if installOpts.addToPATH {
				// TODO add to PATH the latest version of Python (this just chooses the first version listed)
				fullPythonPath := "/opt/python/" + installOpts.versions[0] + "/bin"
				err = system.AddToPATH(fullPythonPath, "python")
				if err != nil {
					return fmt.Errorf("issue adding Python binary to PATH: %w", err)
				}
			}
		}
	} else if program == "workbench" {
		err := workbench.CheckDownloadAndInstallWorkbench(osType)
		if err != nil {
			return fmt.Errorf("issue installing Workbench: %w", err)
		}
	} else if program == "prodrivers" {
		err = prodrivers.CheckPromptDownloadAndInstallProDrivers(osType)
		if err != nil {
			return fmt.Errorf("issue checking, prompting, downloading or installing Pro Drivers: %w", err)
		}
	} else if program == "jupyter" {
		if installOpts.path != "" {
			err := jupyter.InstallAndConfigJupyter(installOpts.path)
			if err != nil {
				return fmt.Errorf("issue installing or configuring Jupyter: %w", err)
			}
		} else {
			err := jupyter.ScanPromptInstallAndConfigJupyter()
			if err != nil {
				return fmt.Errorf("issue scanning, prompting, installing or configuring Jupyter: %w", err)
			}
		}
	}
	return nil
}

func setInstallOpts(installOpts *installOpts) {
	installOpts.versions = viper.GetStringSlice("version")
	installOpts.path = viper.GetString("path")
	installOpts.symlink = viper.GetBool("symlink")
	installOpts.addToPATH = viper.GetBool("add-to-path")
}

func (opts *installOpts) Validate(args []string) error {
	// check args lengths
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided, please provide one argument")
	} else if len(args) > 1 {
		return fmt.Errorf("too many arguments provided, please provide only one argument")
	}

	// only the flag for path is supported for jupyter
	if opts.path != "" && args[0] != "jupyter" {
		return fmt.Errorf("the path flag is only supported for jupyter")
	}

	// only the flag for symlink is supported for r
	if opts.symlink && args[0] != "r" {
		return fmt.Errorf("the symlink flag is only supported for r")
	}

	// only the flag for add-to-path (addToPATH) is supported for python
	if opts.addToPATH && args[0] != "python" {
		return fmt.Errorf("the add-to-path flag is only supported for python")
	}

	// ensure versions are valid if provided for r and python
	if args[0] == "r" && len(opts.versions) != 0 {
		err := languages.ValidateRVersions(opts.versions)
		if err != nil {
			return fmt.Errorf("invalid R versions: %w", err)
		}
	} else if args[0] == "python" && len(opts.versions) != 0 {
		osType, err := operatingsystem.DetectOS()
		if err != nil {
			return fmt.Errorf("issue detecting OS: %w", err)
		}
		err = languages.ValidatePythonVersions(opts.versions, osType)
		if err != nil {
			return fmt.Errorf("invalid Python versions: %w", err)
		}
	}

	// ensure versions are not provided for workbench, prodrivers or jupyter
	if args[0] == "workbench" && len(opts.versions) != 0 {
		return fmt.Errorf("workbench does not support specifying versions")
	} else if args[0] == "prodrivers" && len(opts.versions) != 0 {
		return fmt.Errorf("prodrivers does not support specifying versions")
	} else if args[0] == "jupyter" && len(opts.versions) != 0 {
		return fmt.Errorf("jupyter does not support specifying versions")
	}

	// ensure path is valid if provided
	if opts.path != "" {
		pathExists := system.VerifyFileExists(opts.path)
		if !pathExists {
			return fmt.Errorf("the path provided does not exist")
		}
	}

	// ensure program is valid
	if args[0] != "r" && args[0] != "python" && args[0] != "workbench" && args[0] != "prodrivers" && args[0] != "jupyter" {
		return fmt.Errorf("invalid argument provided")
	}

	return nil
}

func newInstallCmd() *installCmd {
	var installOpts installOpts

	root := &installCmd{opts: installOpts}

	// adding two spaces to have consistent formatting
	exampleText := []string{
		"To start an interactive prompt to select R or Python version(s):",
		"  wbi install r",
		"  wbi install python",
		"",
		"To install a specific R or Python version:",
		"  wbi install r --version 4.2.2",
		"  wbi install python --version 3.11.2",
		"",
		"To install multiple R or Python versions:",
		"  wbi install r --version 4.2.2,4.1.3",
		"  wbi install python --version 3.11.2,3.10.10",
		"",
		"To install Workbench:",
		"  wbi install workbench",
		"",
		"To install Pro Drivers:",
		"  wbi install prodrivers",
		"",
		"To start an interactive prompt to select Jupyter install location:",
		"  wbi install jupyter",
		"",
		"To install Jupyter to a specific Python location:",
		"  wbi install jupyter --path /path/to/python",
	}

	cmd := &cobra.Command{
		Use:     "install [program]",
		Short:   "Install R, Python, Workbench, Pro Drivers, or Jupyter",
		Example: strings.Join(exampleText, "\n"),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setInstallOpts(&root.opts)
			if err := root.opts.Validate(args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("install-opts")
			if err := newInstall(root.opts, strings.ToLower(args[0])); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringSliceP("version", "v", []string{}, "Version(s) of R or Python to install. Multiple values can be passed by seperating each version with a comma.")
	viper.BindPFlag("version", cmd.Flags().Lookup("version"))

	cmd.Flags().StringP("path", "p", "", "Python location to install Jupyter to.")
	viper.BindPFlag("path", cmd.Flags().Lookup("path"))

	cmd.Flags().BoolP("symlink", "s", false, "Symlinks both R and Rscript for the first version of R specified to /usr/local/bin/.")
	viper.BindPFlag("symlink", cmd.Flags().Lookup("symlink"))

	cmd.Flags().BoolP("add-to-path", "a", false, "Adds the first Python version specified to users PATH by adding a file in /etc/profile.d/.")
	viper.BindPFlag("add-to-path", cmd.Flags().Lookup("add-to-path"))

	root.cmd = cmd
	return root
}
