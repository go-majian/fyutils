package manger

import (
	"testing"
	"runtime"
	"time"
	"fmt"
)

func TestRateLimit_CheckValid(t *testing.T) {

	runtime.GOMAXPROCS(runtime.NumCPU())
	rate := NewRate()

	//cmd:限制请求时间间隔（单位：毫秒）
	source:= map[string]int64{
		"listClubRoom":200,
		"我回来了":400,
		"有事离开，等我一下":200,
	}
	rate.SetRate(source)

	go func() {
		for i:=0;i<1000;i++{
			go func() {rate.SetRate(source)}()
		}
	}()
	for i:=0;i<100000;i++{
		go func() {
			fmt.Println(rate.CheckValid("area",10,"有事离开，等我一下"))
		}()
	}

	time.Sleep(100*time.Second)
}
