package impl

import (
	"testing"
	"time"

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

type withinYear2022Datum struct{}

func (d withinYear2022Datum) Value() []byte {
	return []byte("test-data")
}
func (d withinYear2022Datum) EventTime() time.Time {
	return time.Date(2022, 6, 1, 0, 14, 0, 0, time.UTC)

}
func (d withinYear2022Datum) Watermark() time.Time {
	return time.Now()
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

func Test_FilterEventTime(t *testing.T) {
	testKey := "test-key"
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
				functionsdk.MessageTTo(janFisrt2022, "within_year_2022", []byte("test-data")),
			},
		},
		{
			name:  "DatumWithEventTimeAfter2022GetsKeyAndEventTimeUpdated",
			input: afterYear2022Datum{},
			expectedOutput: functionsdk.MessageTs{
				functionsdk.MessageTTo(janFirst2023, "after_year_2022", []byte("test-data")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := FilterEventTime(testKey, tt.input)
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}
