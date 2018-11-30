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
	"github.com/tidwall/gjson"
	"github.com/urfave/cli"
	"log"
	"os"
	"time"
)

func main() {
	var configPath string
	var participantName string
	var writerName string
	var data string
	var count int
	var interval int

	app := cli.NewApp()
	app.Name = "DDS Introspection"
	app.Usage = "Injecting DDS data for diagnostics"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "cfgPath, c",
			Value:       "ShapeExample.xml",
			Usage:       "Path for XML App Creation configuration",
			Destination: &configPath,
		},
		cli.StringFlag{
			Name:        "participant, p",
			Value:       "MyParticipantLibrary::Zero",
			Usage:       "Participant name in XML",
			Destination: &participantName,
		},
		cli.StringFlag{
			Name:        "writer, w",
			Value:       "MyPublisher::MySquareWriter",
			Usage:       "Writer name in XML",
			Destination: &writerName,
		},
		cli.StringFlag{
			Name:        "data, d",
			Value:       `{"color": "BLUE", "x": 10, "y": 20, "shapesize": 30}`,
			Usage:       "Data in JSON format",
			Destination: &data,
		},
		cli.IntFlag{
			Name:        "count",
			Value:       10,
			Usage:       "Number of injecting samples",
			Destination: &count,
		},
		cli.IntFlag{
			Name:        "interval",
			Value:       1,
			Usage:       "Interval between samples in seconds",
			Destination: &interval,
		},
	}

	app.Action = func(c *cli.Context) error {
		// Create a connector defined in the XML configuration
		connector, err := rti.NewConnector(participantName, configPath)
		if err != nil {
			log.Panic(err)
		}
		// Delete the connector when this main function returns
		defer connector.Delete()

		// Get an output from the connector
		output, err := connector.GetOutput(writerName)
		if err != nil {
			log.Panic(err)
		}

		for i := 0; i < count; i++ {
			m, ok := gjson.Parse(data).Value().(map[string]interface{})
			if !ok {
				log.Panic(err)
			}

			output.Instance.Set(&m)
			output.Write()

			log.Println("Writing...")
			time.Sleep(time.Second * time.Duration(interval))
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Panic(err)
	}

}
