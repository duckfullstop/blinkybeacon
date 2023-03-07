package fsbeacon

import (
	"testing"
	"time"
)

func TestFarmBeacon(t *testing.T) {
	beacon, err := OpenFarmBeacon()
	if err != nil {
		if err.Error() == "No HID devices with requested VID/PID found in the system." {
			t.Skip("Farming Simulator USB beacon not found, skipping test.")
		}
		t.Fatal(err)
	}
	defer beacon.Close()

	t.Log("Spinning...")
	err = beacon.Spin()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	t.Log("Flashing...")
	err = beacon.Flash()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	t.Log("Stopping!")
	err = beacon.Stop()
	if err != nil {
		t.Fatal(err)
	}
}
