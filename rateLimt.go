package fyutils

import (
	"time"
	"sync"
	"fmt"
	"runtime/debug"
)

const (
	Close = 0
	Open = 1
)

// 限速
type RateLimit struct {
	sync.RWMutex
	Stats       int  			//状态  0：为关闭     1：打开
	Data   		map[string]map[string]*RateItem      // uid:{cmd:RateItem}
	Source   	map[string]int64					 //设置限速 原始数据  map[cmd]限制请求时间间隔
}

type RateItem struct {
	LastTimestamp	int64		//上次请求时间戳
	Rate 			int64	    //限制请求时间间隔(毫秒秒)
}

func (ri *RateLimit)CheckValid(area string,uid int32,cmd string)(bool,int64)  {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("CacheError, err:%v ", err))
			fmt.Println(string(debug.Stack()))
		}
	}()

	defer ri.Unlock()
	ri.Lock()

	if ri.Data==nil{
		ri.Data = make(map[string]map[string]*RateItem)
		return true,0
	}

	if ri.Stats == Close{
		return true,0
	}

	key := fmt.Sprintf("%s_%d",area,uid)
	if ri.Data[key]==nil{
		ri.Data[key] = make(map[string]*RateItem)
	}

	if ri.Data[key]==nil{
		ri.Data[key] = make(map[string]*RateItem)
		for k,v:=range ri.Source{
			ri.Data[key][k] = new(RateItem)
			ri.Data[key][k].Rate = v
		}
		return true,0
	}

	if ri.Data[key][cmd]==nil{
		return true,0
	}

	tm_now := time.Now().UnixNano()/1000000
	tm := tm_now - ri.Data[key][cmd].LastTimestamp
	if tm>=ri.Data[key][cmd].Rate{

		// 替换当前时间戳
		ri.Data[key][cmd].LastTimestamp = tm_now
		return true,tm
	}

	return false,tm
}

// limitCmd   map[cmd]限速时间
func (ri *RateLimit)SetRate(source map[string]int64)  {
	defer ri.Unlock()
	ri.Lock()

	ri.Source = source
	ri.Data = make(map[string]map[string]*RateItem)
}

func NewRate()*RateLimit  {
	rate:= new(RateLimit)
	rate.Data = make(map[string]map[string]*RateItem)
	rate.Source = make(map[string]int64)
	rate.Stats = Open
	return rate
}
