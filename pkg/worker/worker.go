package worker

import (
	"context"
	"fmt"
	"sync"

	"github.com/gyu-young-park/StoryShift/pkg/log"
)

type WorkerManagable[P any, R any] interface {
	Close() error
	Submit(task Task[P, R])
	GetResultChan() <-chan R
}

type workerPool[P any, R any] struct {
	maxWorker int
	taskQueue chan Task[P, R]
}

type WorkerManager[P any, R any] struct {
	once        sync.Once
	Name        string
	pool        workerPool[P, R]
	resultQueue chan R
	ctx         context.Context
}

func NewWorkerManager[P any, R any](ctx context.Context, name string, maxWorker int) *WorkerManager[P, R] {
	wg := &WorkerManager[P, R]{}
	wg.setting(ctx, name, maxWorker)
	return wg
}
func (w *WorkerManager[P, R]) setting(ctx context.Context, name string, maxWorker int) {
	w.ctx = ctx
	w.Name = name
	w.once = sync.Once{}
	w.pool = workerPool[P, R]{
		maxWorker: maxWorker,
		taskQueue: make(chan Task[P, R], maxWorker),
	}
	w.resultQueue = make(chan R, maxWorker)
	w.initialize()
}

func (w *WorkerManager[P, R]) initialize() {
	for i := 1; i <= w.pool.maxWorker; i++ {
		go func(workerName string) {
			logger := log.GetLogger()
			for {
				select {
				case task, ok := <-w.pool.taskQueue:
					if !ok {
						return
					}
					logger.Debugf("getting task: %s in %s", task.Name, workerName)
					w.resultQueue <- task.Fn(task.Param)
				case <-w.ctx.Done():
					return
				}
			}
		}(fmt.Sprintf("%s-%d", w.Name, i))
	}
}

func (w *WorkerManager[P, R]) Close() error {
	w.once.Do(func() {
		close(w.pool.taskQueue)
		close(w.resultQueue)
	})
	return nil
}

func (w *WorkerManager[P, R]) GetResultChan() <-chan R {
	return w.resultQueue
}

func (w *WorkerManager[P, R]) Submit(task Task[P, R]) {
	go func() {
		w.pool.taskQueue <- task
	}()
}

// TODO: Split Aggregation function to AggregateDecoreator struct
func (w *WorkerManager[P, R]) AggregateAndClose(paramList []P, taskFunc TaskFn[P, R]) []R {
	return w.aggregateAndClose(paramList, taskFunc)
}

func (w *WorkerManager[P, R]) Aggregate(paramList []P, taskFunc TaskFn[P, R]) []R {
	ret := w.aggregateAndClose(paramList, taskFunc)
	w.setting(w.ctx, w.Name, w.pool.maxWorker)
	return ret
}

func (w *WorkerManager[P, R]) aggregateAndClose(paramList []P, taskFunc TaskFn[P, R]) []R {
	var wg sync.WaitGroup
	for i, param := range paramList {
		wg.Add(1)
		w.Submit(Task[P, R]{
			Name:  fmt.Sprintf("aggregation-%s-%v", w.Name, i),
			Param: param,
			Fn:    taskFunc,
		})
	}

	go func() {
		wg.Wait()
		w.Close()
	}()

	ret := []R{}
	for res := range w.resultQueue {
		wg.Done()
		ret = append(ret, res)
	}

	return ret
}
