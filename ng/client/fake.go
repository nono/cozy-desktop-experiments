package client

import (
	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
)

type Fake struct {
	Address           string
	SynchronizedCount int
}

func NewFake(address string) remote.Client {
	return &Fake{
		Address:           address,
		SynchronizedCount: 0,
	}
}

func (f *Fake) Changes(seq *remote.Seq) (*remote.ChangesResponse, error) {
	docs := []*remote.Doc{}
	lastSeq := remote.Seq("0-g1AAAAIDeJyN0EEOgjAQQNGJmKgLz6BHKNBCWclNtMOUIMF2oa71JnoTvYneBEtYAAuTbmaSSf5bTAMAyyogWBtrLOnc2MqeL407zxTgpm3bugoUnNxhQcg06ZhgdTWky6PR9CfFrZu4m9QskkKh8qnzrt5P6pQkQ8596kNX3yZ1KIizED1qM3cT7m454DEIXKQ8Kf2FZy-8Rt9jShSUeQvvXvgMQqy5FIXwFr69MPpDFiUxFdFYqH_ptZpY")
	return &remote.ChangesResponse{Docs: docs, Seq: lastSeq, Pending: 0}, nil
}

func (f *Fake) Refresh() error {
	return nil
}

func (f *Fake) Synchronized() error {
	f.SynchronizedCount++
	return nil
}
