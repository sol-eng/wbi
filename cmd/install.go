package cmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sol-eng/wbi/internal/jupyter"
	"github.com/sol-eng/wbi/internal/languages"
	"github.com/sol-eng/wbi/internal/os"
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
	versions []string
	path     string
}

func newInstall(installOpts installOpts, program string) error {
	// Determine OS
	osType, err := os.DetectOS()
	if err != nil {
		return err
	}

	if program == "r" {
		if len(installOpts.versions) == 0 {
			_, err = languages.ScanAndHandleRVersions(osType)
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
		}
	} else if program == "python" {
		if len(installOpts.versions) == 0 {
			_, err = languages.ScanAndHandlePythonVersions(osType)
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
		}
	} else if program == "workbench" {
		workbenchInstalled := workbench.VerifyWorkbench()
		if !workbenchInstalled {
			err := workbench.DownloadAndInstallWorkbench(osType)
			if err != nil {
				return fmt.Errorf("issue installing Workbench: %w", err)
			}
		} else {
			return fmt.Errorf("workbench is already installed")
		}
	} else if program == "prodrivers" {
		proDriversExistingStatus, err := prodrivers.CheckExistingProDrivers()
		if err != nil {
			return fmt.Errorf("issue in checking for prior pro driver installation: %w", err)
		}
		if !proDriversExistingStatus {
			err := prodrivers.DownloadAndInstallProDrivers(osType)
			if err != nil {
				return fmt.Errorf("issue installing Pro Drivers: %w", err)
			}
		} else {
			return fmt.Errorf("Pro Drivers are already installed")
		}
	} else if program == "jupyter" {
		if installOpts.path != "" {
			err := jupyter.InstallJupyter(installOpts.path)
			if err != nil {
				return fmt.Errorf("issue installing Jupyter: %w", err)
			}
		} else {
			pythonVersions, err := languages.ScanForPythonVersions()
			if err != nil {
				return fmt.Errorf("issue occured in scanning for Python versions: %w", err)
			}
			jupyterPythonTarget, err := jupyter.KernelPrompt(pythonVersions)
			if err != nil {
				return fmt.Errorf("issue selecting Python location for Jupyter: %w", err)
			}
			err = jupyter.InstallJupyter(jupyterPythonTarget)
			if err != nil {
				return fmt.Errorf("issue installing Jupyter: %w", err)
			}
		}
	}
	return nil
}

func setInstallOpts(installOpts *installOpts) {
	installOpts.versions = viper.GetStringSlice("version")
	installOpts.path = viper.GetString("path")
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

	// ensure versions are valid if provided for r and python
	if args[0] == "r" && len(opts.versions) != 0 {
		err := languages.ValidateRVersions(opts.versions)
		if err != nil {
			return fmt.Errorf("invalid R versions: %w", err)
		}
	} else if args[0] == "python" && len(opts.versions) != 0 {
		err := languages.ValidatePythonVersions(opts.versions)
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

	cmd := &cobra.Command{
		Use:   "install",
		Short: "install",
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

	cmd.Flags().StringSliceP("version", "v", []string{}, "")
	viper.BindPFlag("version", cmd.Flags().Lookup("version"))

	cmd.Flags().StringP("path", "p", "", "")
	viper.BindPFlag("path", cmd.Flags().Lookup("path"))

	root.cmd = cmd
	return root
}
