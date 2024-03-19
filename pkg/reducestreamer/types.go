package reducestreamer

import "time"

// handlerDatum implements the Datum interface and is used in the reduceStream functions.
type handlerDatum struct {
	value     []byte
	eventTime time.Time
	watermark time.Time
	headers   map[string]string
}

func NewHandlerDatum(value []byte, eventTime time.Time, watermark time.Time, headers map[string]string) Datum {
	return &handlerDatum{
		value:     value,
		eventTime: eventTime,
		watermark: watermark,
		headers:   headers,
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

// intervalWindow implements IntervalWindow interface which will be passed as metadata
// to reduce handlers
type intervalWindow struct {
	startTime time.Time
	endTime   time.Time
}

func NewIntervalWindow(startTime time.Time, endTime time.Time) IntervalWindow {
	return &intervalWindow{
		startTime: startTime,
		endTime:   endTime,
	}
}

func (i *intervalWindow) StartTime() time.Time {
	return i.startTime
}

func (i *intervalWindow) EndTime() time.Time {
	return i.endTime
}

// metadata implements Metadata interface which will be passed to reduceStream function.
type metadata struct {
	intervalWindow IntervalWindow
}

func NewMetadata(window IntervalWindow) Metadata {
	return &metadata{intervalWindow: window}
}

func (m *metadata) IntervalWindow() IntervalWindow {
	return m.intervalWindow
}
