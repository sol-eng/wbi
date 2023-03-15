package cmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sol-eng/wbi/internal/languages"
	"github.com/spf13/cobra"
)

type scanCmd struct {
	cmd  *cobra.Command
	opts scanOpts
}

type scanOpts struct {
}

func newScan(scanOpts scanOpts, language string) error {
	if language == "r" {
		rVersions, err := languages.ScanForRVersions()
		if err != nil {
			return fmt.Errorf("issue occured in scanning for R versions: %w", err)
		}
		fmt.Println(strings.Join(rVersions, ", "))
	} else if language == "python" {
		pythonVersions, err := languages.ScanForPythonVersions()
		if err != nil {
			return fmt.Errorf("issue occured in scanning for Python versions: %w", err)
		}
		fmt.Println(strings.Join(pythonVersions, ", "))
	} else {
		return fmt.Errorf("language %s is not supported", language)
	}

	return nil
}

func setScanOpts(scanOpts *scanOpts) {

}

func (opts *scanOpts) Validate(args []string) error {
	// check args lengths
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided, please provide one argument")
	} else if len(args) > 1 {
		return fmt.Errorf("too many arguments provided, please provide only one argument")
	}

	// ensure only r or python is provided
	if args[0] != "r" && args[0] != "python" {
		return fmt.Errorf("invalid language provided, please provide one of the following: r, python")
	}
	return nil
}

func newScanCmd() *scanCmd {
	var scanOpts scanOpts

	root := &scanCmd{opts: scanOpts}

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "scan",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setScanOpts(&root.opts)
			if err := root.opts.Validate(args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("scan-opts")
			if err := newScan(root.opts, strings.ToLower(args[0])); err != nil {
				return err
			}
			return nil
		},
	}

	root.cmd = cmd
	return root
}
