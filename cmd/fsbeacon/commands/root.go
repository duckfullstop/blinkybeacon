package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "fsbeacon",
	Short: "fsbeacon is a simple control application for controlling Farming Simulator beacons.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runUntilInterrupted() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Print("Caught SIGTERM, exiting.")
}

func runWithTimeout(t time.Duration) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-c:
		fmt.Print("Caught SIGTERM, exiting.")
		return
	case <-time.After(t):
		return
	}
}
