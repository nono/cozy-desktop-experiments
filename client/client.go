// Package client provides a remote.Client implementation that makes requests
// to cozy-stack. It also provides a fake implementation that mocks it for
// tests.
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/nono/cozy-desktop-experiments/client/jsonapi"
	"github.com/nono/cozy-desktop-experiments/state/remote"
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
func New(address string) *Client {
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

type changesResponse struct {
	LastSeq remote.Seq      `json:"last_seq"`
	Pending int             `json:"pending"`
	Results []changesResult `json:"results"`
}

type changesResult struct {
	ID      string                 `json:"id"`
	Deleted bool                   `json:"deleted"`
	Doc     map[string]interface{} `json:"doc"`
}

// Changes is required by the remote.Client interface.
//
// Note: the design docs are ignored
func (c *Client) Changes(seq *remote.Seq, limit int, skipTrashed bool) (*remote.ChangesResponse, error) {
	params := url.Values{
		"include_docs": {"true"},
		"fields":       {"_rev,name,type,dir_id"},
		"limit":        {fmt.Sprintf("%d", limit)},
	}
	if seq != nil {
		params.Add("since", string(*seq))
	}
	if skipTrashed {
		params.Add("skip_deleted", "true")
		params.Add("skip_trashed", "true")
	}
	path := fmt.Sprintf("/files/_changes?%s", params.Encode())
	res, err := c.NewRequest(http.MethodGet, path).Do()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		// Flush the body to allow reusing the connection with keepalive
		_, _ = io.Copy(ioutil.Discard, res.Body)
		return nil, fmt.Errorf("invalid status code %d for Changes", res.StatusCode)
	}

	var changes changesResponse
	if err := json.NewDecoder(res.Body).Decode(&changes); err != nil {
		return nil, err
	}
	docs := make([]*remote.ChangedDoc, 0, len(changes.Results))
	for _, result := range changes.Results {
		doc := &remote.Doc{
			ID: remote.ID(result.ID),
		}
		if rev, ok := result.Doc["_rev"].(string); ok {
			doc.Rev = remote.Rev(rev)
		}
		if name, ok := result.Doc["name"].(string); ok {
			doc.Name = name
		}
		if typ, ok := result.Doc["type"].(string); ok {
			doc.Type = jsonapi.ConvertType(typ)
		}
		if dirID, ok := result.Doc["dir_id"].(string); ok {
			doc.DirID = remote.ID(dirID)
		}
		changed := &remote.ChangedDoc{Doc: doc, Deleted: result.Deleted}
		docs = append(docs, changed)
	}
	return &remote.ChangesResponse{
		Docs:    docs,
		Seq:     changes.LastSeq,
		Pending: changes.Pending,
	}, nil
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

// DocsByID returns a map of id -> doc (for testing purpose).
func (c *Client) DocsByID() map[remote.ID]*remote.Doc {
	params := url.Values{
		"include_docs": {"true"},
		"limit":        {"10000"},
	}
	path := fmt.Sprintf("/data/io.cozy.files/_all_docs?%s", params.Encode())
	res, err := c.NewRequest(http.MethodGet, path).Do()
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		panic(fmt.Errorf("invalid status code %d for DocsByID", res.StatusCode))
	}
	var list allDocsResponse
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		panic(err)
	}
	byID := map[remote.ID]*remote.Doc{}
	for _, row := range list.Rows {
		id := remote.ID(row.Doc.ID)
		if id.IsDesignDoc() {
			continue
		}
		doc := &remote.Doc{
			ID:    id,
			Rev:   remote.Rev(row.Doc.Rev),
			Type:  jsonapi.ConvertType(row.Doc.Type),
			Name:  row.Doc.Name,
			DirID: remote.ID(row.Doc.DirID),
		}
		byID[doc.ID] = doc
	}
	return byID
}

type allDocsResponse struct {
	Rows []allDocsRow `json:"rows"`
}

type allDocsRow struct {
	Doc struct {
		ID    string `json:"_id"`
		Rev   string `json:"_rev"`
		Type  string `json:"type"`
		Name  string `json:"name"`
		DirID string `json:"dir_id"`
	} `json:"doc"`
}

// NewRequest creates a request representation that can be forged to send an
// HTTP request to the stack.
func (c *Client) NewRequest(verb, path string) *request {
	var hostID uint32
	if host, err := os.Hostname(); err == nil {
		hostID = crc32.ChecksumIEEE([]byte(host))
	}
	headers := map[string]string{
		"authorization": "Bearer " + c.AccessToken,
		"user-agent":    fmt.Sprintf("Cozy-Desktop-NG-%s-%0x", runtime.GOOS, hostID),
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

var _ remote.Client = &Client{}
