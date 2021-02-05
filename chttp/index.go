package chttp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	dialerTimeout             = 30 * time.Second
	dialerKeepAlive           = 30 * time.Second
	maxConnectionPerHost      = 50
	maxIdleConnections        = 300
	maxIdleConnectionsPerHost = 50

	getMethod  = "GET"
	postMethod = "POST"
)

type Client struct {
	client *http.Client
}

var singletonClient *Client

type (
	Config struct {
		MaxIdleConnection        int
		MaxIdleConnectionPerHost int
		MaxConnectionsPerHost    int
		Timeout                  time.Duration
		DialerTimeout            time.Duration
		DialerKeepAlive          time.Duration
	}
)

func NewSingletonClient(
	config Config,
) *Client {
	if singletonClient != nil {
		return singletonClient
	}
	config.manageConfig()

	transport := http.DefaultTransport.(*http.Transport)
	transport.MaxIdleConns = config.MaxIdleConnection
	transport.MaxIdleConnsPerHost = config.MaxConnectionsPerHost
	transport.MaxConnsPerHost = config.MaxConnectionsPerHost

	dialContext := (&net.Dialer{
		Timeout:   config.DialerTimeout,
		KeepAlive: config.DialerKeepAlive,
	}).DialContext
	transport.DialContext = dialContext

	httpClient := &http.Client{
		Transport: transport,
	}

	if config.Timeout != 0 {
		httpClient.Timeout = config.Timeout
	}

	singletonClient = &Client{
		client: httpClient,
	}
	return singletonClient
}

func (cl Client) ServePost(
	ctx context.Context,
	endpoint string,
	headers map[string]string,
	payload []byte,
	username *string,
	password *string,
) ([]byte, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		postMethod,
		endpoint,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, err
	}
	if username != nil && password != nil {
		request.SetBasicAuth(*username, *password)
	}

	populateHeaders(request, headers)

	resp, err := cl.client.Do(request)
	if err != nil {
		return nil, err
	}

	decodedByte, err := decodeResponse(resp)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return decodedByte, err
}

func (cl Client) ServeGet(
	ctx context.Context,
	endpoint string,
	headers map[string]string,
	queryParams map[string]string,
	username *string,
	password *string,
) ([]byte, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		getMethod,
		endpoint,
		nil,
	)
	if err != nil {
		return nil, err
	}
	populateHeaders(request, headers)

	q := request.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	request.URL.RawQuery = q.Encode()
	if username != nil && password != nil {
		request.SetBasicAuth(*username, *password)
	}

	resp, err := cl.client.Do(request)
	if err != nil {
		return nil, err
	}

	decodedResponse, err := decodeResponse(resp)
	defer resp.Body.Close()

	return decodedResponse, err
}

func (entity *Config) manageConfig() {
	if entity.MaxConnectionsPerHost <= 0 {
		entity.MaxConnectionsPerHost = maxConnectionPerHost
	}
	if entity.MaxIdleConnection <= 0 {
		entity.MaxIdleConnection = maxIdleConnections
	}
	if entity.MaxIdleConnectionPerHost <= 0 {
		entity.MaxIdleConnectionPerHost = maxIdleConnectionsPerHost
	}
	if entity.DialerTimeout == 0 {
		entity.DialerTimeout = dialerTimeout
	}
	if entity.DialerKeepAlive == 0 {
		entity.DialerKeepAlive = dialerKeepAlive
	}
}

func populateHeaders(r *http.Request, headers map[string]string) {
	for key, value := range headers {
		r.Header.Add(key, value)
	}
}

func decodeResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode < 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		_, err := io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("external client cerror")
	}

	if resp.StatusCode >= 500 {
		_, err := io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("external server cerror")
	}
	return nil, errors.New("invalid response")
}

