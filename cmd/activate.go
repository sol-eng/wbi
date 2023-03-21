package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/sol-eng/wbi/internal/license"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type activateCmd struct {
	cmd  *cobra.Command
	opts activateOpts
}

type activateOpts struct {
	key string
}

func newActivate(activateOpts activateOpts) error {
	err := license.DetectAndActivateLicense(activateOpts.key)
	if err != nil {
		return fmt.Errorf("issue activating license: %w", err)
	}
	return nil
}

func setActivateOpts(activateOpts *activateOpts) {
	activateOpts.key = viper.GetString("key")
}

func (opts *activateOpts) Validate(args []string) error {
	// check args lengths
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided, please provide one argument")
	} else if len(args) > 1 {
		return fmt.Errorf("too many arguments provided, please provide only one argument")
	}

	// only the key flag is supported for license
	if opts.key != "" && args[0] != "license" {
		return fmt.Errorf("the key flag is only supported for license")
	}
	// the key flag is required for license
	if opts.key == "" && args[0] == "license" {
		return fmt.Errorf("the key flag is required for license")
	}

	return nil
}

func newActivateCmd() *activateCmd {
	var activateOpts activateOpts

	root := &activateCmd{opts: activateOpts}

	cmd := &cobra.Command{
		Use:   "activate license --key [license key]",
		Short: "Activate Workbench with a license key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setActivateOpts(&root.opts)
			if err := root.opts.Validate(args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("activate-opts")
			if err := newActivate(root.opts); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringP("key", "k", "", "License key to activate Workbench with")
	viper.BindPFlag("key", cmd.Flags().Lookup("key"))

	root.cmd = cmd
	return root
}
