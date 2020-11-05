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
	"path"
	"runtime"

	rti "github.com/rticommunity/rticonnextdds-connector-go"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("runtime.Caller error")
	}
	filepath := path.Join(path.Dir(filename), "../ShapeExample.xml")
	connector, err := rti.NewConnector("MyParticipantLibrary::Zero", filepath)
	if err != nil {
		log.Panic(err)
	}
	defer connector.Delete()
	input, err := connector.GetInput("MySubscriber::MySquareReader")
	if err != nil {
		log.Panic(err)
	}

	for {
		connector.Wait(-1)
		input.Take()
		numOfSamples, _ := input.Samples.GetLength()
		for j := 0; j < numOfSamples; j++ {
			valid, _ := input.Infos.IsValid(j)
			if valid {
				json, err := input.Samples.GetJSON(j)
				if err != nil {
					log.Println(err)
				} else {
					log.Println(string(json))
				}
			}
		}
	}
}
