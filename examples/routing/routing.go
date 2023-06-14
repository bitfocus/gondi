package main

import (
	"fmt"
	"time"

	"github.com/bitfocus/gondi"
)

func main() {
	fmt.Println("Initializing NDI")
	gondi.InitLibrary("")

	version := gondi.GetVersion()
	fmt.Printf("NDI version: %s\n", version)

	findInstance, err := gondi.NewFindInstance(true, "", "10.20.10.42")
	if err != nil {
		panic(err)
	}
	defer findInstance.Destroy()

	// Wait for sources to appear
	fmt.Println("Looking for sources...")
	for {
		more := findInstance.WaitForSources(5000)
		if !more {
			break
		}
	}

	// Fetch the sources
	sources := findInstance.GetCurrentSources()

	if len(sources) == 0 {
		fmt.Println("No sources found, cannot continue")
		return
	}

	selectedSource := sources[0]
	fmt.Printf("Source selected: %s\n", selectedSource.Name())

	// Set up a NDI Output called "route1"
	route1, err := gondi.NewRoutingInstance("Output 1", "")
	if err != nil {
		panic(err)
	}
	defer route1.Destroy()

	// Route it to show the content of selectedSource
	route1.Change(selectedSource)

	time.Sleep(60 * time.Second)
}
