package sync

type Operation interface {
	OperationName() string
	Go(events chan Event)
}
