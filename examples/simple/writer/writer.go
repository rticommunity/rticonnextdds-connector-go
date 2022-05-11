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
	"fmt"
	"log"
	"path"
	"runtime"

	rti "github.com/rticommunity/rticonnextdds-connector-go"
)

func main() {
	// Find the file path to the XML configuration
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("runtime.Caller error")
	}
	filepath := path.Join(path.Dir(filename), "../ShapeExample.xml")

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

	totalShapes := 0

	log.Print("Starting to write...")

	// Set values to the instance and write the instance
	for i := 0; true; i++ {

		if totalShapes%99 == 0 {
			i = 0
			log.Print(fmt.Sprintf("written %d shapes", totalShapes))
		}
		totalShapes++

		output.Instance.SetInt("x", i)
		output.Instance.SetInt("y", i*2)
		output.Instance.SetInt("shapesize", 30)
		output.Instance.SetString("color", "BLUE")
		err := output.Write()
		if err != nil {
			log.Print("error writing:", err)
		}
		// time.Sleep(time.Millisecond * 50)
	}
}
