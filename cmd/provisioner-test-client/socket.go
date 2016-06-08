package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// For some reason there's a weird asymmetry between the http package's server
// and client functionality, the former can easily be used with a unix socket
// file, the latter makes it painful.

// Socketclient acts a helper to make life a little saner.

type SocketClient struct {
	client *http.Client
}

func NewSocketClient(path string) *SocketClient {
	transport := &http.Transport{Dial: func(proto, addr string) (net.Conn, error) {
		return net.Dial("unix", path)
	}}

	return &SocketClient{&http.Client{Transport: transport}}
}

func do(uri string, action func(string) (*http.Response, error)) (ret string, err error) {
	var res *http.Response

	// This is where things get silly - we override the dial so the hostname
	// is ignored.
	dummyUrl := fmt.Sprintf("http://provisioner/%s", uri)

	if res, err = action(dummyUrl); err != nil {
		return
	}

	// Vital, otherwise we leak resources.
	defer res.Body.Close()

	var bytes []byte
	if bytes, err = ioutil.ReadAll(res.Body); err == nil {
		ret = string(bytes)
	}

	return
}

func (c *SocketClient) Get(uri string) (ret string, err error) {
	return do(uri, func(url string) (resp *http.Response, err error) {
		resp, err = c.client.Get(url)
		return
	})
}

func (c *SocketClient) PostJsonString(uri, body string) (ret string, err error) {
	reader := strings.NewReader(body)

	return do(uri, func(url string) (resp *http.Response, err error) {
		resp, err = c.client.Post(url, "text/json", reader)
		return
	})
}
