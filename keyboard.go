package main

import (
	"log"

	"github.com/neurlang/wayland/wl"
)

func (app *appState) HandleKeyboardEnter(wl.KeyboardEnterEvent) {
	app.pointerEvent.eventMask &= ^keyboardEventLeave
	app.pointerEvent.eventMask |= keyboardEventEnter

	app.redecorate(true)

}

func (app *appState) HandleKeyboardLeave(wl.KeyboardLeaveEvent) {
	app.pointerEvent.eventMask &= ^keyboardEventEnter
	app.pointerEvent.eventMask |= keyboardEventLeave

	app.redecorate(false)

	app.pointerEvent.moveWindow = false

}
func (app *appState) attachKeyboard() {
	keyboard, err := app.seat.GetKeyboard()
	if err != nil {
		log.Fatal("unable to register keyboard interface")
	}
	app.keyboard = keyboard

	keyboard.AddKeyHandler(app)
	keyboard.AddEnterHandler(app)
	keyboard.AddLeaveHandler(app)

	log.Print("keyboard interface registered")
}
