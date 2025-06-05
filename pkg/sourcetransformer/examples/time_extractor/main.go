package main

import (
	"context"
	"log"
	"time"
	"timeextractionfilter/util"

	"github.com/araddon/dateparse"
	"github.com/numaproj/numaflow-go/pkg/sourcetransformer"
)

type TimeExtractor struct {
	filterExpr      string
	eventTimeExpr   string
	eventTimeFormat string
}

func (f TimeExtractor) Transform(_ context.Context, keys []string, d sourcetransformer.Datum) sourcetransformer.Messages {
	resultMsg, err := f.apply(d.EventTime(), d.Value(), keys)
	if err != nil {
		log.Fatalf("Filter map function apply got an error: %v", err)
	}
	return sourcetransformer.MessagesBuilder().Append(resultMsg)
}

func (f TimeExtractor) apply(et time.Time, payload []byte, keys []string) (sourcetransformer.Message, error) {
	result, err := util.EvalBool(f.filterExpr, payload)
	if err != nil {
		return sourcetransformer.MessageToDrop(et), err
	}
	if result {
		timeStr, err := util.EvalStr(f.eventTimeExpr, payload)
		if err != nil {
			return sourcetransformer.NewMessage(payload, et).WithKeys(keys), err
		}
		var newEventTime time.Time
		time.Local, _ = time.LoadLocation("UTC")
		if f.eventTimeFormat != "" {
			newEventTime, err = time.Parse(f.eventTimeFormat, timeStr)
		} else {
			newEventTime, err = dateparse.ParseStrict(timeStr)
		}
		if err != nil {
			return sourcetransformer.NewMessage(payload, et).WithKeys(keys), err
		} else {
			return sourcetransformer.NewMessage(payload, newEventTime).WithKeys(keys), nil
		}
	}
	return sourcetransformer.MessageToDrop(et), nil
}

func main() {
	f := TimeExtractor{
		filterExpr:    `int(json(payload).id) < 100`,
		eventTimeExpr: `json(payload).time`,
	}
	err := sourcetransformer.NewServer(f).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start source transform server: ", err)
	}
}
