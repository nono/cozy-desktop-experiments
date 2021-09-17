package state

type State struct {
}

func Sync(platform Platform) {
	state := &State{}
	ops := EventStart{}.Update(state)
	for {
		for _, op := range ops {
			op.Go(platform)
		}
		event := platform.NextEvent()
		ops = event.Update(state)
	}
}
