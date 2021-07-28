package onfleet

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const defaultBaseURL = "https://onfleet.com/api/v2/"

type Client struct {
	baseURL *url.URL
	client  *http.Client

	// Onfleet resources endpoints
	Organization *OrganizationService
	Workers      *WorkersService
	Hubs         *HubsService
	Teams        *TeamsService
}

type service struct {
	client *Client
}

func (c *Client) WithBaseURL(baseURL string) {
	url, _ := url.Parse(baseURL)
	c.baseURL = url
}

func NewClient(apiKey string) *Client {
	url, _ := url.Parse(defaultBaseURL)
	transport := basicAuthTransport{Username: apiKey}

	client := Client{
		baseURL: url,
		client:  transport.Client(),
	}
	svc := service{client: &client}

	client.Organization = (*OrganizationService)(&svc)
	client.Workers = (*WorkersService)(&svc)
	client.Hubs = (*HubsService)(&svc)
	client.Teams = (*TeamsService)(&svc)

	return &client
}

func (c *Client) NewRequest(method string, path string, body interface{}) (*http.Request, error) {
	url, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		data, err := ioutil.ReadAll(resp.Body)
		if err == nil && data != nil {
			return fmt.Errorf("error %v - %v - %v", resp.StatusCode, req.URL.Path, string(data))
		}

		return fmt.Errorf("error %v - %v", resp.StatusCode, req.URL.Path)
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	return err
}

func Float64(f float64) *float64 { return &f }

type basicAuthTransport struct {
	Username string
	Password string
}

func (bat basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
			bat.Username, bat.Password)))))
	return http.DefaultTransport.RoundTrip(req)
}

func (bat *basicAuthTransport) Client() *http.Client {
	return &http.Client{Transport: bat}
}
