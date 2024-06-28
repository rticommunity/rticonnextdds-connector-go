/*****************************************************************************
*   (c) 2005-2015 Copyright, Real-Time Innovations.  All rights reserved.    *
*                                                                            *
* No duplications, whole or partial, manual or electronic, may be made       *
* without express written permission.  Any such copies, or revisions thereof,*
* must display this notice unaltered.                                        *
* This code contains trade secrets of Real-Time Innovations, Inc.            *
*                                                                            *
*****************************************************************************/

package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"

	rti "github.com/rticommunity/rticonnextdds-connector-go"
)

func main() {
	// Find the file path to the XML configuration
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("runtime.Caller error")
	}
	filepath := path.Join(path.Dir(filename), "../Bob.xml")

	// Create a channel to receive signals from OS
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Create a connector defined in the XML configuration
	connector, err := rti.NewConnector("MyParticipantLibrary::Zero", filepath)
	if err != nil {
		log.Panic(err)
	}
	// Delete the connector when this main function returns
	defer connector.Delete()

	// Get an input from the connector
	input, err := connector.GetInput("MySubscriber::MySquareReader")
	if err != nil {
		log.Panic(err)
	}

	run := true

	// Get values from a received sample and print them
	for run == true {
		select {
		case sig := <-sigchan:
			log.Printf("Received signal %v: terminating\n", sig)
			run = false
		default:
			input.Take()
			numOfSamples, _ := input.Samples.GetLength()
			for i := 0; i < numOfSamples; i++ {
				valid, _ := input.Infos.IsValid(i)
				if valid {
					color, _ := input.Samples.GetString(i, "color")
					x, _ := input.Samples.GetInt(i, "x")
					y, _ := input.Samples.GetInt(i, "y")
					shapesize, _ := input.Samples.GetInt(i, "shapesize")

					log.Println("---Received Sample---")
					log.Printf("color: %s\n", color)
					log.Printf("x: %d\n", x)
					log.Printf("y: %d\n", y)
					log.Printf("shapesize: %d\n", shapesize)
				}
			}
			time.Sleep(time.Second * 1)
		}
	}
}
