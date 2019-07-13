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
	"github.com/kyoungho/rticonnextdds-connector-go"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
)

func main() {
	// Find the file path to the XML configuration
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("runtime.Caller error")
	}
	filepath := path.Join(path.Dir(filename), "../ShapeExample.xml")

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

	ch := make(chan *rti.Samples, 100)
	input.ChannelSubscribe(ch)

	run := true

	// Get values from a received sample and print them
	for run == true {
		select {
		case sig := <-sigchan:
			log.Printf("Received signal %v: terminating\n", sig)
			run = false
		case samples := <-ch:
			numOfSamples := samples.GetLength()
			for i := 0; i < numOfSamples; i++ {
				//if infos.IsValid(i) {
				json, _ := samples.GetJSON(i)
				log.Printf("---Received Sample---\n, %s", json)
				//}
			}
		}
	}
}
