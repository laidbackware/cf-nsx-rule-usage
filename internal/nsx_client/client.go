package nsx_client

import (
	b64 "encoding/base64"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	HttpClient 	*http.Client
	BaseUrl  		string
	Header			http.Header
	XsrfToken		string
	log 				Logger
}

const httpTimeoutSeconds = 30

func SetupClient(nsxApi, nsxUsername, nsxPassword string, skipVerify, debug bool, log Logger) (*Client, error) {
	baseUrl := "https://" + nsxApi
	CheckConnectivity(baseUrl)


	httpClient := &http.Client{
		Timeout: 		time.Duration(httpTimeoutSeconds) * time.Second,
		Transport: 	&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	
	authstr := b64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", nsxUsername, nsxPassword)))
	
	header := http.Header{
		"accept":       	{"application/json"},
		"authorization": 	{"Basic " + authstr},
	}

	// Test user credebtials
	req, err := http.NewRequest("GET", baseUrl + "/api/v1/cluster-manager/status", nil)
	if err != nil {return nil, err}
	req.Header = header
	resp, err := httpClient.Do(req)
	if err != nil {return nil, err}
	if resp.StatusCode == 403 {
		return nil, errors.New("invalid username or password")
	}
	if resp.StatusCode != 200 {return nil, errors.New("Return code: " + strconv.Itoa(resp.StatusCode))}
	defer resp.Body.Close()

	
	client := &Client{
		HttpClient: httpClient,
		BaseUrl: 		baseUrl,
		Header: 		header,
		log:				log,
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
	defer resp.Body.Close()
	respBody, err = io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return respBody, errors.New(
			"Return code: " + strconv.Itoa(resp.StatusCode) + 
			"\nfrom endpoint: " + endpoint +
			"\nwith response: " + string(respBody[:]))
	}
	
	if err != nil {return respBody, err}

	return respBody, nil
}

func CheckConnectivity(api string) (err error) {
	httpClient := &http.Client{
		Timeout: 		time.Duration(httpTimeoutSeconds) * time.Second,
		Transport: 	&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	var res *http.Response
	res, err = httpClient.Get(api)
	if err != nil {return}

	if res.StatusCode != 200 {
		err = fmt.Errorf("login: server did not return 200 ok returned: %s", strconv.Itoa(res.StatusCode))
	}
	return
}