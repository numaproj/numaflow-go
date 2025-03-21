package accumulator

import "time"

// handlerDatum implements the Datum interface and is used in the accumulator functions.
type handlerDatum struct {
	value     []byte
	id        string
	eventTime time.Time
	watermark time.Time
	keys      []string
	tags      []string
	headers   map[string]string
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

func (h *handlerDatum) Keys() []string {
	return h.keys
}

func (h *handlerDatum) UpdateValue(value []byte) {
	h.value = value
}

func (h *handlerDatum) SetTags(tags []string) {
	h.tags = tags
}

func (h *handlerDatum) Headers() map[string]string {
	return h.headers
}

func (h *handlerDatum) ID() string {
	return h.id
}
