package main

import (
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rajveermalviya/go-wayland/wayland/client"
	"golang.org/x/sys/unix"
)

type App struct {
	registry      *client.Registry
	compositor    *client.Compositor
	display       *client.Display
	shm           *client.Shm
	shell         *client.Shell
	deviceManager *client.DataDeviceManager
	seat          *client.Seat
}

func checkClosedStdio() {
	message := "gowl-clipboard has been launched with a closed standard file descriptor. This is a bug in the software that has launched gowl-clipboard. Aborting."

	valid := true
	for fd := 0; fd <= 2; fd++ {
		if _, err := unix.FcntlInt(uintptr(fd), syscall.F_GETFD, 0); err != nil {
			valid = false
			break
		}
	}
	if valid {
		return
	}

	// Check if stderr is directed to a character device
	stat, err := os.Stderr.Stat()
	if err == nil && (stat.Mode()&os.ModeCharDevice) != 0 {
		log.Println(message)
		os.Exit(1)
	}

	// Maybe there is a tty we could write to?
	if tty, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0666); err == nil {
		defer tty.Close()
		log.New(tty, "", 0).Println(message)
		os.Exit(1)
	}

	// As a last resort, try syslog
	if logger, err := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, "wl-clipboard"); err == nil {
		defer logger.Close()
		logger.Err(message)
		os.Exit(1)
	}

	// If syslog fails, we just log to stderr (which might not be visible)
	log.Fatalf(message)
}

func main() {
	/* Ignore SIGPIPE.
	 * We don't really output anything
	 * to our stdout, yet we don't want
	 * to get killed when writing clipboard
	 * contents to a closed pipe.
	 */
	signal.Ignore(syscall.SIGPIPE)

	display, err := client.Connect("")
	if err != nil {
		log.Fatalf("Failed to connect to Wayland display: %v", err)
	}

	checkClosedStdio()

	registry, err := display.GetRegistry()
	if err != nil {
		log.Fatalf("Failed to get registry: %v", err)
	}
	app := App{
		registry: registry,
		display:  display,
	}

	// Setup global handler
	app.registry.SetGlobalHandler(app.HandleRegistryGlobal)

	app.displayRoundTrip()

	println("all interfaces registered")
}

func (app *App) displayRoundTrip() {
	// Get display sync callback
	callback, err := app.display.Sync()
	if err != nil {
		log.Fatalf("unable to get sync callback: %v", err)
	}
	defer func() {
		if err2 := callback.Destroy(); err2 != nil {
			log.Println("unable to destroy callback: ", err2)
		}
	}()

	done := false
	callback.SetDoneHandler(func(_ client.CallbackDoneEvent) {
		done = true
	})

	// Wait for callback to return
	for !done {
		app.dispatch()
	}
}

func (app *App) context() *client.Context {
	return app.display.Context()
}

func (app *App) dispatch() {
	app.display.Context().Dispatch()
}

func (app *App) HandleRegistryGlobal(event client.RegistryGlobalEvent) {
	switch event.Interface {
	case "wl_compositor":
		compositor := client.NewCompositor(app.context())
		err := app.registry.Bind(event.Name, event.Interface, event.Version, compositor)
		if err != nil {
			log.Fatalf("unable to bind wl_compositor interface: %v", err)
		}
		app.compositor = compositor
	case "wl_shm":
		shm := client.NewShm(app.context())
		err := app.registry.Bind(event.Name, event.Interface, event.Version, shm)
		if err != nil {
			log.Fatalf("unable to bind wl_shm interface: %v", err)
		}
		app.shm = shm
	case "wl_shell":
		shell := client.NewShell(app.context())
		err := app.registry.Bind(event.Name, event.Interface, event.Version, shell)
		if err != nil {
			log.Fatalf("unable to bind wl_shell interface: %v", err)
		}
		app.shell = shell

	case "wl_data_device_manager":
		deviceManager := client.NewDataDeviceManager(app.context())
		err := app.registry.Bind(event.Name, event.Interface, event.Version, deviceManager)
		if err != nil {
			log.Fatalf("unable to bind wl_data_device_manager interface: %v", err)
		}
		app.deviceManager = deviceManager
	case "wl_seat":
		seat := client.NewSeat(app.context())
		err := app.registry.Bind(event.Name, event.Interface, event.Version, seat)
		if err != nil {
			log.Fatalf("unable to bind wl_seat interface: %v", err)
		}
		app.seat = seat

	}
}
