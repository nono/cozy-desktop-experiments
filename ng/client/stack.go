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
	res, err := s.NewRequest(http.MethodPost, path).Do()
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

func (s *Stack) Trash(doc *remote.Doc) (*remote.Doc, error) {
	res, err := s.NewRequest(http.MethodDelete, fmt.Sprintf("/files/%s", doc.ID)).
		AddHeader("if-match", string(doc.Rev)).
		Do()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		// Flush the body to allow reusing the connection with keepalive
		_, _ = io.Copy(ioutil.Discard, res.Body)
		return nil, fmt.Errorf("invalid status code %d for Trash", res.StatusCode)
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
	res, err := s.NewRequest(http.MethodPost, "/auth/access_token").
		ContentType("application/x-www-form-urlencoded").
		Body(body).
		Do()
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
	res, err := s.NewRequest(http.MethodPost, "/settings/synchronized").Do()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (s *Stack) NewRequest(verb, path string) *request {
	headers := map[string]string{
		"authorization": "Bearer " + s.AccessToken,
	}
	return &request{
		verb:    verb,
		path:    path,
		headers: headers,
		client:  s.Client,
	}
}

type request struct {
	verb    string
	path    string
	headers map[string]string
	body    io.Reader
	client  *http.Client
}

func (r *request) ContentType(ctype string) *request {
	r.headers["content-type"] = ctype
	return r
}

func (r *request) AddHeader(key, value string) *request {
	r.headers[key] = value
	return r
}

func (r *request) Body(body io.Reader) *request {
	r.body = body
	return r
}

func (r *request) Do() (*http.Response, error) {
	req, err := http.NewRequest(r.verb, r.path, r.body)
	if err != nil {
		return nil, err
	}
	for k, v := range r.headers {
		req.Header.Add(k, v)
	}
	return r.client.Do(req)
}
