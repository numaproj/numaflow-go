package mapper

import "time"

// handlerDatum implements the Datum interface and is used in the map functions.
type handlerDatum struct {
	value          []byte
	eventTime      time.Time
	watermark      time.Time
	headers        map[string]string
	userMetadata   UserMetadata
	systemMetadata SystemMetadata
}

func NewHandlerDatum(value []byte, eventTime time.Time, watermark time.Time, headers map[string]string, userMetadata UserMetadata, systemMetadata SystemMetadata) Datum {
	return &handlerDatum{
		value:          value,
		eventTime:      eventTime,
		watermark:      watermark,
		headers:        headers,
		userMetadata:   userMetadata,
		systemMetadata: systemMetadata,
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

func (h *handlerDatum) UserMetadata() UserMetadata {
	return h.userMetadata
}

func (h *handlerDatum) SystemMetadata() SystemMetadata {
	return h.systemMetadata
}
