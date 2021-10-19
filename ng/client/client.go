// Package client provides a remote.Client implementation that makes requests
// to cozy-stack. It also provides a fake implementation that mocks it for
// tests.
package client

import (
	"bytes"
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

// Client is a client for cozy-stack.
type Client struct {
	Address      string
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
	Client       *http.Client
}

// New returns a client for cozy-stack.
func New(address string) remote.Client {
	// TODO use a specific user-agent
	// TODO update the OAuth client on new versions of the client
	return &Client{
		Address: address,
		Client: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

// Register will create an OAuth client on the stack.
func (c *Client) Register() error {
	if c.ClientID != "" {
		return nil
	}
	payload, err := json.Marshal(map[string]interface{}{
		"redirect_uris":    []string{"http://localhost:9000/"},
		"client_name":      "Cozy-Desktop-NG",
		"software_id":      "github.com/nono/cozy-desktop-experiments",
		"software_version": "0.0.1",
		"client_kind":      "desktop",
	})
	if err != nil {
		return err
	}
	res, err := c.NewRequest(http.MethodPost, "/auth/register").
		ContentType("application/json").
		Accept("application/json").
		Body(bytes.NewReader(payload)).
		Do()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		// Flush the body to allow reusing the connection with keepalive
		_, _ = io.Copy(ioutil.Discard, res.Body)
		return fmt.Errorf("invalid status code %d for Register", res.StatusCode)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return err
	}
	c.ClientID, _ = body["client_id"].(string)
	c.ClientSecret, _ = body["client_secret"].(string)
	// TODO registration_access_token
	return nil
}

// Changes is required by the remote.Client interface.
func (c *Client) Changes(seq *remote.Seq) (*remote.ChangesResponse, error) {
	return nil, errors.New("Not yet implemented")
}

// CreateDir is required by the remote.Client interface.
func (c *Client) CreateDir(parentID remote.ID, name string) (*remote.Doc, error) {
	params := url.Values{
		"Type": {"directory"},
		"Name": {name},
	}
	path := fmt.Sprintf("/files/%s?%s", parentID, params.Encode())
	res, err := c.NewRequest(http.MethodPost, path).Do()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		// Flush the body to allow reusing the connection with keepalive
		_, _ = io.Copy(ioutil.Discard, res.Body)
		return nil, fmt.Errorf("invalid status code %d for CreateDir", res.StatusCode)
	}
	return jsonapi.ParseDoc(res.Body)
}

// Trash is required by the remote.Client interface.
func (c *Client) Trash(doc *remote.Doc) (*remote.Doc, error) {
	res, err := c.NewRequest(http.MethodDelete, fmt.Sprintf("/files/%s", doc.ID)).
		IfMatch(string(doc.Rev)).
		Do()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		out, _ := io.ReadAll(res.Body)
		fmt.Printf("out = %s\n", out)
		// Flush the body to allow reusing the connection with keepalive
		_, _ = io.Copy(ioutil.Discard, res.Body)
		return nil, fmt.Errorf("invalid status code %d for Trash", res.StatusCode)
	}
	return jsonapi.ParseDoc(res.Body)
}

// Refresh is required by the remote.Client interface.
func (c *Client) Refresh() error {
	body := strings.NewReader(url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {c.RefreshToken},
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
	}.Encode())
	res, err := c.NewRequest(http.MethodPost, "/auth/access_token").
		ContentType("application/x-www-form-urlencoded").
		Body(body).
		Do()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		// Flush the body to allow reusing the connection with keepalive
		_, _ = io.Copy(ioutil.Discard, res.Body)
		return fmt.Errorf("invalid status code %d for Refresh", res.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return err
	}
	if token, ok := data["refresh_token"].(string); ok && token != "" {
		c.RefreshToken = token
	}
	token, ok := data["access_token"].(string)
	if !ok || token == "" {
		return errors.New("invalid response for Refresh")
	}
	c.AccessToken = token
	return nil
}

// Synchronized is required by the remote.Client interface.
func (c *Client) Synchronized() error {
	res, err := c.NewRequest(http.MethodPost, "/settings/synchronized").Do()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

// NewRequest creates a request representation that can be forged to send an
// HTTP request to the stack.
func (c *Client) NewRequest(verb, path string) *request {
	headers := map[string]string{
		"authorization": "Bearer " + c.AccessToken,
	}
	return &request{
		verb:    verb,
		url:     c.Address + path,
		headers: headers,
		client:  c.Client,
	}
}

type request struct {
	verb    string
	url     string
	headers map[string]string
	body    io.Reader
	client  *http.Client
}

func (r *request) Accept(ctype string) *request {
	return r.addHeader("accept", ctype)
}

func (r *request) ContentType(ctype string) *request {
	return r.addHeader("content-type", ctype)
}

func (r *request) IfMatch(etag string) *request {
	return r.addHeader("if-match", etag)
}

func (r *request) addHeader(key, value string) *request {
	r.headers[key] = value
	return r
}

func (r *request) Body(body io.Reader) *request {
	r.body = body
	return r
}

func (r *request) Do() (*http.Response, error) {
	req, err := http.NewRequest(r.verb, r.url, r.body)
	if err != nil {
		return nil, err
	}
	for k, v := range r.headers {
		req.Header.Add(k, v)
	}
	return r.client.Do(req)
}
