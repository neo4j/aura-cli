//go:build windows

package clicfg

import (
	"log"

	"golang.org/x/sys/windows/registry"
)

func init() {
	p, err := registry.ExpandString("%LOCALAPPDATA%")

	if err != nil {
		log.Panic(err)
	}

	ConfigPrefix = p
}
