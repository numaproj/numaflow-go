package impl

import (
	"time"

	functionsdk "github.com/KeranYang/numaflow-go/pkg/function"
)

func FilterEventTime(keys []string, d functionsdk.Datum) functionsdk.MessageTs {
	janFirst2022 := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	janFirst2023 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if d.EventTime().Before(janFirst2022) {
		return functionsdk.MessageTsBuilder().Append(functionsdk.MessageTToDrop())
	} else if d.EventTime().Before(janFirst2023) {
		return functionsdk.MessageTsBuilder().Append(functionsdk.NewMessageT(d.Value(), janFirst2022).WithTags([]string{"within_year_2022"}))
	} else {
		return functionsdk.MessageTsBuilder().Append(functionsdk.NewMessageT(d.Value(), janFirst2023).WithTags([]string{"after_year_2022"}))
	}
}
