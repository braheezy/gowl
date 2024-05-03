package main

import "github.com/rajveermalviya/go-wayland/wayland/client"

type GenericDeviceManager interface {
	GetDataDevice(seat *client.Seat) (GenericDevice, error)
}

type GenericDevice interface {
	SelectionCallback()
}

type DataDeviceManager struct {
	proxy *client.DataDeviceManager
}

func (d *DataDeviceManager) GetDataDevice(seat *client.Seat) (GenericDevice, error) {
	dev, err := d.proxy.GetDataDevice(seat)
	if err != nil {
		return DataDevice{}, err
	}
	return DataDevice{dev}, nil
}

type ZwlrDataControlDeviceManager struct {
	proxy *ZwlrDataControlManagerV1
}

func (d *ZwlrDataControlDeviceManager) GetDataDevice(seat *client.Seat) (GenericDevice, error) {
	dev, err := d.proxy.GetDataDevice(seat)
	if err != nil {
		return DataDevice{}, err
	}
	return ZwlrDataControlDevice{dev}, nil
}

type DataDevice struct {
	proxy *client.DataDevice
}

func (d DataDevice) SelectionCallback() {
	d.proxy.SetSelectionHandler(func(event client.DataDeviceSelectionEvent) {
		if event.Id != nil {
			event.Id.Destroy()
		}
	})
}

type ZwlrDataControlDevice struct {
	proxy *ZwlrDataControlDeviceV1
}

func (d ZwlrDataControlDevice) SelectionCallback() {
	d.proxy.SetSelectionHandler(func(event ZwlrDataControlDeviceV1SelectionEvent) {
		if event.Id != nil {
			event.Id.Destroy()
		}
	})
}
