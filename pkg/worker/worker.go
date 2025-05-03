package worker

import (
	"context"
	"fmt"

	"github.com/gyu-young-park/StoryShift/pkg/log"
)

type WorkerManagable[T any] interface {
	Initialize() error
	Close() error
	Result() chan T
	Submit(task Task[T])
}

type workerPool[T any] struct {
	maxWorker int
	taskQueue chan Task[T]
}

type WorkerManager[T any] struct {
	Name        string
	pool        workerPool[T]
	resultQueue chan T
	ctx         context.Context
}

func NewWorkerManager[T any](ctx context.Context, name string, maxWorker int) *WorkerManager[T] {
	return &WorkerManager[T]{
		Name: name,
		pool: workerPool[T]{
			maxWorker: maxWorker,
			taskQueue: make(chan Task[T], maxWorker),
		},
		resultQueue: make(chan T, maxWorker),
		ctx:         ctx,
	}
}

func (w *WorkerManager[T]) Initialize() {
	for i := 1; i <= w.pool.maxWorker; i++ {
		go func(workerName string) {
			logger := log.GetLogger()
			for {
				select {
				case task := <-w.pool.taskQueue:
					logger.Debugf("getting task: %s in %s", task.Name, workerName)
					w.resultQueue <- task.Fn()
				case <-w.ctx.Done():
					return
				}
			}
		}(fmt.Sprintf("%s-%d", w.Name, i))
	}
}

func (w *WorkerManager[T]) Close() error {
	close(w.pool.taskQueue)
	close(w.resultQueue)
	return nil
}

func (w *WorkerManager[T]) Result() chan T {
	return w.resultQueue
}

func (w *WorkerManager[T]) Submit(task Task[T]) {
	w.pool.taskQueue <- task
}
