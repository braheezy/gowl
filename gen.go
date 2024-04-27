package main

//go:generate go run github.com/rajveermalviya/go-wayland/cmd/go-wayland-scanner -pkg main -o gtk_primary_selection.go -i https://raw.githubusercontent.com/bugaevc/wl-clipboard/master/src/protocol/gtk-primary-selection.xml

//go:generate go run github.com/rajveermalviya/go-wayland/cmd/go-wayland-scanner -pkg main -o wlr_data_control.go -i https://raw.githubusercontent.com/bugaevc/wl-clipboard/master/src/protocol/wlr-data-control-unstable-v1.xml
