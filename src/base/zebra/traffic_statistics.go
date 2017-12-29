package zebra

import (
	"sync/atomic"
)

type TrafficStatistics struct {
	readMsgCount  uint32
	readMsgSize   uint64
	writeMsgCount uint32
	writeMsgSize  uint64
}

func (this *TrafficStatistics) setRead(rl uint64) {
	atomic.AddUint32(&this.readMsgCount, 1)
	atomic.AddUint64(&this.readMsgSize, rl)
}

func (this *TrafficStatistics) setWrite(wl uint64) {
	atomic.AddUint32(&this.writeMsgCount, 1)
	atomic.AddUint64(&this.writeMsgSize, wl)
}

func (this *TrafficStatistics) Get() (rms, wms uint32, rml, wml uint64) {
	rms = atomic.LoadUint32(&this.readMsgCount)
	wms = atomic.LoadUint32(&this.writeMsgCount)
	rml = atomic.LoadUint64(&this.readMsgSize)
	wml = atomic.LoadUint64(&this.writeMsgSize)
	return
}

//var TS = new(TrafficStatistics)
