//go:build windows

package clicfg

import (
	"golang.org/x/sys/windows/registry"
)

func init() {
	p, err := registry.ExpandString("%LOCALAPPDATA%")

	if err != nil {
		panic(err)
	}

	ConfigPrefix = p
}
