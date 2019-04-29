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
	"github.com/rticommunity/rticonnextdds-connector-go"
	"log"
	"path"
	"runtime"
	"time"
)

type Shape struct {
	Color     string `json:"color"`
	X         []int    `json:"x"`
	Y         []int    `json:"y"`
	Shapesize int    `json:"shapesize"`
}

func main() {
	// Find the file path to the XML configuration
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("runtime.Caller error")
	}
	filepath := path.Join(path.Dir(filename), "../ShapeSeqExample.xml")

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

	// Get values from a received sample and print them
	for i := 0; i < 500; i++ {
		input.Take()
		numOfSamples := input.Samples.GetLength()
		for j := 0; j < numOfSamples; j++ {
			if input.Infos.IsValid(j) {
				var shape Shape
				err := input.Samples.Get(j, &shape)
				if err != nil {
					log.Println(err)
				}
				log.Println("---Received Sample---")
				log.Printf("color: %s\n", shape.Color)
				log.Printf("x: %d\n", shape.X)
				log.Printf("y: %d\n", shape.Y)
				log.Printf("shapesize: %d\n", shape.Shapesize)
			}
		}
		time.Sleep(time.Second * 1)
	}
}
