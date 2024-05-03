package main

import (
	"log"
	"os"

	"github.com/neurlang/wayland/wl"
	"github.com/neurlang/wayland/wlclient"
)

type App struct {
	dataDeviceManager *wl.DataDeviceManager
	seat              *wl.Seat
	ctx               *wl.Context
	display           *wl.Display
}

func (app *App) HandleRegistryGlobal(event wl.RegistryGlobalEvent) {
	println("Interface found:", event.Interface)
	if event.Interface == "wl_data_device_manager" {
		app.dataDeviceManager = wl.NewDataDeviceManager(app.ctx)
	} else if event.Interface == "wl_seat" {
		app.seat = wl.NewSeat(app.ctx)
	}
}

func (app *App) HandleDataSourceSend(event wl.DataSourceSendEvent) {
	// Write the hardcoded string to the provided file descriptor
	file := os.NewFile(uintptr(event.Fd), "clipboard")
	defer file.Close()
	content := "Hello, Wayland clipboard!"
	if _, err := file.Write([]byte(content)); err != nil {
		log.Fatalf("Failed to write to clipboard: %v", err)
	}
}

func main() {
	display, err := wl.Connect("")
	if err != nil {
		log.Fatalf("Failed to connect to display: %v", err)
	}
	app := &App{
		display: display,
		ctx:     display.Context(),
	}

	registry, err := app.display.GetRegistry()
	if err != nil {
		log.Fatalf("Failed to get registry: %v", err)
	}

	// Listen for global objects and bind to the wl_data_device_manager interface
	registry.AddGlobalHandler(app)
	app.display.Context().Run() // Process events to ensure dataDeviceManager is bound

	wlclient.DisplayRoundtrip(display)

	if app.dataDeviceManager == nil {
		log.Fatalf("Data device manager not available")
	}

	// Create a data source and offer text data
	dataSource, err := app.dataDeviceManager.CreateDataSource()
	if err != nil {
		log.Fatalf("Failed to create data source: %v", err)
	}
	dataSource.AddSendHandler(app)
	dataSource.Offer("text/plain")

	// Set this data source as the clipboard content
	dataDevice, err := app.dataDeviceManager.GetDataDevice(app.seat)
	if err != nil {
		log.Fatalf("Failed to get data device: %v", err)
	}
	err = dataDevice.SetSelection(dataSource, 0)
	if err != nil {
		log.Fatalf("Failed to set selection: %v", err)
	}

	// Keep the application running to handle the send event
	display.Context().Run()
	wlclient.DisplayRoundtrip(display)

}
