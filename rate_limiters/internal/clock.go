package internal

type Clock interface {
	Now() int64
}
