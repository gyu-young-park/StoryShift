package worker

type Task[T any, P any] struct {
	Name  string
	Param P
	Fn    func(P) T
}
