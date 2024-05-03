package main

//go:generate go-wayland-scanner -pkg main -o gtk_primary_selection.go -i https://raw.githubusercontent.com/bugaevc/wl-clipboard/master/src/protocol/gtk-primary-selection.xml

//go:generate go-wayland-scanner -pkg main -o wlr_data_control.go -i https://raw.githubusercontent.com/bugaevc/wl-clipboard/master/src/protocol/wlr-data-control-unstable-v1.xml

//go:generate go-wayland-scanner -pkg main -o gtk_shell.go -i https://raw.githubusercontent.com/bugaevc/wl-clipboard/master/src/protocol/gtk-shell.xml

//go:generate go-wayland-scanner -pkg main -o wl_roots_data_control.go -i https://gitlab.freedesktop.org/wlroots/wlr-protocols/-/raw/master/unstable/wlr-data-control-unstable-v1.xml?ref_type=heads
