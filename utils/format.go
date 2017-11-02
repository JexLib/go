package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//千位分隔符格式化
func Format_splitK(num string) string {
	var temp = num
	i := strings.Index(temp, ",")
	if i < 0 {
		i = strings.Index(temp, ".")
	}
	if i < 0 {
		i = len(temp)
	}
	if i > 3 {
		temp = temp[0:i-3] + "," + temp[i-3:]
		// fmt.Println("sppk:", num, temp)
		return Format_splitK(temp)
	} else {
		return temp
	}
}

//算力格式化
func FormatHashrate(hashrate float64) string {
	units := []string{"H", "KH", "MH", "GH", "TH", "PH"}
	i := 0
	h := hashrate
	for h > 1000 {
		h = h / 1000
		i++
	}
	return fmt.Sprintf("%6.2f", h) + " " + units[i]
}

//矿池难度格式化
func WithMetricPrefix(params string) string {
	a, err := strconv.Atoi(params)
	if err != nil {
		return params
	}
	n := float64(a)
	if n < 1000 {
		return fmt.Sprintf("%6.2f", n)
	}

	i := 0
	units := []string{"K", "M", "G", "T", "P"}

	for n > 1000 {
		n = n / 1000
		i++
	}
	return fmt.Sprintf("%6.3f", n) + " " + units[i-1]
}

//每天平均收益 例如eth GetAVGEstreward(12234242,3)
func GetAVGEstreward(difficulty string, coinCount int) string {
	d, _ := strconv.Atoi(difficulty)
	return fmt.Sprintf("%.3f", 1000000000*3600*24*float32(coinCount)*0.99/float32(d)) + "GH/s"
}

func MustParseDuration(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		panic("util: Can't parse duration `" + s + "`: " + err.Error())
	}
	return value
}
