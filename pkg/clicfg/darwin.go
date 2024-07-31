//go:build darwin

package clicfg

import (
	"os/user"
	"path/filepath"
)

func init() {
	currentUser, _ := user.Current()
	homeDir := currentUser.HomeDir

	ConfigPrefix = filepath.Join(homeDir, "Library/Preferences")
}
