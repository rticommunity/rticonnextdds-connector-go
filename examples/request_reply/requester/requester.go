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
	"github.com/rticommunity/rticonnextdds-connector-go/types"
	"log"
	"path"
	"runtime"
	"time"
	//"strings"
	"strconv"
)

type Requester struct {
	connector *rti.Connector
	output    *rti.Output
	input     *rti.Input
	identity  string
}

func main() {
	// Find the file path to the XML configuration
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("runtime.Caller error")
	}
	filepath := path.Join(path.Dir(filename), "../RequestReplyExample.xml")

	var requester Requester
	var err error

	// Create a connector defined in the XML configuration
	requester.connector, err = rti.NewConnector("MyParticipantLibrary::Requester", filepath)
	if err != nil {
		log.Panic(err)
	}
	// Delete the connector when this main function returns
	defer requester.connector.Delete()

	// Get an output from the connector
	requester.output, err = requester.connector.GetOutput("RequestPublisher::RequestWriter")
	if err != nil {
		log.Panic(err)
	}
	// Get an input from the connector
	requester.input, err = requester.connector.GetInput("RequestSubscriber::RequestReader")
	if err != nil {
		log.Panic(err)
	}

	// Set values to the instance and write the instance
	for i := 0; i < 500; i++ {
		var requestData types.Shape
		requestData.X = i
		requestData.Y = i * 2
		requestData.Shapesize = 30
		requestData.Color = "BLUE"

		requester.identity = "{\"writer_guid\":[1,1,68,56,24,72,236,10,57,15,60,144,128,0,0,1],\"sequence_number\":" + strconv.Itoa(i+1) + "}"

		requester.sendRequest(&requestData)

		log.Println("---Sent Request---")
		log.Printf("color: %s\n", requestData.Color)
		log.Printf("x: %d\n", requestData.X)
		log.Printf("y: %d\n", requestData.Y)
		log.Printf("shapesize: %d\n", requestData.Shapesize)

		var replyData types.Shape
		err = requester.receiveReply(&replyData)

		if err != nil {
			log.Println(err)
		} else {
			log.Println("---Received Reply---")
			log.Printf("color: %s\n", replyData.Color)
			log.Printf("x: %d\n", replyData.X)
			log.Printf("y: %d\n", replyData.Y)
			log.Printf("shapesize: %d\n", replyData.Shapesize)
		}

		time.Sleep(time.Second * 1)
	}
}

func (requester *Requester) sendRequest(data *types.Shape) {
	requester.output.Instance.Set(data)
	params := "{\"identity\":" + requester.identity + "}"
	requester.output.WriteWithParams(params)
	log.Println("Requested...")
}

func (requester *Requester) receiveReply(data *types.Shape) error {
	err := requester.connector.Wait(10000)
	if err == nil {
		requester.input.Take()
		valid, _ := requester.input.Infos.IsValid(0)
		if valid {
			related_identity, _ := requester.input.Infos.GetRelatedIdentityJSON(0)
			log.Println(requester.identity)
			log.Println(related_identity)

			//if strings.Compare(requester.identity, related_identity) == 0 {
			if requester.identity == related_identity {
				err := requester.input.Samples.Get(0, data)
				if err != nil {
					log.Println(err)
				}
			} else {
				log.Println("Ignored an unrelated sample")
			}
		}
	}
	return err
}
