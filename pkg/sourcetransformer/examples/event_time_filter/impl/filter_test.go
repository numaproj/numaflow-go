package impl

import (
	"testing"
	"time"

	functionsdk "github.com/numaproj/numaflow-go/pkg/mapper"
	"github.com/numaproj/numaflow-go/pkg/sourcetransformer"

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
	testKeys := []string{"test-key"}
	janFirst2022 := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	janFirst2023 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name           string
		input          functionsdk.Datum
		expectedOutput sourcetransformer.Messages
	}{
		{
			name:  "DatumWithEventTimeBefore2022GetsDropped",
			input: beforeYear2022Datum{},
			expectedOutput: sourcetransformer.Messages{
				sourcetransformer.MessageToDrop(time.Date(2021, 8, 1, 0, 14, 0, 0, time.UTC)),
			},
		},
		{
			name:  "DatumWithEventTimeWithin2022GetsKeyAndEventTimeUpdated",
			input: withinYear2022Datum{},
			expectedOutput: sourcetransformer.Messages{
				sourcetransformer.NewMessage([]byte("test-data"), janFirst2022).WithTags([]string{"within_year_2022"}),
			},
		},
		{
			name:  "DatumWithEventTimeAfter2022GetsKeyAndEventTimeUpdated",
			input: afterYear2022Datum{},
			expectedOutput: sourcetransformer.Messages{
				sourcetransformer.NewMessage([]byte("test-data"), janFirst2023).WithTags([]string{"after_year_2022"}),
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
