package function

import (
	"context"

	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/function/v1"
)

type functionService struct {
	v1.UnimplementedFunctionServiceServer

	mapper  MapHandler
	reducer ReduceHandler
}

func (fs *functionService) DoFn(ctx context.Context, d *v1.Datum) (*v1.DatumList, error) {
	msgs, err := fs.mapper.HandleDo(ctx, d.GetKey(), d.GetValue())
	if err != nil {
		return nil, err
	}
	l := []*v1.Datum{}
	for _, m := range msgs.Items() {
		l = append(l, &v1.Datum{
			Key:      m.Key,
			Value:    m.Value,
			PaneInfo: d.PaneInfo, // should provide a way for the user to update the event time.
		})
	}
	datumList := &v1.DatumList{
		Items: l,
	}
	return datumList, nil
}

func (fs *functionService) ReduceFn(v1.FunctionService_ReduceFnServer) error {

	return nil
}
