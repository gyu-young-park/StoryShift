package worker

import (
	"context"
	"fmt"
	"sync"

	"github.com/gyu-young-park/StoryShift/pkg/log"
)

type WorkerManagable[P any, R any] interface {
	Initialize() error
	Close() error
	Result() chan []R
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
	wg.Initialize()
	return wg
}

func (w *WorkerManager[P, R]) Initialize() {
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
	w.once.Do(func() {
		close(w.pool.taskQueue)
		close(w.resultQueue)
	})
	return nil
}

func (w *WorkerManager[P, R]) Result() []R {
	ret := []R{}
	for res := range w.resultQueue {
		ret = append(ret, res)
	}
	return ret
}

func (w *WorkerManager[P, R]) Submit(task Task[P, R]) {
	go func() {
		w.pool.taskQueue <- task
	}()
}
