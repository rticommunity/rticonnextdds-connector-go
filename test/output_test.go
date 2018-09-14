package connector_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvalidWriter(t *testing.T) {
	invalidWriterName := "invalidWriter"

	connector := newTestConnector()
	output, err := connector.GetOutput(invalidWriterName)
        assert.Nil(t, output)
        assert.NotNil(t, err)
}

func TestCreateWriter(t *testing.T) {
        writerName := "MyPublisher::MySquareWriter"

        connector := newTestConnector()
        output, err := connector.GetOutput(writerName)
        assert.NotNil(t, output)
        assert.NotNil(t, output.Instance)
        assert.Nil(t, err)
	deleteTestConnector(connector)
}
