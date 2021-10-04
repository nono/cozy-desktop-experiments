package client

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/nono/cozy-desktop-experiments/ng/state/remote"
)

type Stack struct {
	Address      string
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
	Client       *http.Client
}

func NewStack(address string) remote.Client {
	return &Stack{
		Address: address,
		Client: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

func (s *Stack) Changes(seq *remote.Seq) (*remote.ChangesResponse, error) {
	return nil, errors.New("Not yet implemented")
}

func (s *Stack) Refresh() error {
	return nil // TODO
}

func (s *Stack) Synchronized() error {
	res, err := s.Post("/settings/synchronized", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (s *Stack) Post(path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+s.AccessToken)
	return s.Client.Do(req)
}
