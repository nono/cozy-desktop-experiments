package state

type State struct {
	local *LocalState
}

func Sync(platform Platform) {
	state := &State{
		local: NewLocalState(),
	}
	ops := EventStart{}.Update(state)
	for {
		for _, op := range ops {
			op.Go(platform)
		}
		event := platform.NextEvent()
		ops = event.Update(state)
	}
}
