package queue

/**
  排队管理
*/

import (
	"time"

	"log"
)

type Queues struct {
	MaxGoroutine chan int
	items        []func() error
	breaktag     bool
	onCallback   OnFinishTaskCallback
}

type OnFinishTaskCallback func()

func NewQueues(maxGoroutineCount int, onCallback ...OnFinishTaskCallback) *Queues {
	t := new(Queues)
	t.MaxGoroutine = make(chan int, maxGoroutineCount)
	if len(onCallback) > 0 {
		t.onCallback = onCallback[0]
	}
	return t
}

func (t *Queues) Start() {
	log.Println("start")
	t.breaktag = false
	go func() {
		for {
			for len(t.items) > 0 {

				if t.breaktag {
					break
				}
				v := t.items[0]
				t.MaxGoroutine <- 1
				go func(r func() error) {
					//fmt.Println("start task:",k)
					log.Println("start task")
					err := r()
					if t.onCallback != nil {
						t.onCallback()
					}
					log.Println("finish task:", err)
					<-t.MaxGoroutine
				}(v)
				//delete(t.tasks,k)
				t.items = append(t.items[:0], t.items[0+1:]...)
			}
			if t.breaktag {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	
}

func (t *Queues) AddItem(item func() error) {
	t.items = append(t.items, item)
}

func (t *Queues) Stop() {
	t.breaktag = true
}
