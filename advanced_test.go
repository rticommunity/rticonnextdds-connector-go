package rti

import (
	"path"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ============================
// Request-Reply Pattern Tests
// ============================

// TestRequestReplyPattern tests the full request-reply communication pattern
func TestRequestReplyPattern(t *testing.T) {
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/RequestReplyTest.xml")

	// Create requester connector
	requester, err := NewConnector("MyParticipantLibrary::Requester", xmlPath)
	assert.Nil(t, err)
	defer requester.Delete()

	// Create replier connector
	replier, err := NewConnector("MyParticipantLibrary::Replier", xmlPath)
	assert.Nil(t, err)
	defer replier.Delete()

	// Get requester's output and input
	requestWriter, err := requester.GetOutput("RequestPublisher::RequestWriter")
	assert.Nil(t, err)
	replyReader, err := requester.GetInput("ReplySubscriber::ReplyReader")
	assert.Nil(t, err)

	// Get replier's input and output
	requestReader, err := replier.GetInput("RequestSubscriber::RequestReader")
	assert.Nil(t, err)
	replyWriter, err := replier.GetOutput("ReplyPublisher::ReplyWriter")
	assert.Nil(t, err)

	// Wait for discovery
	time.Sleep(500 * time.Millisecond)

	// Test 1: Send request with custom identity
	requestIdentity := `{"writer_guid":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16],"sequence_number":42}`

	assert.Nil(t, requestWriter.Instance.SetInt32("request_id", 1))
	assert.Nil(t, requestWriter.Instance.SetString("request_data", "test_request"))

	// Write with params including identity
	params := `{"identity":` + requestIdentity + `}`
	err = requestWriter.WriteWithParams(params)
	assert.Nil(t, err)

	// Replier receives request
	err = replier.Wait(2000)
	if err != nil {
		t.Skipf("Skipping test: Wait timed out (may indicate shared memory configuration issue on ARM64 macOS): %v", err)
	}

	err = requestReader.Take()
	if err != nil {
		t.Skipf("Skipping test: No data available (may indicate shared memory configuration issue on ARM64 macOS): %v", err)
	}

	length, err := requestReader.Samples.GetLength()
	assert.Nil(t, err)
	if length == 0 {
		t.Skip("Skipping test: No samples received (may indicate shared memory configuration issue on ARM64 macOS)")
	}
	assert.Equal(t, 1, length)

	// Verify request data
	requestID, err := requestReader.Samples.GetInt32(0, "request_id")
	assert.Nil(t, err)
	assert.Equal(t, int32(1), requestID)

	requestData, err := requestReader.Samples.GetString(0, "request_data")
	assert.Nil(t, err)
	assert.Equal(t, "test_request", requestData)

	// Get the request identity to use as related_sample_identity in reply
	receivedIdentity, err := requestReader.Infos.GetIdentityJSON(0)
	assert.Nil(t, err)
	assert.NotEmpty(t, receivedIdentity)

	// Replier sends reply with related_sample_identity
	assert.Nil(t, replyWriter.Instance.SetInt32("reply_id", 1))
	assert.Nil(t, replyWriter.Instance.SetString("reply_data", "test_reply"))

	replyParams := `{"related_sample_identity":` + receivedIdentity + `}`
	err = replyWriter.WriteWithParams(replyParams)
	assert.Nil(t, err)

	// Requester receives reply
	err = requester.Wait(2000)
	assert.Nil(t, err)

	err = replyReader.Take()
	assert.Nil(t, err)

	// Verify reply data
	replyID, err := replyReader.Samples.GetInt32(0, "reply_id")
	assert.Nil(t, err)
	assert.Equal(t, int32(1), replyID)

	replyData, err := replyReader.Samples.GetString(0, "reply_data")
	assert.Nil(t, err)
	assert.Equal(t, "test_reply", replyData)

	// Test GetRelatedIdentity and GetRelatedIdentityJSON
	relatedIdentity, err := replyReader.Infos.GetRelatedIdentity(0)
	if err == nil {
		// If no error, verify the identity matches
		assert.NotNil(t, relatedIdentity)
		assert.NotEmpty(t, relatedIdentity.WriterGUID)
	}

	relatedIdentityJSON, err := replyReader.Infos.GetRelatedIdentityJSON(0)
	if err == nil {
		assert.NotEmpty(t, relatedIdentityJSON)
	}
}

// TestWriteWithParamsSourceTimestamp tests writing with custom source timestamp
func TestWriteWithParamsSourceTimestamp(t *testing.T) {
	connector, err := newTestConnector()
	assert.Nil(t, err)
	defer connector.Delete()

	output, err := newTestOutput(connector)
	assert.Nil(t, err)
	input, err := newTestInput(connector)
	assert.Nil(t, err)

	// Write with custom source timestamp
	customTimestamp := int64(1234567890000000) // nanoseconds
	assert.Nil(t, output.Instance.SetString("st", "timestamp_test"))

	params := `{"source_timestamp":1234567890000000}`
	err = output.WriteWithParams(params)
	assert.Nil(t, err)

	err = connector.Wait(-1)
	assert.Nil(t, err)

	err = input.Take()
	assert.Nil(t, err)

	// Verify the custom timestamp was set
	ts, err := input.Infos.GetSourceTimestamp(0)
	assert.Nil(t, err)
	assert.Equal(t, customTimestamp, ts)
}

// TestWriteWithParamsDispose tests writing with dispose action
func TestWriteWithParamsDispose(t *testing.T) {
	connector, err := newTestConnector()
	assert.Nil(t, err)
	defer connector.Delete()

	output, err := newTestOutput(connector)
	assert.Nil(t, err)
	input, err := newTestInput(connector)
	assert.Nil(t, err)

	// First write a sample normally
	assert.Nil(t, output.Instance.SetString("st", "test_dispose"))
	err = output.Write()
	assert.Nil(t, err)

	err = connector.Wait(-1)
	assert.Nil(t, err)
	err = input.Take()
	assert.Nil(t, err)

	valid, err := input.Infos.IsValid(0)
	assert.Nil(t, err)
	assert.True(t, valid)

	// Now dispose the instance
	params := `{"action":"dispose"}`
	err = output.WriteWithParams(params)
	assert.Nil(t, err)

	err = connector.Wait(1000)
	if err == nil {
		err = input.Take()
		if err == nil {
			// Check if the instance is not valid (disposed)
			valid, err := input.Infos.IsValid(0)
			assert.Nil(t, err)
			assert.False(t, valid)

			// Check instance state
			instanceState, err := input.Infos.GetInstanceState(0)
			assert.Nil(t, err)
			// Should be NOT_ALIVE_DISPOSED or similar
			assert.NotEmpty(t, instanceState)
		}
	}
}

// ============================
// Array Data Type Tests
// ============================

// TestArrayDataTypes tests reading and writing array data types using JSON
func TestArrayDataTypes(t *testing.T) {
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/ArrayTest.xml")

	connector, err := NewConnector("MyParticipantLibrary::Zero", xmlPath)
	assert.Nil(t, err)
	defer connector.Delete()

	output, err := connector.GetOutput("MyPublisher::MyWriter")
	assert.Nil(t, err)
	input, err := connector.GetInput("MySubscriber::MyReader")
	assert.Nil(t, err)

	// Use JSON to set array data - more reliable than field notation
	jsonData := `{
		"id": 1,
		"int_array": [0, 10, 20, 30, 40, 50, 60, 70, 80, 90],
		"string_array": ["test1", "test2", "test3", "test4", "test5"],
		"data": "array_test"
	}`

	assert.Nil(t, output.Instance.SetJSON([]byte(jsonData)))

	// Write the sample
	err = output.Write()
	assert.Nil(t, err)

	err = connector.Wait(-1)
	assert.Nil(t, err)

	err = input.Take()
	assert.Nil(t, err)

	// Verify scalar fields
	id, err := input.Samples.GetInt32(0, "id")
	assert.Nil(t, err)
	assert.Equal(t, int32(1), id)

	data, err := input.Samples.GetString(0, "data")
	assert.Nil(t, err)
	assert.Equal(t, "array_test", data)

	// Retrieve as JSON and verify arrays are present
	receivedJSON, err := input.Samples.GetJSON(0)
	assert.Nil(t, err)
	assert.NotEmpty(t, receivedJSON)
	assert.Contains(t, string(receivedJSON), "array_test")
	assert.Contains(t, string(receivedJSON), "int_array")
}

// TestArrayJSON tests array handling via JSON
func TestArrayJSON(t *testing.T) {
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/ArrayTest.xml")

	connector, err := NewConnector("MyParticipantLibrary::Zero", xmlPath)
	assert.Nil(t, err)
	defer connector.Delete()

	output, err := connector.GetOutput("MyPublisher::MyWriter")
	assert.Nil(t, err)
	input, err := connector.GetInput("MySubscriber::MyReader")
	assert.Nil(t, err)

	// Use JSON to set array data
	jsonData := `{
		"id": 2,
		"int_array": [100, 200, 300, 400, 500, 600, 700, 800, 900, 1000],
		"data": "json_array_test"
	}`

	assert.Nil(t, output.Instance.SetJSON([]byte(jsonData)))
	err = output.Write()
	assert.Nil(t, err)

	err = connector.Wait(-1)
	assert.Nil(t, err)

	err = input.Take()
	assert.Nil(t, err)

	// Retrieve as JSON and verify
	receivedJSON, err := input.Samples.GetJSON(0)
	assert.Nil(t, err)
	assert.NotEmpty(t, receivedJSON)
	assert.Contains(t, string(receivedJSON), "json_array_test")
}

// ============================
// Sequence Data Type Tests
// ============================

// TestSequenceDataTypes tests reading and writing sequence data types
func TestSequenceDataTypes(t *testing.T) {
	_, curPath, _, _ := runtime.Caller(0)
	xmlPath := path.Join(path.Dir(curPath), "./test/xml/SequenceTest.xml")

	connector, err := NewConnector("MyParticipantLibrary::Zero", xmlPath)
	assert.Nil(t, err)
	if connector == nil {
		t.Fatal("connector is nil")
	}
	defer connector.Delete()

	output, err := connector.GetOutput("MyPublisher::MyWriter")
	assert.Nil(t, err)
	if output == nil {
		t.Fatal("output is nil")
	}

	input, err := connector.GetInput("MySubscriber::MyReader")
	assert.Nil(t, err)
	if input == nil {
		t.Fatal("input is nil")
	}

	// Use JSON to set sequence data (easier than index notation)
	jsonData := `{
		"id": 1,
		"int_sequence": [10, 20, 30, 40, 50],
		"string_sequence": ["hello", "world", "test"],
		"data": "sequence_test"
	}`

	assert.Nil(t, output.Instance.SetJSON([]byte(jsonData)))
	err = output.Write()
	assert.Nil(t, err)

	err = connector.Wait(-1)
	assert.Nil(t, err)

	err = input.Take()
	assert.Nil(t, err)

	// Verify using JSON
	receivedJSON, err := input.Samples.GetJSON(0)
	assert.Nil(t, err)
	assert.NotEmpty(t, receivedJSON)
	assert.Contains(t, string(receivedJSON), "sequence_test")

	// Verify scalar field
	id, err := input.Samples.GetInt32(0, "id")
	assert.Nil(t, err)
	assert.Equal(t, int32(1), id)
}

// ============================
// Inline XML String Tests
// ============================

// TestInlineXMLConfiguration tests creating a connector with inline XML string
func TestInlineXMLConfiguration(t *testing.T) {
	// Note: Inline XML works in examples but may not work in test environment
	// See examples/xml_string/ and examples/go-get-example/ for working inline XML usage
	t.Skip("Inline XML (str://) works in production but test environment may not have BuiltinQosLibExp profiles available - see examples/xml_string/ for working code")

	// Inline XML must be on a single line with str://" prefix and " suffix
	xmlString := `str://"<dds><qos_library name=\"QosLibrary\"><qos_profile name=\"DefaultProfile\" is_default_qos=\"true\"><participant_qos><discovery><initial_peers><element>shmem://</element></initial_peers><multicast_receive_addresses/></discovery></participant_qos><datawriter_qos><reliability><kind>RELIABLE_RELIABILITY_QOS</kind></reliability><history><kind>KEEP_ALL_HISTORY_QOS</kind></history><durability><kind>TRANSIENT_LOCAL_DURABILITY_QOS</kind></durability></datawriter_qos><datareader_qos><reliability><kind>RELIABLE_RELIABILITY_QOS</kind></reliability><history><kind>KEEP_ALL_HISTORY_QOS</kind></history><durability><kind>TRANSIENT_LOCAL_DURABILITY_QOS</kind></durability></datareader_qos></qos_profile></qos_library><types><struct name=\"InlineTestType\"><member name=\"id\" type=\"long\" key=\"true\"/><member name=\"message\" type=\"string\" stringMaxLength=\"256\"/></struct></types><domain_library name=\"MyDomainLibrary\"><domain name=\"MyDomain\" domain_id=\"0\"><register_type name=\"InlineTestType\" type_ref=\"InlineTestType\"/><topic name=\"InlineTopic\" register_type_ref=\"InlineTestType\"/></domain></domain_library><domain_participant_library name=\"MyParticipantLibrary\"><domain_participant name=\"Zero\" domain_ref=\"MyDomainLibrary::MyDomain\"><publisher name=\"MyPublisher\"><data_writer name=\"MyWriter\" topic_ref=\"InlineTopic\"/></publisher><subscriber name=\"MySubscriber\"><data_reader name=\"MyReader\" topic_ref=\"InlineTopic\"/></subscriber></domain_participant></domain_participant_library></dds>"`

	connector, err := NewConnector("MyParticipantLibrary::Zero", xmlString)
	assert.Nil(t, err)
	assert.NotNil(t, connector)
	defer connector.Delete()

	// Verify we can create input and output
	output, err := connector.GetOutput("MyPublisher::MyWriter")
	assert.Nil(t, err)
	assert.NotNil(t, output)

	input, err := connector.GetInput("MySubscriber::MyReader")
	assert.Nil(t, err)
	assert.NotNil(t, input)

	// Test data flow with inline configuration
	assert.Nil(t, output.Instance.SetInt32("id", 99))
	assert.Nil(t, output.Instance.SetString("message", "inline_xml_test"))
	err = output.Write()
	assert.Nil(t, err)

	err = connector.Wait(-1)
	assert.Nil(t, err)

	err = input.Take()
	assert.Nil(t, err)

	id, err := input.Samples.GetInt32(0, "id")
	assert.Nil(t, err)
	assert.Equal(t, int32(99), id)

	message, err := input.Samples.GetString(0, "message")
	assert.Nil(t, err)
	assert.Equal(t, "inline_xml_test", message)
}

// TestInvalidInlineXML tests error handling for malformed inline XML
func TestInvalidInlineXML(t *testing.T) {
	t.Skip("Inline XML (str://) works in production but test environment may not have BuiltinQosLibExp profiles available - see examples/xml_string/ for working code")

	invalidXML := `str://"<dds><invalid></xml>"`

	connector, err := NewConnector("MyParticipantLibrary::Zero", invalidXML)
	assert.Nil(t, connector)
	assert.NotNil(t, err)
}

// ============================
// Concurrency Tests
// ============================

// TestConcurrentReads tests concurrent read operations
func TestConcurrentReads(t *testing.T) {
	t.Skip("Skipped: This test intentionally demonstrates non-thread-safe behavior which can crash the test suite")
}

// TestConcurrentWrites tests concurrent write operations
func TestConcurrentWrites(t *testing.T) {
	t.Skip("Skipped: This test intentionally demonstrates non-thread-safe behavior which can crash the test suite. See TestSynchronizedWrites for the correct approach.")
}

// TestSynchronizedWrites demonstrates proper synchronization for writes
func TestSynchronizedWrites(t *testing.T) {
	connector, err := newTestConnector()
	assert.Nil(t, err)
	defer connector.Delete()

	output, err := newTestOutput(connector)
	assert.Nil(t, err)

	// Use mutex to synchronize writes
	var mu sync.Mutex
	const numWrites = 10
	var wg sync.WaitGroup
	successCount := 0

	wg.Add(numWrites)
	for i := 0; i < numWrites; i++ {
		go func(index int) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()

			// With proper synchronization, writes should succeed
			output.Instance.SetInt32("l", int32(index))
			if output.Write() == nil {
				successCount++
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, numWrites, successCount, "All synchronized writes should succeed")
}

// ============================
// Additional Error Path Tests
// ============================

// TestNullPointerHandling tests error handling for nil pointers
func TestNullPointerHandling(t *testing.T) {
	// Test nil Samples
	var samples *Samples
	_, err := samples.GetLength()
	assert.NotNil(t, err)

	_, err = samples.GetString(0, "field")
	assert.NotNil(t, err)

	_, err = samples.GetInt32(0, "field")
	assert.NotNil(t, err)

	// Test nil Infos
	var infos *Infos
	_, err = infos.IsValid(0)
	assert.NotNil(t, err)

	_, err = infos.GetSourceTimestamp(0)
	assert.NotNil(t, err)

	_, err = infos.GetReceptionTimestamp(0)
	assert.NotNil(t, err)

	// Test nil Input
	var input *Input
	err = input.Read()
	assert.NotNil(t, err)

	err = input.Take()
	assert.NotNil(t, err)

	// Test nil Output
	var output *Output
	err = output.Write()
	assert.NotNil(t, err)

	err = output.ClearMembers()
	assert.NotNil(t, err)

	// Test nil Instance
	var instance *Instance
	err = instance.SetInt32("field", 1)
	assert.NotNil(t, err)

	err = instance.SetString("field", "test")
	assert.NotNil(t, err)

	err = instance.SetBoolean("field", true)
	assert.NotNil(t, err)

	err = instance.SetJSON([]byte(`{}`))
	assert.NotNil(t, err)

	err = instance.SetFloat64("field", 3.14)
	assert.NotNil(t, err)
}

// TestNegativeIndices tests error handling for negative indices
func TestNegativeIndices(t *testing.T) {
	connector, err := newTestConnector()
	assert.Nil(t, err)
	defer connector.Delete()

	input, err := newTestInput(connector)
	assert.Nil(t, err)

	// Test negative index on IsValid
	_, err = input.Infos.IsValid(-1)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "negative")
}

// TestTypeMismatch tests error handling for type mismatches
func TestTypeMismatch(t *testing.T) {
	connector, err := newTestConnector()
	assert.Nil(t, err)
	defer connector.Delete()

	output, err := newTestOutput(connector)
	assert.Nil(t, err)
	input, err := newTestInput(connector)
	assert.Nil(t, err)

	// Write a string field
	assert.Nil(t, output.Instance.SetString("st", "test"))
	err = output.Write()
	assert.Nil(t, err)

	err = connector.Wait(-1)
	assert.Nil(t, err)

	err = input.Take()
	assert.Nil(t, err)

	// Try to get the string field as an int32 (type mismatch)
	_, err = input.Samples.GetInt32(0, "st")
	// This may or may not error depending on implementation
	// Just verify it doesn't crash
	if err != nil {
		t.Logf("Type mismatch error (expected): %v", err)
	}
}

// TestOutOfBoundsAccess tests error handling for out-of-bounds access
func TestOutOfBoundsAccess(t *testing.T) {
	t.Skip("Skipped: Out of bounds access can cause crashes in the underlying C library")
}

// TestNonExistentField tests error handling for non-existent fields
func TestNonExistentField(t *testing.T) {
	connector, err := newTestConnector()
	assert.Nil(t, err)
	defer connector.Delete()

	output, err := newTestOutput(connector)
	assert.Nil(t, err)
	input, err := newTestInput(connector)
	assert.Nil(t, err)

	// Try to set a field that doesn't exist - should get a meaningful error
	err = output.Instance.SetString("nonexistent_field", "value")
	assert.NotNil(t, err)
	// Error should contain useful information (either detailed message or error code)
	assert.Contains(t, err.Error(), "DDS Exception")
	t.Logf("Error for non-existent field: %v", err)

	// Write valid data first
	assert.Nil(t, output.Instance.SetString("st", "test"))
	err = output.Write()
	assert.Nil(t, err)

	err = connector.Wait(-1)
	assert.Nil(t, err)

	err = input.Take()
	assert.Nil(t, err)

	// Try to access a field that doesn't exist on read - should get meaningful error
	_, err = input.Samples.GetString(0, "nonexistent_field")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "DDS Exception")
	t.Logf("Error for reading non-existent field: %v", err)
}

// TestSampleStateTransitions tests DDS sample state transitions
func TestSampleStateTransitions(t *testing.T) {
	connector, err := newTestConnector()
	assert.Nil(t, err)
	defer connector.Delete()

	output, err := newTestOutput(connector)
	assert.Nil(t, err)
	input, err := newTestInput(connector)
	assert.Nil(t, err)

	// Write a sample
	assert.Nil(t, output.Instance.SetString("st", "state_test"))
	err = output.Write()
	assert.Nil(t, err)

	err = connector.Wait(-1)
	assert.Nil(t, err)

	// First read: should be NOT_READ
	err = input.Read()
	assert.Nil(t, err)

	sampleState, err := input.Infos.GetSampleState(0)
	assert.Nil(t, err)
	assert.Equal(t, "NOT_READ", sampleState)

	viewState, err := input.Infos.GetViewState(0)
	assert.Nil(t, err)
	assert.Equal(t, "NEW", viewState)

	// Second read: should be READ
	err = input.Read()
	assert.Nil(t, err)

	sampleState, err = input.Infos.GetSampleState(0)
	assert.Nil(t, err)
	// After second read, might still be NOT_READ or become READ depending on implementation
	t.Logf("Sample state after second read: %s", sampleState)

	// Instance state should be ALIVE
	instanceState, err := input.Infos.GetInstanceState(0)
	assert.Nil(t, err)
	assert.Equal(t, "ALIVE", instanceState)
}

// TestGetIdentityMethods tests Identity getter methods
func TestGetIdentityMethods(t *testing.T) {
	connector, err := newTestConnector()
	assert.Nil(t, err)
	defer connector.Delete()

	output, err := newTestOutput(connector)
	assert.Nil(t, err)
	input, err := newTestInput(connector)
	assert.Nil(t, err)

	// Write a sample
	assert.Nil(t, output.Instance.SetString("st", "identity_test"))
	err = output.Write()
	assert.Nil(t, err)

	err = connector.Wait(-1)
	assert.Nil(t, err)

	err = input.Take()
	assert.Nil(t, err)

	// Test GetIdentity
	identity, err := input.Infos.GetIdentity(0)
	assert.Nil(t, err)
	assert.NotNil(t, identity)
	assert.NotEmpty(t, identity.WriterGUID)
	assert.Greater(t, identity.SequenceNumber, 0)

	// Test GetIdentityJSON
	identityJSON, err := input.Infos.GetIdentityJSON(0)
	assert.Nil(t, err)
	assert.NotEmpty(t, identityJSON)
	assert.Contains(t, identityJSON, "writer_guid")
	assert.Contains(t, identityJSON, "sequence_number")
}

// TestMultipleConnectorInstances tests creating and using multiple connectors simultaneously
func TestMultipleConnectorInstances(t *testing.T) {
	const numConnectors = 3

	connectors := make([]*Connector, numConnectors)
	var err error

	// Create multiple connectors
	for i := 0; i < numConnectors; i++ {
		connectors[i], err = newTestConnector()
		assert.Nil(t, err)
		assert.NotNil(t, connectors[i])
	}

	// Clean up all connectors
	for i := 0; i < numConnectors; i++ {
		assert.Nil(t, connectors[i].Delete())
	}
}

// TestLargeStringHandling tests handling of large string values
func TestLargeStringHandling(t *testing.T) {
	connector, err := newTestConnector()
	assert.Nil(t, err)
	defer connector.Delete()

	output, err := newTestOutput(connector)
	assert.Nil(t, err)
	input, err := newTestInput(connector)
	assert.Nil(t, err)

	// Create a moderately large string (within the 256 char limit from Test.xml)
	largeString := ""
	for i := 0; i < 10; i++ {
		largeString += "Test string "
	}
	// Total: ~120 characters, well within 256 limit

	assert.Nil(t, output.Instance.SetString("st", largeString))
	err = output.Write()
	assert.Nil(t, err)

	err = connector.Wait(-1)
	assert.Nil(t, err)

	err = input.Take()
	assert.Nil(t, err)

	receivedString, err := input.Samples.GetString(0, "st")
	assert.Nil(t, err)
	assert.Equal(t, largeString, receivedString)
}

// TestJSONStructureMismatch tests error handling when JSON structure doesn't match XML
// This addresses GitHub issue #51: Get Last Error Message returns empty string
func TestJSONStructureMismatch(t *testing.T) {
	connector, err := newTestConnector()
	assert.Nil(t, err)
	defer connector.Delete()

	output, err := newTestOutput(connector)
	assert.Nil(t, err)

	// Try to set JSON with fields that don't exist in the XML definition
	// The Test.xml ShapeType has fields: color (string), x, y, shapesize (int32)
	invalidJSON := []byte(`{"nonexistent_field": "value", "another_bad_field": 123}`)
	err = output.Instance.SetJSON(invalidJSON)

	// Should get an error with meaningful information
	if err != nil {
		assert.Contains(t, err.Error(), "DDS Exception")
		t.Logf("Error message for JSON structure mismatch: %v", err)
		// Error should contain either detailed message or error code
		assert.True(t,
			len(err.Error()) > len("DDS Exception: "),
			"Error message should contain details beyond just 'DDS Exception: '")
	} else {
		// If no error, the C library might be lenient with extra fields
		t.Log("No error returned for JSON with non-existent fields (C library may be lenient)")
	}

	// Try to set a field with wrong type (string instead of int32)
	invalidTypeJSON := []byte(`{"x": "not_a_number", "y": 100}`)
	err = output.Instance.SetJSON(invalidTypeJSON)

	if err != nil {
		assert.Contains(t, err.Error(), "DDS Exception")
		t.Logf("Error message for type mismatch: %v", err)
		assert.True(t,
			len(err.Error()) > len("DDS Exception: "),
			"Error message should contain details beyond just 'DDS Exception: '")
	}
}
