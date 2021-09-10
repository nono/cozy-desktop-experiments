package sync

func Start(platform *Platform) {
	state, ops := EventStart{}.Update(&State{})
	for {
		for _, op := range ops {
			op.Go(platform)
		}
		event := <-platform.Events
		state, ops = event.Update(state)
	}
}
