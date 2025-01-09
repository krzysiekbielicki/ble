package dev

import (
	"github.com/krzysiekbielicki/ble"
	"github.com/krzysiekbielicki/ble/linux"
)

// DefaultDevice ...
func DefaultDevice() (d ble.Device, err error) {
	return linux.NewDevice()
}
