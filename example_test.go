package rti

import (
	"fmt"
	"github.com/rticommunity/rticonnextdds-connector-go/types"
	"path"
	"runtime"
)

func ExampleNewConnector() {
	// Get the path for test XML configs
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/Test.xml")

	participantProfile := "MyParticipantLibrary::Zero"
	connector, err := NewConnector(participantProfile, xmlPath)
	if err != nil {
		// Handle an error
	}
	connector.Delete()
	// Output:
}

func ExampleConnector_GetInput() {
	// Create Connector
	connector := newTestConnector()
	defer connector.Delete()

	input, err := connector.GetInput("MySubscriber::MyReader")
	if input == nil || err != nil {
		// Handle an error
	}
	// Output:
}

func ExampleConnector_GetOutput() {
	// Create Connector
	connector := newTestConnector()
	defer connector.Delete()

	output, err := connector.GetInput("MyPublisher::MyWriter")
	if output == nil || err != nil {
		// Handle an error
	}
	// Output:
}

func ExampleInfos_GetIdentity() {
	// Create Connector
	connector := newTestConnector()
	defer connector.Delete()

	input := newTestInput(connector)
	output := newTestOutput(connector)

	var outputTestData types.Test
	outputTestData.St = "test"
	output.Instance.Set(&outputTestData)
	output.Write()

	connector.Wait(-1)
	input.Take()
	writerId := input.Infos.GetIdentity(0)
	//fmt.Printf("wrtier_guid: %x\n", writerId.WriterGuid)
	fmt.Printf("seuqnece_number: %d\n", writerId.SequenceNumber)
	// Output:
	//seuqnece_number: 1
}

/*
func ExampleInput_AsyncSubscribe() {
	connector := newTestConnector()
	defer connector.Delete()
	input := newTestInput(connector)
	output := newTestOutput(connector)

	var outputTestData types.Test
	outputTestData.St = "test"
	output.Instance.Set(&outputTestData)
	output.Write()

	input.AsyncSubscribe(func(samples *Samples, infos *Infos) {
		numOfSamples := samples.GetLength()
		for i := 0; i < numOfSamples; i++ {
			if infos.IsValid(i) {
				json, _ := samples.GetJSON(i)
				fmt.Printf("---Received Sample---\n%s", json)
				return
			}
		}
	})

	// Output:
	//  ---Received Sample---
	//{
	//"st":"test",
	//"b":false,
	//"c":0,
	//"s":0,
	//"us":0,
	//"l":0,
	//"ul":0,
	//"ll":0,
	//"ull":0,
	//"f":0,
	//"d":0
	//}
}

func ExampleInput_ChannelSubscribe() {
	connector := newTestConnector()
	defer connector.Delete()
	input := newTestInput(connector)
	output := newTestOutput(connector)

	ch := make(chan *Samples)
	input.ChannelSubscribe(ch)

	var outputTestData types.Test
	outputTestData.St = "test"
	output.Instance.Set(&outputTestData)
	output.Write()

	for {
		samples := <-ch
		numOfSamples := samples.GetLength()
		for i := 0; i < numOfSamples; i++ {
			json, _ := samples.GetJSON(i)
			fmt.Printf("---Received Sample---\n%s", json)

			// ---
			// This return should not be in your actual code,
			// but for example test, we need to break out of
			// infinite for loop
			return
			// ---
		}
	}
	// Output:
	//  ---Received Sample---
	//{
	//"st":"test",
	//"b":false,
	//"c":0,
	//"s":0,
	//"us":0,
	//"l":0,
	//"ul":0,
	//"ll":0,
	//"ull":0,
	//"f":0,
	//"d":0
	//}
}
*/
