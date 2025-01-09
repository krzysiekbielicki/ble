package dev

import (
	"github.com/krzysiekbielicki/ble"
	"github.com/krzysiekbielicki/ble/darwin"
)

// DefaultDevice ...
func DefaultDevice() (d ble.Device, err error) {
	return darwin.NewDevice()
}
