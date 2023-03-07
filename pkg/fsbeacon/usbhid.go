// Package fsbeacon implements a control interface for using Farming Simulator USB beacon light toys.
// It requires that host USB HID access is available.
// fsbeacon cannot differentiate between multiple beacons plugged into the same device due to a hardware limitation.
package fsbeacon

import (
	hid "github.com/sstallion/go-hid"
	"sync"
	"time"
)

// Full credit to https://gist.github.com/steve228uk/873d653f1ecec0456ea3f475b6e54f68 for the main information on this
// Did a little more digging on my own with an OpenVizla USB MITM device, but not really much to be seen -
// the game just repeatedly sends the same command every few seconds to ensure the beacon stays awake.

const (
	fs22vid uint16 = 0x340d
	fs22pid uint16 = 0x1710
)

var fs22MagicFlash []byte = []byte{0x0, 0xFF, 0x07, 0xFF, 0x50, 0xFF, 0x1C, 0x8E, 0xB0, 0xB8}
var fs22MagicSpin []byte = []byte{0x0, 0xFF, 0x01, 0x66, 0xC8, 0xFF, 0xAD, 0x52, 0x81, 0xD6}
var fs22MagicStop []byte = []byte{0x0, 0xFF, 0x00, 0x00, 0x64, 0x00, 0x32, 0x9E, 0xD7, 0x0D}

// Beacon is any device that can do beacon-like activities.
type Beacon interface {
	// Flash causes the device to strobe or flash.
	Flash() error
	// Spin causes the device to produce or animate a spinning effect.
	Spin() error
	// Stop causes the device to return to an idle, off state.
	Stop() error
	// Close ends the session with the given device. The device should not be used again after calling this function.
	Close() error
}

// FarmBeacon represents a Farm Simulator USB beacon device.
type FarmBeacon struct {
	sync.Mutex
	stop      chan bool
	hidDevice *hid.Device
}

// OpenFarmBeacon opens the connected Farm Simulator USB beacon for use, and returns a handler struct.
func OpenFarmBeacon() (beacon *FarmBeacon, err error) {
	// Initialize the hid package.
	if err = hid.Init(); err != nil {
		return
	}

	// Open the beacon.
	var bcn FarmBeacon
	bcn.hidDevice, err = hid.OpenFirst(fs22vid, fs22pid)
	if err != nil {
		return
	}
	return &bcn, nil
}

func (beacon *FarmBeacon) worker(magic []byte) {
	beacon.Mutex.Lock()
	defer beacon.Mutex.Unlock()

	beacon.stop = make(chan bool)
	defer close(beacon.stop)

	ticker := time.NewTicker(5 * time.Second)
	// Set the initial device state
	beacon.hidDevice.Write(magic)
	for {
		select {
		case <-ticker.C:
			// Update HID state with desired magic
			beacon.hidDevice.Write(magic)
		case updateHID := <-beacon.stop:
			if updateHID {
				beacon.hidDevice.Write(fs22MagicStop)
			}
			return
		}
	}
}

// Close destroys the connection to the USB beacon.
func (beacon *FarmBeacon) Close() error {
	return beacon.hidDevice.Close()
}

// Flash causes the beacon to flash indefinitely.
func (beacon *FarmBeacon) Flash() (err error) {
	if err = beacon.ready(); err != nil {
		return err
	}
	if beacon.stop != nil {
		// Beacon worker already running, stop it quietly.
		beacon.stop <- false
	}
	// Set the initial state so that if an error occurs, we know about it.
	_, err = beacon.hidDevice.Write(fs22MagicFlash)
	if err != nil {
		return
	}
	go beacon.worker(fs22MagicFlash)
	return
}

// Spin causes the beacon to display a spinning effect indefinitely.
func (beacon *FarmBeacon) Spin() (err error) {
	if err = beacon.ready(); err != nil {
		return err
	}
	if beacon.stop != nil {
		// Beacon worker already running, stop it quietly.
		beacon.stop <- false
	}
	// Set the initial state so that if an error occurs, we know about it.
	_, err = beacon.hidDevice.Write(fs22MagicSpin)
	if err != nil {
		return
	}
	go beacon.worker(fs22MagicSpin)
	return
}

// Stop causes the beacon to turn off.
func (beacon *FarmBeacon) Stop() (err error) {
	if err = beacon.ready(); err != nil {
		return err
	}
	if beacon.stop != nil {
		// Beacon worker already running, stop it along with the animation.
		beacon.stop <- true
	}
	return
}

// ready is a helper function that determines whether the beacon is connected and ready for further commands.
func (beacon *FarmBeacon) ready() (err error) {
	if err = beacon.hidDevice.Error(); err.Error() != "Success" {
		return beacon.hidDevice.Error()
	}
	return nil
}
