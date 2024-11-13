package nsx_client

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"time"
)

type Client struct {
	HttpClient 	*http.Client
	BaseUrl  		string
	Header			http.Header
	XsrfToken		string
}

const httpTimeoutSeconds = 5

func SetupClient(nsxApi, nsxUsername, nsxPassword string) (*Client, error) {
	baseUrl := "https://" + nsxApi
	CheckConnectivity(baseUrl)

	jar, _ := cookiejar.New(nil)

	httpClient := &http.Client{
		Timeout: 		time.Duration(httpTimeoutSeconds) * time.Second,
		Transport: 	&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Jar: jar,
	}

	postBody := []byte(fmt.Sprintf("j_username=%s&j_password=%s", nsxUsername, nsxPassword))
	req, err := http.NewRequest("POST", baseUrl + "/api/session/create", bytes.NewBuffer(postBody))
	if err != nil {return nil, err}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {return nil, err}
	if resp.StatusCode != 200 {return nil, errors.New("Return code: " + strconv.Itoa(resp.StatusCode))}
	defer resp.Body.Close()

	header := http.Header{
		"accept":       	{"application/json"},
		"x-xsrf-token":		{resp.Header.Get("X-Xsrf-Token")},
	}

	client := &Client{
		HttpClient: httpClient,
		BaseUrl: 		baseUrl,
		Header: 		header,
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
	defer resp.Body.Close()
	
	respBody, err = io.ReadAll(resp.Body)
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