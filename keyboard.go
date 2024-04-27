package main

import (
	"log"
	"syscall"

	"github.com/rajveermalviya/go-wayland/wayland/client"
)

type Keyboard struct {
	Proxy   *client.Keyboard
	OnFocus func(kb *Keyboard, serial uint32)
	Data    interface{}
}

func NewKeyboard(proxy *client.Keyboard) *Keyboard {
	kb := &Keyboard{Proxy: proxy}
	kb.init()
	return kb
}

func (kb *Keyboard) init() {
	kb.Proxy.SetKeymapHandler(kb.keymapHandler)
	kb.Proxy.SetEnterHandler(kb.enterHandler)
	kb.Proxy.SetLeaveHandler(kb.leaveHandler)
	kb.Proxy.SetKeyHandler(kb.keyHandler)
	kb.Proxy.SetModifiersHandler(kb.modifiersHandler)
}

func (kb *Keyboard) keymapHandler(event client.KeyboardKeymapEvent) {
	// Close the file descriptor immediately to prevent resource leaks
	if err := syscall.Close(int(event.Fd)); err != nil {
		log.Printf("Failed to close file descriptor: %v", err)
	}
}

func (kb *Keyboard) enterHandler(event client.KeyboardEnterEvent) {
	if kb.OnFocus != nil {
		kb.OnFocus(kb, event.Serial)
	}
}

func (kb *Keyboard) leaveHandler(event client.KeyboardLeaveEvent) {
	// Optionally handle leave events
}

func (kb *Keyboard) keyHandler(event client.KeyboardKeyEvent) {
	// Handle key press or release
}

func (kb *Keyboard) modifiersHandler(event client.KeyboardModifiersEvent) {
	// Handle modifier keys
}
