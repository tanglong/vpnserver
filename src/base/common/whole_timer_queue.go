package common

import (
	"container/heap"
	"time"
)

const (
	WHOLE_TIME_START_TYPE_1 uint32 = 1 //服务器重启上个时间点开始算,重启时会执行一次
	WHOLE_TIME_START_TYPE_2 uint32 = 2 //服务器重启下个时间点开始算,重启后时间点到了才会执行
)

const (
	WHOLE_TIME_TYPE_EACH_5_MIN  uint32 = 1 //5分钟定时器
	WHOLE_TIME_TYPE_DAY_22_HOUR uint32 = 2 //22点定时器
	WHOLE_TIME_TYPE_DAY_5_HOUR  uint32 = 3 //5点定时器
)

type WholeTimeOuter interface {
	TimeOut(int64)
}

type WholeTimer struct {
	WholeTimeOuter
	id       uint32
	timet    int64 //开始时间
	interval int64 //迭代周期
	index    int
}

type WholeTimerQueue []*WholeTimer

func (this WholeTimerQueue) Len() int {
	return len(this)
}

func (this WholeTimerQueue) Less(i, j int) bool {
	return this[i].timet < this[j].timet
}

func (this WholeTimerQueue) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
	this[i].index = i
	this[j].index = j
}

func (this *WholeTimerQueue) Push(x interface{}) {
	tmp := *this
	n := len(tmp)
	tmp = tmp[0 : n+1]
	timer := x.(*WholeTimer)
	timer.index = n
	tmp[n] = timer
	*this = tmp
}

func (this *WholeTimerQueue) Pop() interface{} {
	tmp := *this
	n := len(tmp)
	timer := tmp[n-1]
	tmp[n-1] = nil
	timer.index = -1
	*this = tmp[0 : n-1]
	return timer
}

type WholeTimerManager struct {
	id uint32
	tq WholeTimerQueue
}

func NewWholeTimerManager(size int) *WholeTimerManager {
	if size == 0 {
		size = TIMER_QUEUE_SIZE
	}
	return &WholeTimerManager{tq: make([]*WholeTimer, 0, size)}
}

/*
i: 回调方法
start: 计时器类型
timet: 执行时间点
*/
func (this *WholeTimerManager) AddTimer(i WholeTimeOuter, timet, startt uint32) uint32 {
	if cap(this.tq) <= len(this.tq) {
		return 0
	}
	var iv, s int64
	switch timet {
	case WHOLE_TIME_TYPE_EACH_5_MIN:
		iv = int64(5 * 60)
		s = int64(DailyWholeMinute())
		if startt == WHOLE_TIME_START_TYPE_2 {
			s += iv
		}
	case WHOLE_TIME_TYPE_DAY_22_HOUR:
		iv = int64(24 * 3600)
		s = int64(DailyWhole22Hour())
		if startt == WHOLE_TIME_START_TYPE_2 {
			s += iv
		}
	case WHOLE_TIME_TYPE_DAY_5_HOUR:
		iv = int64(24 * 3600)
		s = int64(DailyWhole5Hour())
		if startt == WHOLE_TIME_START_TYPE_2 {
			s += iv
		}
	}
	timer := &WholeTimer{WholeTimeOuter: i, interval: iv, timet: s}
	this.id++
	timer.id = this.id
	heap.Push(&this.tq, timer)
	return timer.id
}

var queueEx *Queue = &Queue{}

func (this *WholeTimerManager) Run(now int64, limit int) {
	for len(this.tq) > 0 {
		tmp := this.tq[0]
		if tmp.timet <= now {
			timer := heap.Pop(&this.tq).(*WholeTimer)
			queueEx.Push(timer.WholeTimeOuter)
			if timer.interval > 0 {
				timer.timet += timer.interval
				heap.Push(&this.tq, timer)
			}
		} else {
			break
		}
		if limit > 0 && queueEx.Len() >= limit {
			break
		}
	}

	for queueEx.Len() > 0 {
		queueEx.Pop().(TimeOuter).TimeOut(now)
	}
}

//得到前一个整点5分的时间
func DailyWholeMinute() uint32 {
	now := time.Now()
	mins := []int{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60}
	hour, minute, _ := now.Hour(), now.Minute(), now.Second()
	min := mins[0]
	for index := 1; index < len(mins); index++ {
		if minute >= min && minute < mins[index] {
			break
		} else {
			min = mins[index]
		}
	}
	return DailyZero() + uint32(hour*3600) + uint32(min*60)
}

//得到前一个整点22时的时间
func DailyWhole22Hour() uint32 {
	now := time.Now()
	hour, minute, second := now.Hour(), now.Minute(), now.Second()
	if hour < 22 {
		return uint32(now.Unix() - int64((hour*3600)+(minute*60)+second) - 2*3600)
	}
	return uint32(now.Unix() - int64((hour*3600)+(minute*60)+second) + 22*3600)
}

//得到前一个整点5时的时间
func DailyWhole5Hour() uint32 {
	now := time.Now()
	hour, minute, second := now.Hour(), now.Minute(), now.Second()
	if hour < 5 {
		return uint32(now.Unix() - int64((hour*3600)+(minute*60)+second) - 19*3600)
	}
	return uint32(now.Unix() - int64((hour*3600)+(minute*60)+second) + 5*3600)
}
