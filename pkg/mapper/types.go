package mapper

import "time"

// handlerDatum implements the Datum interface and is used in the map functions.
type handlerDatum struct {
	value     []byte
	eventTime time.Time
	watermark time.Time
	headers   map[string]string
	metadata  Metadata
}

func NewHandlerDatum(value []byte, eventTime time.Time, watermark time.Time, headers map[string]string, metadata Metadata) Datum {
	return &handlerDatum{
		value:     value,
		eventTime: eventTime,
		watermark: watermark,
		headers:   headers,
		metadata:  metadata,
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

func (h *handlerDatum) Metadata() Metadata {
	return h.metadata
}
