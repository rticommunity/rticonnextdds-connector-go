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
)

type Replier struct {
	connector       *rti.Connector
	output          *rti.Output
	input           *rti.Input
	relatedIdentity string
}

func main() {
	// Find the file path to the XML configuration
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("runtime.Caller error")
	}
	filepath := path.Join(path.Dir(filename), "../RequestReplyExample.xml")

	var replier Replier
	var err error

	// Create a connector defined in the XML configuration
	replier.connector, err = rti.NewConnector("MyParticipantLibrary::Replier", filepath)
	if err != nil {
		log.Panic(err)
	}
	// Delete the connector when this main function returns
	defer replier.connector.Delete()

	// Get an output from the connector
	replier.output, err = replier.connector.GetOutput("ReplyPublisher::ReplyWriter")
	if err != nil {
		log.Panic(err)
	}
	// Get an input from the connector
	replier.input, err = replier.connector.GetInput("ReplySubscriber::ReplyReader")
	if err != nil {
		log.Panic(err)
	}

	// Set values to the instance and write the instance
	for i := 0; i < 500; i++ {
		var requestData types.Shape
		err = replier.receiveRequest(&requestData)

		if err != nil {
			log.Println(err)
		} else {
			log.Println("---Received Request---")
			log.Printf("color: %s\n", requestData.Color)
			log.Printf("x: %d\n", requestData.X)
			log.Printf("y: %d\n", requestData.Y)
			log.Printf("shapesize: %d\n", requestData.Shapesize)
		}

		var replyData types.Shape
		replyData.X = requestData.Y
		replyData.Y = requestData.X
		replyData.Shapesize = requestData.Shapesize
		replyData.Color = requestData.Color

		replier.sendReply(&replyData)

		log.Println("---Sent Reply---")
		log.Printf("color: %s\n", requestData.Color)
		log.Printf("x: %d\n", requestData.X)
		log.Printf("y: %d\n", requestData.Y)
		log.Printf("shapesize: %d\n", requestData.Shapesize)

		time.Sleep(time.Second * 1)
	}
}

func (replier *Replier) sendReply(data *types.Shape) {
	replier.output.Instance.Set(data)
	params := "{\"related_sample_identity\":" + replier.relatedIdentity + "}"
	log.Println(params)
	replier.output.WriteWithParams(params)
	log.Println("Replied...")
}

func (replier *Replier) receiveRequest(data *types.Shape) error {
	err := replier.connector.Wait(-1)
	if err == nil {
		replier.input.Take()
		valid, _ := replier.input.Infos.IsValid(0)
		if valid {
			err := replier.input.Samples.Get(0, data)
			if err != nil {
				log.Println(err)
			}
			log.Println(replier.input.Infos.GetIdentityJSON(0))
			identity, _ := replier.input.Infos.GetIdentityJSON(0)
			replier.relatedIdentity = identity
		}
	}
	return err
}
