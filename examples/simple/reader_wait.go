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
