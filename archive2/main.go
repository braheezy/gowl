package main

import (
	"log"

	"github.com/rajveermalviya/go-wayland/wayland/client"
	"golang.org/x/sys/unix"
)

// AppState holds the application state
type appState struct {
	display     *client.Display
	registry    *client.Registry
	seat        *client.Seat
	seatVersion uint32
	exit        bool
	keyboard    *client.Keyboard

	compositor    *client.Compositor
	deviceManager *client.DataDeviceManager
}

// HandleDisplayError handles client.Display errors
func (*appState) HandleDisplayError(e client.DisplayErrorEvent) {
	// Just log.Fatal for now
	log.Fatalf("display error event: %v", e)
}
func (app *appState) displayRoundTrip() {
	// Get display sync callback
	callback, err := app.display.Sync()
	if err != nil {
		log.Fatalf("unable to get sync callback: %v", err)
	}
	defer func() {
		if err2 := callback.Destroy(); err2 != nil {
			log.Println("unable to destroy callback:", err2)
		}
	}()

	done := false
	callback.SetDoneHandler(func(_ client.CallbackDoneEvent) {
		done = true
	})

	// Wait for callback to return
	for !done {
		app.display.Context().Dispatch()
	}
}

func main() {
	// Connect to wayland server
	display, err := client.Connect("")
	if err != nil {
		log.Fatalf("unable to connect to wayland server: %v", err)
	}
	app := &appState{display: display}
	display.SetErrorHandler(app.HandleDisplayError)

	// Get global interfaces registry
	app.registry, err = app.display.GetRegistry()
	if err != nil {
		log.Fatalf("unable to get global registry object: %v", err)
	}

	// Add global interfaces registrar handler
	app.registry.SetGlobalHandler(app.HandleRegistryGlobal)
	// Wait for interfaces to register
	app.displayRoundTrip()
	// Wait for handler events
	app.displayRoundTrip()

	// Start the dispatch loop
	if app.seat != nil {
		log.Println("setting up clipboard")
		app.setupClipboard("Hello, clipboard!")
	}

	log.Println("closing")
	app.cleanup()
}

func (app *appState) setupClipboard(text string) {
	dataDevice, err := app.deviceManager.GetDataDevice(app.seat)
	if err != nil {
		log.Fatalf("Failed to get data device: %v", err)
	}

	// Create a data source and offer text/plain
	dataSource, err := app.deviceManager.CreateDataSource()
	if err != nil {
		log.Fatalf("Failed to create data source: %v", err)
	}
	dataSource.Offer("text/plain")

	// Set the data when requested
	dataSource.SetSendHandler(func(event client.DataSourceSendEvent) {
		_, err := unix.Write(int(event.Fd), []byte(text))
		if err != nil {
			log.Printf("Failed to write to clipboard: %v", err)
		}
		unix.Close(int(event.Fd))
	})

	// Set the clipboard to our data source
	if err := dataDevice.SetSelection(dataSource, 0); err != nil { // Using '0' for the serial might not be valid in all contexts
		log.Fatalf("Failed to set selection: %v", err)
	}

	app.displayRoundTrip() // Ensure the request is processed before exiting
	app.displayRoundTrip()
}

func (app *appState) HandleRegistryGlobal(e client.RegistryGlobalEvent) {
	log.Printf("discovered an interface: %q\n", e.Interface)

	switch e.Interface {
	case "wl_compositor":
		compositor := client.NewCompositor(app.display.Context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, compositor)
		if err != nil {
			log.Fatalf("unable to bind wl_compositor interface: %v", err)
		}
		app.compositor = compositor

	case "wl_seat":
		seat := client.NewSeat(app.display.Context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, seat)
		if err != nil {
			log.Fatalf("unable to bind wl_seat interface: %v", err)
		}
		app.seat = seat
		app.seatVersion = e.Version
		// Add Keyboard & Pointer handlers
		seat.SetCapabilitiesHandler(app.HandleSeatCapabilities)
		seat.SetNameHandler(app.HandleSeatName)

	case "wl_data_device_manager":
		dataDeviceManager := client.NewDataDeviceManager(app.display.Context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, dataDeviceManager)
		if err != nil {
			log.Fatalf("unable to bind data device manager interface: %v", err)
		}
		app.deviceManager = dataDeviceManager
	}
}
func (*appState) HandleSeatName(e client.SeatNameEvent) {
	log.Printf("seat name: %v", e.Name)
}
func (app *appState) HandleSeatCapabilities(e client.SeatCapabilitiesEvent) {

	haveKeyboard := (e.Capabilities * uint32(client.SeatCapabilityKeyboard)) != 0

	if haveKeyboard && app.keyboard == nil {
		app.attachKeyboard()
	} else if !haveKeyboard && app.keyboard != nil {
		app.releaseKeyboard()
	}
}

func (app *appState) attachKeyboard() {
	keyboard, err := app.seat.GetKeyboard()
	if err != nil {
		log.Fatal("unable to register keyboard interface")
	}
	app.keyboard = keyboard

	keyboard.SetKeyHandler(app.HandleKeyboardKey)
	keyboard.SetKeymapHandler(app.HandleKeyboardKeymap)

	log.Println("keyboard interface registered")
}

func (app *appState) releaseKeyboard() {
	if err := app.keyboard.Release(); err != nil {
		log.Println("unable to release keyboard interface")
	}
	app.keyboard = nil

	log.Println("keyboard interface released")
}

func (app *appState) HandleKeyboardKey(e client.KeyboardKeyEvent) {
	// close on "esc"
	if e.Key == 1 {
		app.exit = true
	}
	log.Printf("key event: %v", e)
}
func (app *appState) HandleKeyboardKeymap(e client.KeyboardKeymapEvent) {
	defer unix.Close(e.Fd)
}

func (app *appState) cleanup() {
	// Release the keyboard if registered
	if app.keyboard != nil {
		app.releaseKeyboard()
	}

	// Release wl_seat handlers
	if app.seat != nil {
		if err := app.seat.Release(); err != nil {
			log.Println("unable to destroy wl_seat:", err)
		}
		app.seat = nil
	}

	if app.compositor != nil {
		if err := app.compositor.Destroy(); err != nil {
			log.Println("unable to destroy wl_compositor:", err)
		}
		app.compositor = nil
	}

	if app.registry != nil {
		if err := app.registry.Destroy(); err != nil {
			log.Println("unable to destroy wl_registry:", err)
		}
		app.registry = nil
	}

	if app.display != nil {
		if err := app.display.Destroy(); err != nil {
			log.Println("unable to destroy wl_display:", err)
		}
	}

	if app.deviceManager != nil {
		if err := app.deviceManager.Destroy(); err != nil {
			log.Println("unable to destroy wl_data_device_manager:", err)
		}
	}

	// Close the wayland server connection
	if err := app.display.Context().Close(); err != nil {
		log.Println("unable to close wayland context:", err)
	}
}
