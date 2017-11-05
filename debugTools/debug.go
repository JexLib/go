package debugTools

import (
	"fmt"
	"time"
)

/**
 计算函数/过程运行时间成本
 用法：
  func test() {
		defer TimeCost(time.Now())
		.......
  }
**/
func TimeCost(start time.Time, label ...string) {
	terminal := time.Since(start)
	if len(label) == 0 {
		label = append(label, "TimeCost:")
	}
	fmt.Println(label, terminal)
}
