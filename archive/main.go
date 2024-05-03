package main

import (
	"log/syslog"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/log"

	"github.com/rajveermalviya/go-wayland/wayland/client"
	xdg_shell "github.com/rajveermalviya/go-wayland/wayland/stable/xdg-shell"
	xdg_activation "github.com/rajveermalviya/go-wayland/wayland/staging/xdg-activation-v1"
	primary_selection "github.com/rajveermalviya/go-wayland/wayland/unstable/primary-selection-v1"
	"golang.org/x/sys/unix"
)

type Config struct {
	HasXdgShell            bool
	HasWpPrimarySelection  bool
	HasGtkPrimarySelection bool
	HasWlrDataControl      bool
	HasGtkShell            bool
	HasXdgActivation       bool
}

type Options struct {
	primary bool
}

type App struct {
	registry      *client.Registry
	compositor    *client.Compositor
	display       *client.Display
	shm           *client.Shm
	shell         *client.Shell
	seats         []*client.Seat
	seat          *client.Seat
	config        *Config
	xdgActivation *xdg_activation.Activation
	xdgWmBase     *xdg_shell.WmBase
	gtkShell1     *GtkShell1
	// clipboardDevice                  ClipboardDevice
	options       *Options
	deviceManager GenericDeviceManager
	device        GenericDevice
	// supportsSelection                bool
	dataDeviceManager                *client.DataDeviceManager
	gtkPrimarySelectionDeviceManager *GtkPrimarySelectionDeviceManager
	zwpPrimarySelectionDeviceManager *primary_selection.PrimarySelectionDeviceManager
	zwlrDataControlManager           *ZwlrDataControlManagerV1
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
		log.Printf("%s\n", message)
		os.Exit(1)
	}

	// Maybe there is a tty we could write to?
	if tty, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0666); err == nil {
		defer tty.Close()
		log.New(tty).Printf("%s\n", message)
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
		// TODO: complain_about_wayland_connection()
	}

	checkClosedStdio()

	app := App{
		display: display,
		options: &Options{},
		// clipboardDevice: &ClipboardDevice{},
		config: &Config{},
	}
	app.options.primary = false
	registry, err := display.GetRegistry()
	if err != nil {
		log.Fatalf("Failed to get registry: %v", err)
	}
	app.registry = registry
	app.registry.SetGlobalHandler(app.HandleRegistryGlobal)

	// Wait for the initial set of globals to appear
	app.displayRoundTrip()

	seat := app.findRegistrySeat()

	deviceManager := app.findDeviceManager()
	if deviceManager == nil {
		// TODO: complain_about_selection_support
		log.Info("No suitable device manager found for the requested selection mode.")
	}
	device, err := deviceManager.GetDataDevice(seat)
	if err != nil {
		log.Fatalf("Failed to get data device: %v", err)
	}
	device.SelectionCallback()
}

func (app *App) findRegistrySeat() *client.Seat {
	app.displayRoundTrip()
	//TODO: If option `--seat` is supported, this is where it would be found and set

	return app.seat
}

func (app *App) findDeviceManager() GenericDeviceManager {
	/* For regular selection, we just look at the two supported
	 * protocols. We prefer wlr-data-control, as it doesn't require
	 * us to use the popup surface hack.
	 */
	if !app.options.primary {
		if app.config.HasWlrDataControl {
			return &ZwlrDataControlDeviceManager{app.zwlrDataControlManager}
		}
		return &DataDeviceManager{app.dataDeviceManager}
	} else {
		/* For primary selection, it's a bit more complicated. We also
		 * prefer wlr-data-control, but we don't know in advance whether
		 * the compositor supports primary selection, as unlike with
		 * other protocols here, the mere presence of wlr-data-control
		 * does not imply primary selection support. However, we assume
		 * that if a compositor supports primary selection at all, then
		 * if it supports wlr-data-control v2 it also supports primary
		 * selection over wlr-data-control; which is only reasonable.
		 */
		// if app.config.HasWlrDataControl && app.zwlrDataControlManager.Version >= 2 {
		// 	return app.zwlrDataControlManager
		// }
		// if app.config.HasWpPrimarySelection {
		// 	return app.zwpPrimarySelectionDeviceManager
		// }
		// if app.config.HasGtkPrimarySelection {
		// 	return app.gtkPrimarySelectionDeviceManager
		// }
		return nil
	}
}

func (app *App) displayRoundTrip() {
	// Get display sync callback
	callback, err := app.display.Sync()
	if err != nil {
		log.Fatalf("unable to get sync callback: %v", err)
	}
	defer func() {
		if err2 := callback.Destroy(); err2 != nil {
			log.Printf("unable to destroy callback: %v\n", err2)
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

func HandleSeatName(e client.SeatNameEvent) {
	log.Printf("seat name: %v", e.Name)
}

func (app *App) HandleRegistryGlobal(event client.RegistryGlobalEvent) {
	log.Debugf("discovered an interface: %q", event.Interface)
	switch event.Interface {
	case "wl_compositor":
		if event.Version > 2 {
			compositor := client.NewCompositor(app.context())
			err := app.registry.Bind(event.Name, event.Interface, event.Version, compositor)
			if err != nil {
				log.Fatalf("unable to bind wl_compositor interface: %v", err)
			}
			app.compositor = compositor
		}
	case "wl_shm":
		if event.Version >= 1 {
			shm := client.NewShm(app.context())
			err := app.registry.Bind(event.Name, event.Interface, event.Version, shm)
			if err != nil {
				log.Fatalf("unable to bind wl_shm interface: %v", err)
			}
			app.shm = shm
		}
	case "wl_shell":
		if event.Version > 1 {
			shell := client.NewShell(app.context())
			err := app.registry.Bind(event.Name, event.Interface, event.Version, shell)
			if err != nil {
				log.Fatalf("unable to bind wl_shell interface: %v", err)
			}
			app.shell = shell
		}
	case "wl_data_device_manager":
		if event.Version > 1 {
			deviceManager := client.NewDataDeviceManager(app.context())
			err := app.registry.Bind(event.Name, event.Interface, event.Version, deviceManager)
			if err != nil {
				log.Fatalf("unable to bind wl_data_device_manager interface: %v", err)
			}
			app.dataDeviceManager = deviceManager
		}
	case "wl_seat":
		if event.Version >= 2 {
			seat := client.NewSeat(app.context())
			err := app.registry.Bind(event.Name, event.Interface, event.Version, seat)
			if err != nil {
				log.Fatalf("unable to bind wl_seat interface: %v", err)
			}
			seat.SetNameHandler(HandleSeatName)
			app.seat = seat
			app.seats = append(app.seats, seat)
		}

	case "zwp_primary_selection_device_manager_v1":
		if event.Version >= 1 {
			zwpPrimarySelectionDeviceManager := primary_selection.NewPrimarySelectionDeviceManager(app.context())
			err := app.registry.Bind(event.Name, event.Interface, event.Version, zwpPrimarySelectionDeviceManager)
			if err != nil {
				log.Fatalf("unable to bind zwp_primary_selection_device_manager_v1 interface: %v", err)
			}
			app.zwpPrimarySelectionDeviceManager = zwpPrimarySelectionDeviceManager
			app.config.HasWpPrimarySelection = true
		}
	case "xdg_activation_v1":
		if event.Version >= 1 {
			xdgActivation := xdg_activation.NewActivation(app.context())
			err := app.registry.Bind(event.Name, event.Interface, event.Version, xdgActivation)
			if err != nil {
				log.Fatalf("unable to bind xdg_activation_v1 interface: %v", err)
			}
			app.xdgActivation = xdgActivation
			app.config.HasXdgActivation = true
		}
	case "gtk_primary_selection_device_manager":
		if event.Version >= 1 {
			gtkPrimarySelectionDeviceManager := NewGtkPrimarySelectionDeviceManager(app.context())
			err := app.registry.Bind(event.Name, event.Interface, event.Version, gtkPrimarySelectionDeviceManager)
			if err != nil {
				log.Fatalf("unable to bind gtk_primary_selection_device_manager interface: %v", err)
			}
			app.gtkPrimarySelectionDeviceManager = gtkPrimarySelectionDeviceManager
			app.config.HasGtkPrimarySelection = true
		}
	case "zwlr_data_control_manager_v1":

		zwlrDataControlManager := NewZwlrDataControlManagerV1(app.context())
		err := app.registry.Bind(event.Name, event.Interface, event.Version, zwlrDataControlManager)
		if err != nil {
			log.Fatalf("unable to bind zwlr_data_control_manager_v1 interface: %v", err)
		}
		zwlrDataControlManager.Version = event.Version
		app.zwlrDataControlManager = zwlrDataControlManager
		app.config.HasWlrDataControl = true
	case "gtk_shell1":
		if event.Version > 4 {
			gtkShell1 := NewGtkShell1(app.context())
			err := app.registry.Bind(event.Name, event.Interface, event.Version, gtkShell1)
			if err != nil {
				log.Fatalf("unable to bind gtk_shell1 interface: %v", err)
			}
			app.gtkShell1 = gtkShell1
			app.config.HasGtkShell = true
		}
	case "xdg_wm_base":
		app.config.HasXdgShell = true

	}
}
