package gohttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JexLib/golang/cache"

	"github.com/JexLib/golang/utils"
)

type Client struct {
	client *http.Client
}

func NewClient(mcache ...cache.Cache) *Client {
	mClient := &http.Client{}

	return &Client{
		client: mClient,
	}
}

func (r *Client) SetTimeOut(timeout string) *Client {
	r.client.Timeout = utils.MustParseDuration(timeout)
	return r
}

func (r *Client) Post(url string, data_ptr interface{}) ([]byte, error) {
	data, _ := json.Marshal(data_ptr)
	return r.doRequest(url, "POST", data)
}

func (r *Client) Get(url string, result_ptr ...interface{}) (bytes []byte, err error) {
	if bytes, err = r.doRequest(url, "GET", nil); err == nil && len(result_ptr) > 0 {
		err = json.Unmarshal(bytes, result_ptr[0]) // JSON to Struct
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return
}

func (r *Client) doRequest(url string, method string, data []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if method == "POST" {
		req.Header.Set("Content-Length", (string)(len(data)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if bytes, err := ioutil.ReadAll(resp.Body); err == nil {
		return bytes, nil
	}

	return nil, err
}
