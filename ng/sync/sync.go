package sync

func Start(platform *Platform) {
	state := &State{}
	ops := EventStart{}.Update(state)
	for {
		for _, op := range ops {
			op.Go(platform)
		}
		event := <-platform.Events
		ops = event.Update(state)
	}
}
