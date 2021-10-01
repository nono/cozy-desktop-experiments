package client

import (
	"errors"

	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
)

type Stack struct {
	Address string
}

func NewStack(address string) remote.Client {
	return &Stack{Address: address}
}

func (s *Stack) Changes(seq *remote.Seq) (*remote.ChangesResponse, error) {
	return nil, errors.New("Not yet implemented")
}
