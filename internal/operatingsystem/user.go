package operatingsystem

import (
	"fmt"

	"github.com/sol-eng/wbi/internal/system"
)

func PromptAndVerifyUser() (string, error) {
	userAccount, err := PromptUserAccount()
	if err != nil {
		return "", err
	}
	// lookup user account details
	user, err := UserLookup(userAccount)
	if err != nil {
		return "", fmt.Errorf("user %s account not found", userAccount)
	}
	// verify non-root and a home directory exists
	if user.Uid == "0" {
		return "", fmt.Errorf("user %s account is root. A non-root user is required", userAccount)
	}
	if user.HomeDir == "" {
		return "", fmt.Errorf("user %s account does not have a home directory. A home directory is required", userAccount)
	}
	system.PrintAndLogInfo(fmt.Sprintf("user %s account found and validated", userAccount))

	return user.Username, nil
}
