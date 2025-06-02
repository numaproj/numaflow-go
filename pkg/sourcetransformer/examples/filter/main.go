package main

import (
	"context"
	"filter/util"
	"log"
	"time"

	"github.com/numaproj/numaflow-go/pkg/sourcetransformer"
)

type Filter struct {
	expression string
}

func (f Filter) Transform(_ context.Context, keys []string, d sourcetransformer.Datum) sourcetransformer.Messages {
	resultMsg, err := f.apply(d.EventTime(), d.Value(), keys)
	if err != nil {
		log.Fatalf("Filter map function apply got an error: %v", err)
	}
	return sourcetransformer.MessagesBuilder().Append(resultMsg)
}

func (f Filter) apply(et time.Time, msg []byte, keys []string) (sourcetransformer.Message, error) {
	result, err := util.EvalBool(f.expression, msg)
	if err != nil {
		return sourcetransformer.MessageToDrop(et), err
	}
	if result {
		return sourcetransformer.NewMessage(msg, et).WithKeys(keys), nil
	}
	return sourcetransformer.MessageToDrop(et), nil
}

func main() {
	f := Filter{
		expression: `int(json(payload).id) < 100 && json(payload).msg == 'hello' && sprig.contains('good', string(json(payload).desc))`,
	}
	err := sourcetransformer.NewServer(f).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start source transform server: ", err)
	}
}
