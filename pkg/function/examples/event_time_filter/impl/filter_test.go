package impl

import (
	"testing"
	"time"

	"github.com/google/uuid"
	functionsdk "github.com/numaproj/numaflow-go/pkg/function"

	"github.com/stretchr/testify/assert"
)

type beforeYear2022Datum struct{}

func (d beforeYear2022Datum) Value() []byte {
	return []byte("test-data")
}
func (d beforeYear2022Datum) EventTime() time.Time {
	return time.Date(2021, 8, 1, 0, 14, 0, 0, time.UTC)

}
func (d beforeYear2022Datum) Watermark() time.Time {
	return time.Now()
}

func (d beforeYear2022Datum) Metadata() functionsdk.DatumMetadata {
	return beforeYear2022DatumMetadata{
		id:           "id",
		numDelivered: 1,
		uniqueID:     uuid.New().String(),
	}
}

type beforeYear2022DatumMetadata struct {
	id           string
	numDelivered uint64
	uniqueID     string
}

func (d beforeYear2022DatumMetadata) UUID() string {
	return d.uniqueID
}

func (d beforeYear2022DatumMetadata) ID() string {
	return d.id
}

func (d beforeYear2022DatumMetadata) NumDelivered() uint64 {
	return d.numDelivered
}

type withinYear2022Datum struct{}

func (d withinYear2022Datum) ID() string {
	return "id"
}

func (d withinYear2022Datum) Value() []byte {
	return []byte("test-data")
}
func (d withinYear2022Datum) EventTime() time.Time {
	return time.Date(2022, 6, 1, 0, 14, 0, 0, time.UTC)

}
func (d withinYear2022Datum) Watermark() time.Time {
	return time.Now()
}

func (d withinYear2022Datum) Metadata() functionsdk.DatumMetadata {
	return withinYear2022DatumMetadaa{
		id:           "id",
		numDelivered: 1,
		uniqueID:     uuid.New().String(),
	}
}

type withinYear2022DatumMetadaa struct {
	id           string
	numDelivered uint64
	uniqueID     string
}

func (d withinYear2022DatumMetadaa) ID() string {
	return d.id
}

func (d withinYear2022DatumMetadaa) NumDelivered() uint64 {
	return d.numDelivered
}

func (d withinYear2022DatumMetadaa) UUID() string {
	return d.uniqueID
}

type afterYear2022Datum struct{}

func (d afterYear2022Datum) Value() []byte {
	return []byte("test-data")
}
func (d afterYear2022Datum) EventTime() time.Time {
	return time.Date(2023, 6, 1, 13, 14, 0, 0, time.UTC)

}
func (d afterYear2022Datum) Watermark() time.Time {
	return time.Now()
}

func (d afterYear2022Datum) Metadata() functionsdk.DatumMetadata {
	return afterYear2022DatumMetadata{
		id:           "id",
		numDelivered: 1,
		uniqueID:     uuid.New().String(),
	}
}

type afterYear2022DatumMetadata struct {
	id           string
	numDelivered uint64
	uniqueID     string
}

func (d afterYear2022DatumMetadata) ID() string {
	return d.id
}

func (d afterYear2022DatumMetadata) NumDelivered() uint64 {
	return d.numDelivered
}

func (d afterYear2022DatumMetadata) UUID() string {
	return d.uniqueID
}

func Test_FilterEventTime(t *testing.T) {
	testKeys := []string{"test-key"}
	janFisrt2022 := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	janFirst2023 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name           string
		input          functionsdk.Datum
		expectedOutput functionsdk.MessageTs
	}{
		{
			name:  "DatumWithEventTimeBefore2022GetsDropped",
			input: beforeYear2022Datum{},
			expectedOutput: functionsdk.MessageTs{
				functionsdk.MessageTToDrop(),
			},
		},
		{
			name:  "DatumWithEventTimeWithin2022GetsKeyAndEventTimeUpdated",
			input: withinYear2022Datum{},
			expectedOutput: functionsdk.MessageTs{
				functionsdk.NewMessageT([]byte("test-data"), janFisrt2022).WithTags([]string{"within_year_2022"}),
			},
		},
		{
			name:  "DatumWithEventTimeAfter2022GetsKeyAndEventTimeUpdated",
			input: afterYear2022Datum{},
			expectedOutput: functionsdk.MessageTs{
				functionsdk.NewMessageT([]byte("test-data"), janFirst2023).WithTags([]string{"after_year_2022"}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := FilterEventTime(testKeys, tt.input)
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}
