package main

import (
	"context"
	"filter/util"
	"log"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

type Filter struct {
	expression string
}

func (f Filter) Map(ctx context.Context, keys []string, d mapper.Datum) mapper.Messages {
	resultMsg, err := f.apply(d.Value())
	if err != nil {
		log.Fatalf("Filter map function apply got an error: %v", err)
	}
	return mapper.MessagesBuilder().Append(resultMsg)
}

func (f Filter) apply(msg []byte) (mapper.Message, error) {
	result, err := util.EvalBool(f.expression, msg)
	if err != nil {
		return mapper.MessageToDrop(), err
	}
	if result {
		return mapper.NewMessage(msg), nil
	}
	return mapper.MessageToDrop(), nil
}

func main() {

	f := Filter{
		expression: `int(json(payload).id) < 100 && json(payload).msg == 'hello' && sprig.contains('good', string(json(payload).desc))`,
	}
	err := mapper.NewServer(f).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start filter server: ", err)
	}
}
