package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nono/cozy-desktop-experiments/ng/client/jsonapi"
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
	// TODO use a specific user-agent
	return &Stack{
		Address: address,
		Client: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

// TODO update the OAuth client on new versions of the client

func (s *Stack) Changes(seq *remote.Seq) (*remote.ChangesResponse, error) {
	// TODO add an option to skip deleted docs
	return nil, errors.New("Not yet implemented")
}

func (s *Stack) CreateDir(parentID remote.ID, name string) (*remote.Doc, error) {
	params := url.Values{
		"Type": {"directory"},
		"Name": {name},
	}
	path := fmt.Sprintf("/files/%s?%s", parentID, params.Encode())
	res, err := s.post(path, "", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		// Flush the body to allow reusing the connection with keepalive
		_, _ = io.Copy(ioutil.Discard, res.Body)
		return nil, fmt.Errorf("invalid status code %d for CreateDir", res.StatusCode)
	}
	return jsonapi.ParseDoc(res.Body)
}

func (s *Stack) Refresh() error {
	params := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {s.RefreshToken},
		"client_id":     {s.ClientID},
		"client_secret": {s.ClientSecret},
	}
	body := strings.NewReader(params.Encode())
	res, err := s.post("/auth/access_token", "application/x-www-form-urlencoded", body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		// Flush the body to allow reusing the connection with keepalive
		_, _ = io.Copy(ioutil.Discard, res.Body)
		return fmt.Errorf("invalid status code %d for Refresh", res.StatusCode)
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
		return errors.New("invalid response for Refresh")
	}
	s.AccessToken = token
	return nil
}

func (s *Stack) Synchronized() error {
	res, err := s.post("/settings/synchronized", "", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (s *Stack) post(path string, ctype string, body io.Reader) (*http.Response, error) {
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
