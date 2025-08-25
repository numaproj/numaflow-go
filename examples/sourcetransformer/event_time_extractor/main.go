package main

import (
	"context"
	"event_time_extractor/util"
	"log"
	"time"

	"github.com/araddon/dateparse"
	"github.com/numaproj/numaflow-go/pkg/sourcetransformer"
)

type EventTimeExtractor struct {
	// expression is used to extract the string representation of the event time from message payload.
	// e.g. `json(payload).metadata.time`
	expression string
	// format specifies the layout of extracted time string.
	// with format, eventTimeExtractor uses the time.Parse function to translate the event time string representation to time.Time object.
	// otherwise if format is not specified, eventTimeExtractor uses dateparse to find format based on the time string.
	format string
}

func (e EventTimeExtractor) Transform(_ context.Context, keys []string, d sourcetransformer.Datum) sourcetransformer.Messages {
	resultMsg, err := e.apply(d.Value(), d.EventTime(), keys)
	if err != nil {
		log.Fatalf("event time extractor got an error: %v", err)
	}
	return sourcetransformer.MessagesBuilder().Append(resultMsg)
}

// apply compiles the payload to extract the new event time. If there is any error during extraction,
// we pass on the original input event time. Otherwise, we assign the new event time to the message.
func (e EventTimeExtractor) apply(payload []byte, et time.Time, keys []string) (sourcetransformer.Message, error) {
	timeStr, err := util.EvalStr(e.expression, payload)
	if err != nil {
		return sourcetransformer.NewMessage(payload, et).WithKeys(keys), err
	}

	var newEventTime time.Time
	time.Local, _ = time.LoadLocation("UTC")
	if e.format != "" {
		newEventTime, err = time.Parse(e.format, timeStr)
	} else {
		newEventTime, err = dateparse.ParseStrict(timeStr)
	}
	if err != nil {
		return sourcetransformer.NewMessage(payload, et).WithKeys(keys), err
	} else {
		return sourcetransformer.NewMessage(payload, newEventTime).WithKeys(keys), nil
	}
}

func main() {
	e := EventTimeExtractor{
		expression: `json(payload).item[1].time`,
	}
	err := sourcetransformer.NewServer(e).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start source transform server: ", err)
	}
}
