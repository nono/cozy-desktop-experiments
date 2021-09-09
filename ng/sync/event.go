package sync

type Event interface {
	EventName() string
}

type EventTick struct{}

func (e EventTick) EventName() string { return "tick" }
