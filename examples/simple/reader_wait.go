/*****************************************************************************
*    (c) 2005-2015 Copyright, Real-Time Innovations, All rights reserved.    *
*                                                                            *
*  RTI grants Licensee a license to use, modify, compile, and create         *
*  derivative works of the Software.  Licensee has the right to distribute   *
*  object form only for use with RTI products. The Software is provided      *
*  "as is", with no warranty of any type, including any warranty for fitness *
*  for any purpose. RTI is under no obligation to maintain or support the    *
*  Software.  RTI shall not be liable for any incidental or consequential    *
*  damages arising out of the use or inability to use the software.          *
*                                                                            *
*****************************************************************************/

package main

import (
	"github.com/kyoungho/rticonnextdds-connector"
	"path"
	"runtime"
	"log"
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
		numOfSamples := input.Samples.GetLength()
		for j := 0; j < numOfSamples; j++ {
			if input.Infos.IsValid(j) {
				json, err := input.Samples.GetJson(j)
				if (err != nil) {
					log.Println(err)
				} else {
					log.Println(string(json))
				}
			}
		}
	}
}
