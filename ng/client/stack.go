package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
	params := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {s.RefreshToken},
		"client_id":     {s.ClientID},
		"client_secret": {s.ClientSecret},
	}
	body := strings.NewReader(params.Encode())
	res, err := s.Post("/auth/access_token", "application/x-www-form-urlencoded", body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Invalid status code %d for refresh_token", res.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return err
	}
	if token, ok := data["refresh_token"].(string); ok && token != "" {
		s.AccessToken = token
	}
	token, ok := data["access_token"].(string)
	if !ok || token == "" {
		return errors.New("Invalid response for refresh_token")
	}
	s.AccessToken = token
	return nil
}

func (s *Stack) Synchronized() error {
	res, err := s.Post("/settings/synchronized", "", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (s *Stack) Post(path string, ctype string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	if ctype != "" {
		req.Header.Add("content-type", ctype)
	}
	req.Header.Add("authorization", "Bearer "+s.AccessToken)
	return s.Client.Do(req)
}
