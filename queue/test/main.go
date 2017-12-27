package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func main() {
	// urls := []string{
	// 	"http://stock.591hx.com/images/hnimg/201412/03/64/13418266510200941552.jpg",
	// 	"http://stock.591hx.com/images/hnimg/201412/03/7/10886141709285175583.jpg",
	// 	"http://stock.591hx.com/images/hnimg/201412/03/0/5391574706574741364.jpg",
	// 	"http://stock.591hx.com/images/hnimg/201412/03/1/5405780767459854941.jpg",
	// 	"http://stock.591hx.com/images/hnimg/201412/03/16/13722698761317688276.jpg",
	// 	"http://stock.591hx.com/images/hnimg/201412/03/7/16853951343108680551.jpg",
	// 	"http://stock.591hx.com/images/hnimg/201412/03/50/17680852843413447062.jpg",
	// 	"http://stock.591hx.com/images/hnimg/201412/03/34/14366548421421579970.jpg",
	// 	"http://stock.591hx.com/images/hnimg/201412/03/4/17141924098089490820.jpg",
	// }
	// var q *queue.Queues
	// q = queue.NewQueues(2, func() {
	// 	log.Println("end:", len(q.MaxGoroutine))
	// })
	// q.Start()
	// for _, v := range urls {
	// 	if (2 - len(q.MaxGoroutine)) > 0 {
	// 		log.Println("add:", v)
	// 		q.AddItem(func() error {
	// 			GetUrl(v)
	// 			time.Sleep(2000)
	// 			return nil
	// 		})
	// 	}
	// }

	// select {}
	Queue(3)
	// var zeroHash = regexp.MustCompile("^0?x?0+$")
	// var addressPattern = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	// fmt.Println("zeroHash:", zeroHash.MatchString("0x0063Bbe659a157bD7e9f0182870e9c6E235f2A55"))
	// fmt.Println("math:", addressPattern.MatchString("0x0063Bbe659a157bD7e9f0182870e9c6E235f2A55"))

}

func GetUrl(url string) []byte {
	ret, err := http.Get(url)
	if err != nil {
		log.Println(url)
		status := map[string]string{}
		status["status"] = "400"
		status["url"] = url
		panic(status)
	}
	body := ret.Body
	data, _ := ioutil.ReadAll(body)
	return data
}

func Queue(max int) {
	urls := []string{
		"http://autoload.bank.ecitic.com/ebank/perbank/helpmate/HelpmateSetup.exe",
		"https://codeload.github.com/MakeGo/error/zip/master",
		"https://github-production-release-asset-2e65be.s3.amazonaws.com/3402186/bb47f4a2-3fac-11e6-9e71-9a4261699bd5?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20171108%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20171108T113119Z&X-Amz-Expires=300&X-Amz-Signature=348ba8910285d33d164c551fdc4042feb48910cec363033b170aa31e1b88cd5a&X-Amz-SignedHeaders=host&actor_id=1274558&response-content-disposition=attachment%3B%20filename%3DRedis-x64-3.2.100.zip&response-content-type=application%2Foctet-stream",
		"https://codeload.github.com/shafarizkyf/tokobiru/zip/2d8d9611f380723fcde366baa55ee28e3473fa5a",
		"https://download.microsoft.com/download/6/7/D/67DC6A4A-0EF9-4385-A1A0-45A19728BCC4/MicrosoftCameraCodecPack-x64.msi",
		// "http://www2.isl.co.jp/SILKYPIX/p/download/files/SILKYPIX4SE44102.exe",
		"http://down.cssmoban.com/cssthemes3/mbts_15_kit.zip",
		"http://stock.591hx.com/images/hnimg/201412/03/34/14366548421421579970.jpg",
		"http://stock.591hx.com/images/hnimg/201412/03/4/17141924098089490820.jpg",
	}

	maxGoroutine := make(chan int, max)
	wg := sync.WaitGroup{}
	log.Println("ccccc:", len(maxGoroutine))
	for _, v := range urls {
		//	fmt.Println(v)
		wg.Add(1)
		maxGoroutine <- 1
		go func(u string) {

			defer func() {
				<-maxGoroutine
				wg.Done()
				log.Println("end:", len(maxGoroutine), u)
			}()
			log.Println("add:", len(maxGoroutine), u)
			GetUrl(u)

		}(v)
	}
	log.Println("count2:", len(maxGoroutine))
	log.Println("finnn:")

	// for len(maxGoroutine) > 0 {
	// 	log.Println("finnn:")
	// 	time.Sleep(time.Second)
	// }
	wg.Wait()
	log.Println("finnn11111")
}
