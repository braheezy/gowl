package main

import "github.com/rajveermalviya/go-wayland/wayland/client"

type GenericShellSurface interface{}

type WlShellSurface struct {
	proxy *client.ShellSurface
}
