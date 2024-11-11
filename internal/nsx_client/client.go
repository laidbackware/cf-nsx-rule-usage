package nsx_client

import (

	b64 "encoding/base64"
	// "encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	// "bytes"
	"crypto/tls"
	// "io"
	"net/http"
	// "net/url"
	// "strings"
	"time"

	// "github.com/andybalholm/cascadia"
	// "github.com/orirawlings/persistent-cookiejar"
	// "golang.org/x/net/html"
)

type Client struct {
	HttpClient *http.Client
	XsrfToken  string
}

const (
	httpTimeoutSeconds = 3
)

var (
	ErrorConnectionFailure = errors.New("login: server did not return 200 ok")
	
	api 			string
	username 	string
	password 	string
)

// username, password
func SetupClient(nsxApi, nsxUsername, nsxPassword string) (client *Client, err error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	api, username, password = nsxApi, nsxUsername, nsxPassword
	CheckConnectivity(api)
	
	httpClient := &http.Client{Timeout: time.Duration(httpTimeoutSeconds) * time.Second}
	//TODO setup auth session
	
	return &Client{HttpClient: httpClient}, nil

}

func (c *Client) makeRequest(url string) ([]byte, error) {
	var respBody []byte
	if api == "" || username == "" || password == "" {
		return respBody, errors.New("Client not initialized")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {return respBody, err}

		// Add request headers
	authstr := b64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	req.Header = http.Header{
		"accept":        {"application/json"},
		"authorization": {"Basic " + authstr},
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {return respBody, err}
	if resp.StatusCode != 200 {return respBody, errors.New("Return code: " + strconv.Itoa(resp.StatusCode))}
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {return respBody, err}

	return respBody, nil
}

func CheckConnectivity(api string) (err error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	httpClient := &http.Client{Timeout: time.Duration(httpTimeoutSeconds) * time.Second}
	var res *http.Response
	res, err = httpClient.Get(fmt.Sprintf("https://%s", api))
	if err != nil {return}

	if res.StatusCode != 200 {
		err = ErrorConnectionFailure
	}
	return
}