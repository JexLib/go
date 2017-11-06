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

type HttpClient struct {
	client *http.Client
}

func NewHttpClient(timeout string, mcache ...cache.Cache) *HttpClient {
	timeoutIntv := utils.MustParseDuration(timeout)
	mClient := &http.Client{
		Timeout: timeoutIntv,
	}

	if len(mcache) > 0 {
		mClient.Transport = cache.NewHttpCacheTransport(mcache[0])
	}

	return &HttpClient{
		client: mClient,
	}
}

func (r *HttpClient) Post(url string, data_ptr interface{}) ([]byte, error) {
	data, _ := json.Marshal(data_ptr)
	return r.doHttpRequest(url, "POST", data)
}

func (r *HttpClient) Get(url string, result_ptr ...interface{}) (bytes []byte, err error) {
	if bytes, err = r.doHttpRequest(url, "GET", nil); err == nil && len(result_ptr) > 0 {
		err = json.Unmarshal(bytes, result_ptr[0]) // JSON to Struct
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return
}

func (r *HttpClient) doHttpRequest(url string, method string, data []byte) ([]byte, error) {
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
