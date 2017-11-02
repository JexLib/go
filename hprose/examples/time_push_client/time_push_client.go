package main

import (
	"fmt"

	"github.com/JexLib/golang/hprose/rpc"
)

func main() {
	client := rpc.NewTCPClient("tcp4://127.0.0.1:2016/")
	count := 0
	done := make(chan struct{})
	client.Subscribe("time", "", nil, func(data string) {
		count++
		if count > 10 {
			client.Unsubscribe("time")
			done <- struct{}{}
		}
		fmt.Println(data)
	})
	<-done
}
