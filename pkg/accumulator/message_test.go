package accumulator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessageFromDatum(t *testing.T) {
	eventTime := time.UnixMilli(60000)
	watermark := time.UnixMilli(59000)
	datum := &handlerDatum{
		value:     []byte("hello"),
		id:        "test-id",
		eventTime: eventTime,
		watermark: watermark,
		keys:      []string{"key1", "key2"},
		headers:   map[string]string{"x": "y"},
	}

	message := MessageFromDatum(datum)

	assert.Equal(t, []byte("hello"), message.Value())
	assert.Equal(t, []string{"key1", "key2"}, message.Keys())
	assert.Nil(t, message.Tags())
	assert.Equal(t, "test-id", message.id)
	assert.Equal(t, eventTime, message.eventTime)
	assert.Equal(t, watermark, message.watermark)
	assert.Equal(t, map[string]string{"x": "y"}, message.headers)
}

func TestMessageToDrop(t *testing.T) {
	eventTime := time.UnixMilli(60000)
	watermark := time.UnixMilli(59000)
	datum := &handlerDatum{
		value:     []byte("hello"),
		id:        "test-id",
		eventTime: eventTime,
		watermark: watermark,
		keys:      []string{"key1", "key2"},
		headers:   map[string]string{"x": "y"},
	}

	message := MessageToDrop(datum)

	// The DROP tag must be set so the message is not forwarded downstream.
	assert.Equal(t, []string{DROP}, message.Tags())
	// No value is forwarded, but the identifying/watermark fields are carried over
	// so the accumulator can advance the watermark and release the tracked state.
	assert.Equal(t, []byte{}, message.Value())
	assert.Equal(t, []string{"key1", "key2"}, message.Keys())
	assert.Equal(t, "test-id", message.id)
	assert.Equal(t, eventTime, message.eventTime)
	assert.Equal(t, watermark, message.watermark)
	assert.Equal(t, map[string]string{"x": "y"}, message.headers)
}
