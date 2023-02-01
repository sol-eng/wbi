package cmd

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dpastoor/wbi/internal/langscanner"
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

	//TODO: check if workbench installed

	var qs = []*survey.Question{
		{
			Name: "languages",
			Prompt: &survey.MultiSelect{
				Message: "What languages will you use",
				Options: []string{"R", "python"},
				Default: []string{"R", "python"},
			},
		},
	}
	answers := struct {
		Languages []string `survey:"languages"`
	}{}
	err := survey.Ask(qs, &answers)
	fmt.Println("You just chose the languages: ", answers.Languages)

	// TODO: conditionally scan for R versions
	rVersions, err := langscanner.ScanForRVersions()
	if err != nil {
		log.Fatal(err)
	}
	if len(rVersions) == 0 {
		// TODO: if no R version found, send link to R installation doc
		log.Fatal("no R versions found at locations: \n", strings.Join(langscanner.GetRootDirs(), "\n"))
	}

	fmt.Println("found R versions: ", strings.Join(rVersions, ", "))
	// TODO: scan for python versions

	// TODO: If python found -- setup jupyter or ask to setup jupyter or check

	// TODO: Handle SSL cert
	// * ask if want SSL
	// * if yes, ask for cert and key
	// * make sure cert and key are valid

	// TODO:

	return err
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
