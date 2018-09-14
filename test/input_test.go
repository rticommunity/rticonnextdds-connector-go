package connector_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvalidDR(t *testing.T) {
	invalidReaderName := "invalidDR"

	connector := newTestConnector()
	input, err := connector.GetInput(invalidReaderName)
        assert.Nil(t, input)
        assert.NotNil(t, err)
}

func TestCreateDR(t *testing.T) {
        readerName := "MySubscriber::MySquareReader"

        connector := newTestConnector()
        input, err := connector.GetInput(readerName)
        assert.NotNil(t, input)
        assert.NotNil(t, input.Samples)
        assert.NotNil(t, input.Infos)
        assert.Nil(t, err)
	deleteTestConnector(connector)
}
