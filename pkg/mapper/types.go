package mapper

import "time"

// handlerDatum implements the Datum interface and is used in the map functions.
type handlerDatum struct {
	value     []byte
	eventTime time.Time
	watermark time.Time
	headers   map[string]string
	id        string
	keys      []string
}

func NewHandlerDatum(value []byte, eventTime time.Time, watermark time.Time, headers map[string]string, id string, keys []string) Datum {
	return &handlerDatum{
		value:     value,
		eventTime: eventTime,
		watermark: watermark,
		headers:   headers,
		id:        id,
		keys:      keys,
	}
}

func (h *handlerDatum) Value() []byte {
	return h.value
}

func (h *handlerDatum) EventTime() time.Time {
	return h.eventTime
}

func (h *handlerDatum) Watermark() time.Time {
	return h.watermark
}

func (h *handlerDatum) Headers() map[string]string {
	return h.headers
}

func (h *handlerDatum) Id() string {
	return h.id
}

func (h *handlerDatum) Keys() []string {
	return h.keys
}
