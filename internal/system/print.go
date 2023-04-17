package system

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func PrintAndLogInfo(message string) {
	fmt.Println(message)
	log.Info(message)
}
