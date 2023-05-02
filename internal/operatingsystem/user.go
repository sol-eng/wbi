package operatingsystem

import (
	"fmt"
	"strings"

	"github.com/sol-eng/wbi/internal/system"
)

func PromptAndVerifyUser() (string, bool, error) {
	var overallSkip bool
	var userAccount string
	for {
		username, err := PromptUserAccount()
		userAccount = username
		if err != nil {
			return "", false, err
		}
		if strings.Contains(userAccount, "skip") {
			overallSkip = true
			break
		}
		// lookup user account details
		user, err := UserLookup(userAccount)
		if err != nil {
			system.PrintAndLogInfo(fmt.Sprintf(`The user account "%s" you entered cannot be found. Please try again. To skip this section type "skip".`, userAccount))
		} else if user.Uid == "0" {
			system.PrintAndLogInfo(fmt.Sprintf(`The user account "%s" is root. A non-root user is required. Please try again. To skip this section type "skip".`, userAccount))
		} else if user.HomeDir == "" {
			system.PrintAndLogInfo(fmt.Sprintf(`The user account "%s" does not have a home directory. A home directory is required. Please try again. To skip this section type "skip".`, userAccount))
		} else {
			system.PrintAndLogInfo(fmt.Sprintf("user %s account found and validated", userAccount))
			break
		}
	}

	return userAccount, overallSkip, nil
}
