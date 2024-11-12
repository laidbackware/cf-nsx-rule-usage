package nsx_client

import (

	b64 "encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"crypto/tls"
	"net/http"
	"time"
)

type Client struct {
	HttpClient *http.Client
	BaseUrl  	string
	Header		http.Header
}

const (
	httpTimeoutSeconds = 3
)

var (
	ErrorConnectionFailure = errors.New("login: server did not return 200 ok")
)

// username, password
func SetupClient(nsxApi, nsxUsername, nsxPassword string) (*Client, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	CheckConnectivity(nsxApi)

	httpClient := &http.Client{
		Timeout: time.Duration(httpTimeoutSeconds) * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	authHeader := b64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", nsxUsername, nsxPassword)))

	header := http.Header{
		"accept":        {"application/json"},
		"authorization": {"Basic " + authHeader},
	}

	client := &Client{
		HttpClient: httpClient,
		BaseUrl: "https://" + nsxApi,
		Header: header,
	}

	return client, nil
}

func (c *Client) makeGetRequest(endpoint string) ([]byte, error) {
	var respBody []byte
	if c.BaseUrl == "" || c.Header == nil {
		return respBody, errors.New("Client not initialized")
	}

	req, err := http.NewRequest("GET", c.BaseUrl + endpoint, nil)
	if err != nil {return respBody, err}
	req.Header = c.Header

	resp, err := c.HttpClient.Do(req)
	if err != nil {return respBody, err}
	if resp.StatusCode != 200 {return respBody, errors.New("Return code: " + strconv.Itoa(resp.StatusCode))}
	
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {return respBody, err}

	return respBody, nil
}

func CheckConnectivity(api string) (err error) {
	httpClient := &http.Client{
		Timeout: time.Duration(httpTimeoutSeconds) * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	var res *http.Response
	res, err = httpClient.Get(fmt.Sprintf("https://%s", api))
	if err != nil {return}

	if res.StatusCode != 200 {
		err = ErrorConnectionFailure
	}
	return
}