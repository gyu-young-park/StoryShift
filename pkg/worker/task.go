package worker

type Task[T any] struct {
	Name string
	Fn   func() T
}
