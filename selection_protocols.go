package main

import (
	client "github.com/rajveermalviya/go-wayland/wayland/unstable/primary-selection-v1"
)

func foo() {
	client.NewPrimarySelectionDevice(nil)
	NewGtkPrimarySelectionDevice(nil)
	NewZwlrDataControlDeviceV1(nil)
}
