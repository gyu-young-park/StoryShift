package worker

type Task[P any, R any] struct {
	Name  string
	Param P
	Fn    func(P) R
}
