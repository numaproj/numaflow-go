package impl

import (
	"time"

	"github.com/numaproj/numaflow-go/pkg/sourcetransformer"
)

func FilterEventTime(_ []string, d sourcetransformer.Datum) sourcetransformer.Messages {
	janFirst2022 := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	janFirst2023 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if d.EventTime().Before(janFirst2022) {
		return sourcetransformer.MessagesBuilder().Append(sourcetransformer.MessageToDrop(d.EventTime()))
	} else if d.EventTime().Before(janFirst2023) {
		return sourcetransformer.MessagesBuilder().Append(sourcetransformer.NewMessage(d.Value(), janFirst2022).WithTags([]string{"within_year_2022"}))
	} else {
		return sourcetransformer.MessagesBuilder().Append(sourcetransformer.NewMessage(d.Value(), janFirst2023).WithTags([]string{"after_year_2022"}))
	}
}
