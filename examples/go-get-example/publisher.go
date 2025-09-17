package main

import (
	"fmt"
	"log"

	rti "github.com/rticommunity/rticonnextdds-connector-go"
)

func main() {
	// Simple XML configuration inline (using str:// prefix for XML string)
	xmlConfig := `str://"<dds>
  <qos_library name="QosLibrary">
    <qos_profile name="DefaultProfile" base_name="BuiltinQosLibExp::Generic.StrictReliable" is_default_qos="true"/>
  </qos_library>
  
  <types>
    <struct name="TestType">
      <member name="message" type="string"/>
      <member name="count" type="long"/>
    </struct>
  </types>
  
  <domain_library name="MyDomainLibrary">
    <domain name="MyDomain" domain_id="0">
      <register_type name="TestType" type_ref="TestType"/>
      <topic name="TestTopic" register_type_ref="TestType"/>
    </domain>
  </domain_library>
  
  <domain_participant_library name="MyParticipantLibrary">
    <domain_participant name="Zero" domain_ref="MyDomainLibrary::MyDomain">
      <publisher name="MyPublisher">
        <data_writer name="MyWriter" topic_ref="TestTopic"/>
      </publisher>
    </domain_participant>
  </domain_participant_library>
</dds>"`

	fmt.Println("Creating RTI Connector...")

	// Create connector from XML string
	connector, err := rti.NewConnector("MyParticipantLibrary::Zero", xmlConfig)
	if err != nil {
		log.Fatal("Failed to create connector:", err)
	}
	defer connector.Delete()

	fmt.Println("✅ RTI Connector created successfully!")

	// Get output (writer) and publish a simple message
	output, err := connector.GetOutput("MyPublisher::MyWriter")
	if err != nil {
		log.Fatal("Failed to get output:", err)
	}

	// Publish one test message
	output.Instance.SetString("message", "Hello from Go get user!")
	output.Instance.SetInt("count", 42)
	output.Write()

	fmt.Println("✅ Successfully published test message!")
	fmt.Println("RTI Connector Go is working with libraries downloaded via go get workflow!")
}
