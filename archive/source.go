package main

import "github.com/rajveermalviya/go-wayland/wayland/client"

type GenericSource interface {
}

type WlDataSource struct {
	proxy *client.DataSource
}
