package sync

import "time"

func Start(localDir string) {
	state := NewState()
	ticker := time.Tick(10 * time.Millisecond)
	events := make(chan Event)
	for {
		var ops []Operation
		select {
		case <-ticker:
			state, ops = state.Update(EventTick{})
		case event := <-events:
			state, ops = state.Update(event)
		}
		for _, op := range ops {
			op.Go(events)
		}
	}
}
