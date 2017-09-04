package statistics

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var (
	statistics *Statistics
	baseName   = "statistics"
	stoday     = ""
	update     = false
)

type (
	Statistics struct {
		tags map[string]interface{}
		mu   sync.RWMutex
	}

	StatisticsNote map[string]float64
)

func GetStatistics() interface{} {
	if statistics == nil {
		statistics = &Statistics{
			tags: make(map[string]interface{}),
		}
		statistics.read()
		go statistics.write()
	}
	return statistics.tags
}

func (st *Statistics) getNote(tags ...string) StatisticsNote {
	st.mu.Lock()
	defer st.mu.Unlock()

	tmSt := st.tags
	for i := 0; i < len(tags)-1; i++ {
		if tmSt[tags[i]] == nil {
			tmSt[tags[i]] = make(map[string]interface{})
		}
		tmSt = tmSt[tags[i]].(map[string]interface{})
	}

	if tmSt[tags[len(tags)-1]] == nil {
		tmSt[tags[len(tags)-1]] = make(StatisticsNote)
	} else {
		switch tmSt[tags[len(tags)-1]].(type) {
		case map[string]interface{}:
			tmNote := make(StatisticsNote)
			for sk, sv := range tmSt[tags[len(tags)-1]].(map[string]interface{}) {
				tmNote[sk] = (sv).(float64)
			}
			tmSt[tags[len(tags)-1]] = tmNote
		}
	}
	return tmSt[tags[len(tags)-1]].(StatisticsNote)
}

func NewStatistics(tags ...string) StatisticsNote {
	if statistics == nil {
		statistics = &Statistics{
			tags: make(map[string]interface{}),
		}
		statistics.read()
		go statistics.write()
	}

	return statistics.getNote(tags...)
}

func (note StatisticsNote) Add(key string) {
	statistics.mu.Lock()
	defer statistics.mu.Unlock()
	note[key]++
	update = true
}

func (st *Statistics) Clean() {
	st.mu.Lock()
	defer st.mu.Unlock()
	clear(st.tags)
}

func clear(m interface{}) {
	for k, v := range m.(map[string]interface{}) {
		switch v.(type) {
		case float64:
			m.(map[string]interface{})[k] = float64(0)
		case StatisticsNote:
			for k1, _ := range v.(StatisticsNote) {
				v.(StatisticsNote)[k1] = float64(0)
			}
		default:
			clear(v)
		}

	}
}

func (st *Statistics) write() {
	//30秒保存一次
	timer := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-timer.C:
			if update {
				saveData, _ := json.Marshal(&st.tags)
				//json格式化
				var out bytes.Buffer
				json.Indent(&out, saveData, "", "\t")

				filename := baseName + stoday
				file, _ := os.Create(filename)
				out.WriteTo(file)
				file.Close()
				update = false

				if stoday != time.Now().Format(".20060102") {
					stoday = time.Now().Format(".20060102")
					st.Clean()
					update = true
				}
			}
		}
	}

}

func (st *Statistics) read() error {
	stoday = time.Now().Format(".20060102")
	filename := baseName + stoday
	//	jsondata := make(map[string]interface{})
	if bts, err := ioutil.ReadFile(filename); err == nil {
		return json.Unmarshal(bts, &st.tags)
	}

	return nil
}
