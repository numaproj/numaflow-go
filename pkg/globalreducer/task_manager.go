package globalreducer

//
//import (
//	"context"
//	"fmt"
//	"strings"
//	"sync"
//	"time"
//
//	v1 "github.com/numaproj/numaflow-go/pkg/apis/proto/sessionreduce/v1"
//)
//
//type reduceTask struct {
//	combinedKey    string
//	partition      *v1.Partition
//	reduceStreamer SessionReducer
//	watermark      time.Time
//	inputCh        chan Datum
//	outputCh       chan Message
//	Done           chan struct{}
//}
//
//func (rt *reduceTask) buildSessionReduceResponse(message Message) *v1.SessionReduceResponse {
//	response := &v1.SessionReduceResponse{
//		Result: &v1.SessionReduceResponse_Result{
//			Keys:  message.Keys(),
//			Value: message.Value(),
//			Tags:  message.Tags(),
//		},
//		Partition:   rt.partition,
//		CombinedKey: rt.combinedKey,
//		Cob:         false,
//	}
//
//	return response
//}
//
//func (rt *reduceTask) buildCobSessionReduceResponse() *v1.SessionReduceResponse {
//	response := &v1.SessionReduceResponse{
//		Result: &v1.SessionReduceResponse_Result{
//			Keys:  []string{},
//			Value: []byte{},
//			Tags:  []string{},
//		},
//		Partition:   rt.partition,
//		CombinedKey: rt.combinedKey,
//		Cob:         true,
//	}
//
//	return response
//}
//
//func (rt *reduceTask) uniqueKey() string {
//	return fmt.Sprintf("%d:%d:%s",
//		rt.partition.GetStart().AsTime().UnixMilli(),
//		rt.partition.GetEnd().AsTime().UnixMilli(),
//		rt.combinedKey)
//}
//
//type reduceTaskManager struct {
//	Tasks  map[string]*reduceTask
//	Output chan *v1.SessionReduceResponse
//	Mutex  sync.RWMutex
//}
//
//func newReduceTaskManager() *reduceTaskManager {
//	return &reduceTaskManager{
//		Tasks:  make(map[string]*reduceTask),
//		Output: make(chan *v1.SessionReduceResponse),
//	}
//}
//
//func (rtm *reduceTaskManager) CreateTask(ctx context.Context, request *v1.SessionReduceRequest, reduceStreamer SessionReducer, partition *v1.Partition) {
//	rtm.Mutex.Lock()
//	defer rtm.Mutex.Unlock()
//
//	task := &reduceTask{
//		combinedKey:    strings.Join(request.GetPayload().GetKeys(), delimiter),
//		partition:      partition,
//		reduceStreamer: reduceStreamer,
//		inputCh:        make(chan Datum),
//		outputCh:       make(chan Message),
//		Done:           make(chan struct{}),
//	}
//
//	go func() {
//		defer close(task.Done)
//		var wg sync.WaitGroup
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			for {
//				select {
//				case <-ctx.Done():
//					return
//				case message, ok := <-task.outputCh:
//					if !ok {
//						rtm.Output <- task.buildCobSessionReduceResponse()
//						return
//					}
//					rtm.Output <- task.buildSessionReduceResponse(message)
//				}
//			}
//		}()
//		reduceStreamer.SessionReduce(ctx, request.GetPayload().GetKeys(), task.inputCh, task.outputCh)
//		close(task.outputCh)
//		wg.Wait()
//	}()
//
//	rtm.Tasks[task.uniqueKey()] = task
//}
//
//func (rtm *reduceTaskManager) AppendToTask(request *v1.SessionReduceRequest) error {
//	rtm.Mutex.RLock()
//	task, ok := rtm.Tasks[generateKey(request.Operation.Partitions[0], request.Payload.Keys)]
//	rtm.Mutex.RUnlock()
//
//	if !ok {
//		return fmt.Errorf("task not found")
//	}
//
//	task.inputCh <- buildDatum(request)
//	return nil
//}
//
//func (rtm *reduceTaskManager) CloseTask(request *v1.SessionReduceRequest) error {
//	rtm.Mutex.Lock()
//	tasksToBeClosed := make([]*reduceTask, 0, len(request.Operation.Partitions))
//	for _, partition := range request.Operation.Partitions {
//		key := generateKey(partition, request.Payload.Keys)
//		task, ok := rtm.Tasks[key]
//		if ok {
//			tasksToBeClosed = append(tasksToBeClosed, task)
//		}
//		delete(rtm.Tasks, key)
//	}
//	rtm.Mutex.Unlock()
//
//	for _, task := range tasksToBeClosed {
//		close(task.inputCh)
//	}
//
//	return nil
//}
//
//func (rtm *reduceTaskManager) MergeTasks(ctx context.Context, request *v1.SessionReduceRequest) error {
//	rtm.Mutex.Lock()
//	tasks := make([]*reduceTask, 0, len(request.Operation.Partitions))
//	for _, partition := range request.Operation.Partitions {
//		key := generateKey(partition, request.Payload.Keys)
//		task, ok := rtm.Tasks[key]
//		if !ok {
//			rtm.Mutex.Unlock()
//			return fmt.Errorf("task not found")
//		}
//		tasks = append(tasks, task)
//	}
//	rtm.Mutex.Unlock()
//
//	if len(tasks) == 0 {
//		return nil
//	}
//
//	mainTask := tasks[0]
//	aggregators := make([][]byte, 0, len(tasks)-1)
//
//	for _, task := range tasks[1:] {
//		close(task.inputCh)
//		aggregators = append(aggregators, task.reduceStreamer.Aggregator(ctx))
//	}
//
//	for _, aggregator := range aggregators {
//		mainTask.reduceStreamer.MergeAggregator(ctx, aggregator)
//	}
//
//	if request.Payload != nil {
//		mainTask.inputCh <- buildDatum(request)
//	}
//
//	return nil
//}
//
//func (rtm *reduceTaskManager) OutputChannel() <-chan *v1.SessionReduceResponse {
//	return rtm.Output
//}
//
//func (rtm *reduceTaskManager) WaitAll() {
//	rtm.Mutex.RLock()
//	tasks := make([]*reduceTask, 0, len(rtm.Tasks))
//	for _, task := range rtm.Tasks {
//		tasks = append(tasks, task)
//	}
//	rtm.Mutex.RUnlock()
//
//	for _, task := range tasks {
//		<-task.Done
//	}
//}
//
//func generateKey(partition *v1.Partition, keys []string) string {
//	return fmt.Sprintf("%d:%d:%s",
//		partition.GetStart().AsTime().UnixMilli(),
//		partition.GetEnd().AsTime().UnixMilli(),
//		strings.Join(keys, delimiter))
//}
//
//func buildDatum(request *v1.SessionReduceRequest) Datum {
//	return NewHandlerDatum(request.Payload.GetValue(), request.Payload.EventTime.AsTime(), request.Payload.Watermark.AsTime())
//}
