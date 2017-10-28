package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//获取远程http json数据
func HttpGetJson(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = json.Unmarshal(body, v) // JSON to Struct
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
