package message

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestString(t *testing.T) {
// 	expectedMsg := ChatMessage{
// 		Author: "dan",
// 		Message: "btw i'm dan",
// 	}

// 	b, err := expectedMsg.WriteTo()
// 	require.NoError(t, err)

// 	var newMsg ChatMessage
// 	newMsg.UnmarshalBinary(b)

// 	assert.Equal(t, expectedMsg.Author, newMsg.Author)
// 	assert.Equal(t, expectedMsg.Message, newMsg.Message)
// }

// func TestMarshalUnMarshalChatMessage_AuthorTooLong(t *testing.T) {
// 	author := ""
// 	for i := 0; i < 256; i++ {
// 		author += "A"

// 	}
// 	expectedMsg := ChatMessage{
// 		Author: author,
// 		Message: "btw i'm dan",
// 	}

// 	_, err := expectedMsg.MarshalBinary()
// 	assert.Error(t, err)
// }

// func TestMarshalUnMarshalChatMessage_MessageTooLong(t *testing.T) {
// 	message := ""
// 	for i := 0; i < 256; i++ {
// 		message += "A"

// 	}
// 	expectedMsg := ChatMessage{
// 		Author: "dan",
// 		Message: message,
// 	}

// 	_, err := expectedMsg.MarshalBinary()
// 	assert.Error(t, err)
// }