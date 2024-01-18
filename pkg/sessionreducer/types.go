package sessionreducer

import "time"

// handlerDatum implements the Datum interface and is used in the SessionReduce functions.
type handlerDatum struct {
	value     []byte
	eventTime time.Time
	watermark time.Time
}

func NewHandlerDatum(value []byte, eventTime time.Time, watermark time.Time) Datum {
	return &handlerDatum{
		value:     value,
		eventTime: eventTime,
		watermark: watermark,
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
