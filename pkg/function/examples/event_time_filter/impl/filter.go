package impl

import (
	"time"

	"github.com/numaproj/numaflow-go/pkg/source"
)

func FilterEventTime(keys []string, d source.Datum) source.MessageTs {
	janFirst2022 := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	janFirst2023 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if d.EventTime().Before(janFirst2022) {
		return source.MessageTsBuilder().Append(source.MessageTToDrop())
	} else if d.EventTime().Before(janFirst2023) {
		return source.MessageTsBuilder().Append(source.NewMessageT(d.Value(), janFirst2022).WithTags([]string{"within_year_2022"}))
	} else {
		return source.MessageTsBuilder().Append(source.NewMessageT(d.Value(), janFirst2023).WithTags([]string{"after_year_2022"}))
	}
}
