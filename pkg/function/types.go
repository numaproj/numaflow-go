package function

import "time"

// handlerDatum implements the Datum interface and is used in the map and reduce handlers.
type handlerDatum struct {
	value     []byte
	eventTime time.Time
	watermark time.Time
	metadata  handlerDatumMetadata
}

func NewHandlerDatum(value []byte, eventTime time.Time, watermark time.Time, metadata DatumMetadata) Datum {
	return &handlerDatum{
		value:     value,
		eventTime: eventTime,
		watermark: watermark,
		metadata:  metadata.(handlerDatumMetadata),
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

func (h *handlerDatum) Metadata() DatumMetadata {
	return h.metadata
}

// handlerDatumMetadata implements the DatumMetadata interface and is used in the map and reduce handlers.
type handlerDatumMetadata struct {
	id           string
	numDelivered uint64
}

func NewHandlerDatumMetadata(id string, numDelivered uint64) DatumMetadata {
	return handlerDatumMetadata{
		id:           id,
		numDelivered: numDelivered,
	}
}

// ID returns the ID of the datum.
func (h handlerDatumMetadata) ID() string {
	return h.id
}

// NumDelivered returns the number of times the datum has been delivered.
func (h handlerDatumMetadata) NumDelivered() uint64 {
	return h.numDelivered
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

// metadata implements Metadata interface which will be passed to reduce handlers
type metadata struct {
	intervalWindow IntervalWindow
}

func NewMetadata(window IntervalWindow) Metadata {
	return &metadata{intervalWindow: window}
}

func (m *metadata) IntervalWindow() IntervalWindow {
	return m.intervalWindow
}
