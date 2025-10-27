package message

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteReadString_Happy(t *testing.T) {
	expectedString := "here is a great string"

	w := new(bytes.Buffer)
	err := writeString(w, expectedString)
	require.NoError(t, err)

	r := bytes.NewReader(w.Bytes())

	result, err := readString(r)
	require.NoError(t, err)

	assert.Equal(t, expectedString, result)
}

func TestWriteReadString_255(t *testing.T) {
	var expectedString string
	for i := 0; i < 255; i++ {
		expectedString += "A"
	}

	w := new(bytes.Buffer)
	err := writeString(w, expectedString)
	require.NoError(t, err)

	r := bytes.NewReader(w.Bytes())

	result, err := readString(r)
	require.NoError(t, err)

	assert.Equal(t, expectedString, result)
}
func TestWriteReadString_TooLong(t *testing.T) {
	var expectedString string
	for i := 0; i < 256; i++ {
		expectedString += "A"
	}

	w := new(bytes.Buffer)
	err := writeString(w, expectedString)
	assert.NotNil(t, err)
}

