package cmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sol-eng/wbi/internal/languages"
	"github.com/sol-eng/wbi/internal/system"
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
		system.PrintAndLogInfo(strings.Join(rVersions, "\n"))
	} else if language == "python" {
		pythonVersions, err := languages.ScanForPythonVersions()
		if err != nil {
			return fmt.Errorf("issue occured in scanning for Python versions: %w", err)
		}
		system.PrintAndLogInfo(strings.Join(pythonVersions, "\n"))
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

	// adding two spaces to have consistent formatting
	exampleText := []string{
		"To scan for existing R and Python installations:",
		"  wbi scan r",
		"  wbi scan python",
	}

	cmd := &cobra.Command{
		Use:     "scan [lanaguage]",
		Short:   "Scan for installed versions of R or Python",
		Example: strings.Join(exampleText, "\n"),
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
