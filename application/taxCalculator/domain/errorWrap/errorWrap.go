package errorWrap

type ErrorWrap[T any] struct {
	Value T
	Error error
}
