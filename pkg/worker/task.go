package worker

type TaskFn[P any, R any] func(P) R

type Task[P any, R any] struct {
	Name  string
	Param P
	Fn    func(P) R
}
