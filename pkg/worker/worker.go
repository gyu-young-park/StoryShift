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
	wg := &WorkerManager[P, R]{
		Name: name,
		once: sync.Once{},
		pool: workerPool[P, R]{
			maxWorker: maxWorker,
			taskQueue: make(chan Task[P, R], maxWorker),
		},
		resultQueue: make(chan R, maxWorker),
		ctx:         ctx,
	}
	wg.initialize()
	return wg
}

func (w *WorkerManager[P, R]) initialize() {
	for i := 1; i <= w.pool.maxWorker; i++ {
		go func(workerName string) {
			logger := log.GetLogger()
			for {
				select {
				case task := <-w.pool.taskQueue:
					logger.Debugf("getting task: %s in %s", task.Name, workerName)
					w.resultQueue <- task.Fn(task.Param)
				case <-w.ctx.Done():
					w.Close()
					return
				}
			}
		}(fmt.Sprintf("%s-%d", w.Name, i))
	}
}

func (w *WorkerManager[P, R]) Close() error {
	if w.pool.taskQueue != nil {
		close(w.pool.taskQueue)
	}
	if w.resultQueue != nil {
		close(w.resultQueue)
	}
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

func (w *WorkerManager[P, R]) Aggregate(cancel context.CancelFunc, paramList []P, taskFunc TaskFn[P, R]) []R {
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
		defer cancel()
	}()

	ret := []R{}
	for res := range w.resultQueue {
		wg.Done()
		ret = append(ret, res)
	}

	return ret
}

func (w *WorkerManager[P, R]) Reload(ctx context.Context) {
	w.ctx = ctx
	w.pool.taskQueue = make(chan Task[P, R], w.pool.maxWorker)
	w.resultQueue = make(chan R, w.pool.maxWorker)
	w.initialize()
}
