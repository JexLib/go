package memory

import "testing"
import "time"
import "fmt"

func TestMemCache(t *testing.T) {
	c := NewMemoryCache(time.Second*5, time.Second*2)
	c.Set("111", 1111)
	c.Set("222", 222, time.Second*15)

	for {
		fmt.Println("111", c.Get("111"))
		fmt.Println("222", c.Get("222"))
		fmt.Println("count:", c.Count())
		time.Sleep(time.Second)
	}
}
