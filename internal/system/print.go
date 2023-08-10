package system

import (
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
)

func PrintAndLogInfo(message string) {
	pterm.DefaultBasicText.Println(message)
	log.Info(message)
}
