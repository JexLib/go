// Copyright 2014-2015 Liu Dong <ddliuhb@gmail.com>.
// Licensed under the MIT license.

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JexLib/golang/cache/memory"
	httpclient "github.com/JexLib/golang/gohttpclient"
)

const (
	USERAGENT       = "my awsome httpclient"
	TIMEOUT         = 10
	CONNECT_TIMEOUT = 5
	SERVER          = "https://github.com"
)

func StartTestServer() {
	http.HandleFunc("/now", func(w http.ResponseWriter, r *http.Request) {
		d := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println("收到now请求", d)
		w.Header().Set("Cache-Control", "max-age=3600")
		time.Sleep(time.Second * 15)
		fmt.Fprintf(w, d)
	})
	fmt.Println("listen on port 9090")
	err := http.ListenAndServe("127.0.0.1:9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func TestGet_cache(clientUseCache *httpclient.HttpClient) {

	// get
	res, err := clientUseCache.Get("http://127.0.0.1:9090/now", nil)

	if err != nil {
		log.Println("get failed", err)
		return
	}

	if res.StatusCode != 200 {
		log.Println("Status Code not 200")
	} else {
		str, _ := res.ToString()
		log.Println("返回数据:", res.Header, str)
	}
}

func main() {
	httpclient.Defaults(httpclient.Map{
		"opt_useragent":   USERAGENT,
		"opt_timeout":     TIMEOUT,
		"Accept-Encoding": "gzip,deflate,sdch",
	})

	res, _ := httpclient.
		WithHeader("Accept-Language", "en-us").
		WithCookie(&http.Cookie{
			Name:  "name",
			Value: "github",
		}).
		WithHeader("Referer", "http://163.com").
		Get(SERVER, nil)

	fmt.Println("Cookies:")
	for k, v := range httpclient.CookieValues(SERVER) {
		fmt.Println(k, ":", v)
	}

	fmt.Println("Response:")
	fmt.Println(res.ToString())

	//缓存测试
	mCache := memory.NewMemoryCache(time.Second*5, time.Second*2)
	clientUseCache := httpclient.NewHttpClient()
	clientUseCache.Defaults(
		httpclient.Map{
			"opt_useragent":   USERAGENT,
			"opt_timeout":     TIMEOUT,
			"Accept-Encoding": "gzip,deflate,sdch",
		},
	)
	// clientUseCache.WithOptions(httpclient.Map{
	// 	"opt_useragent":   USERAGENT,
	// 	"opt_timeout":     TIMEOUT,
	// 	"Accept-Encoding": "gzip,deflate,sdch",
	// })

	clientUseCache.WithHeader("Cache-Control", "max-age=3600")
	clientUseCache.WithCache(mCache)
	go StartTestServer()
	for {

		TestGet_cache(clientUseCache)

		time.Sleep(time.Second)
	}

}
