package chronix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"time"
)

// A SolrClient allows updating documents in Solr.
type SolrClient interface {
	Update(data []map[string]interface{}, commit bool) error
}

// CancelableTransport is like net.Transport but provides
// per-request cancelation functionality.
type CancelableTransport interface {
	http.RoundTripper
	CancelRequest(req *http.Request)
}

// DefaultTransport is used by the solrClient when no explicit transport is provided.
var DefaultTransport CancelableTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 10 * time.Second,
}

type solrClient struct {
	url        *url.URL
	httpClient http.Client
}

// NewSolrClient creates a new Solr client.
func NewSolrClient(url *url.URL, transport CancelableTransport) SolrClient {
	if transport == nil {
		transport = DefaultTransport
	}
	return &solrClient{
		url: url,
		httpClient: http.Client{
			Transport: transport,
		},
	}
}

// Update implements SolrClient.
func (c *solrClient) Update(data []map[string]interface{}, commit bool) error {
	u := c.url
	u.Path = path.Join(c.url.Path, "/update")
	if commit {
		q := u.Query()
		q.Set("commit", "true")
		u.RawQuery = q.Encode()
	}

	buf, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad HTTP response code: %s", resp.Status)
	}
	return nil
}
