package worker

import (
	"context"
	"fmt"
	"sync"

	"github.com/gyu-young-park/StoryShift/pkg/log"
)

type WorkerManagable[T any, P any] interface {
	Initialize() error
	Close() error
	Result() chan []T
	Submit(task Task[T, P])
}

type workerPool[T any, P any] struct {
	maxWorker int
	taskQueue chan Task[T, P]
}

type WorkerManager[T any, P any] struct {
	once        sync.Once
	Name        string
	pool        workerPool[T, P]
	resultQueue chan T
	ctx         context.Context
}

func NewWorkerManager[T any, P any](ctx context.Context, name string, maxWorker int) *WorkerManager[T, P] {
	wg := &WorkerManager[T, P]{
		Name: name,
		once: sync.Once{},
		pool: workerPool[T, P]{
			maxWorker: maxWorker,
			taskQueue: make(chan Task[T, P], maxWorker),
		},
		resultQueue: make(chan T, maxWorker),
		ctx:         ctx,
	}
	wg.Initialize()
	return wg
}

func (w *WorkerManager[T, P]) Initialize() {
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

func (w *WorkerManager[T, P]) Close() error {
	w.once.Do(func() {
		close(w.pool.taskQueue)
		close(w.resultQueue)
	})
	return nil
}

func (w *WorkerManager[T, P]) Result() []T {
	ret := []T{}
	for res := range w.resultQueue {
		ret = append(ret, res)
	}
	return ret
}

func (w *WorkerManager[T, P]) Submit(task Task[T, P]) {
	go func() {
		w.pool.taskQueue <- task
	}()
}
