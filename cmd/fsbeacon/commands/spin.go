package commands

import (
	"fmt"
	"github.com/duckfullstop/blinkybeacon/pkg/fsbeacon"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func init() {
	rootCmd.AddCommand(spinCmd)
}

var spinCmd = &cobra.Command{
	Use:   "spin (seconds)",
	Short: "Spins the beacon for a set length of time. Defaults to 5 seconds.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  handleSpinBeacon,
}

func handleSpinBeacon(_ *cobra.Command, args []string) (err error) {
	var runtime time.Duration
	if len(args) > 0 {
		var arg float64
		arg, err = strconv.ParseFloat(args[0], 32)
		if err != nil {
			return
		}
		runtime = time.Duration(arg) * time.Second
	}

	var d fsbeacon.Beacon
	d, err = fsbeacon.OpenFarmBeacon()
	if err != nil {
		return
	}
	defer d.Close()

	// Start flashing the beacon - this starts a routine in the package to ensure the beacon keeps running indefinitely
	if runtime != 0 {
		fmt.Printf("Spinning beacon for %s.", runtime)
	} else {
		fmt.Printf("Spinning beacon - press ^C to stop.")
	}

	err = d.Spin()
	if err != nil {
		return err
	}

	// Wait for the configured runtime (the aforementioned goroutine is making sure the beacon is whipped)
	// or a SIGTERM from somewhere
	if runtime != 0 {
		runWithTimeout(runtime)
	} else {
		runUntilInterrupted()
	}

	// Now stop the beacon with the stop command. The connection gets closed by the earlier defer.
	err = d.Stop()
	return err
}
