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
	"math/rand"
	"path"
	"runtime"
	"time"
)

type Shape struct {
	Color     string `json:"color"`
	X         []int  `json:"x"`
	Y         []int  `json:"y"`
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

	// Get an output from the connector
	output, err := connector.GetOutput("MyPublisher::MySquareWriter")
	if err != nil {
		log.Panic(err)
	}

	// Set values to the instance and write the instance
	for i := 0; i < 500; i++ {
		var shape Shape
		for j := 0; j < rand.Intn(10); j++ {
			shape.X = append(shape.X, i)
			shape.Y = append(shape.Y, i*2)
		}
		shape.Shapesize = 30
		shape.Color = "BLUE"
		output.Instance.Set(&shape)
		output.Write()

		output.ClearMembers()
		log.Println("Writing...")
		time.Sleep(time.Second * 1)
	}
}
