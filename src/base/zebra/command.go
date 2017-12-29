package zebra

import (
	l4g "base/log4go"
	"time"
)

var (
	DEBUG       bool            = true
	IGNORE_CODE map[uint32]bool = map[uint32]bool{60001: true, 50000: true, 70001: true, 80001: true}
)

type Services func(Sessioner, *PackHead, []byte) bool

type CommandM struct {
	cmdm map[uint32]Services
}

func NewCommandM() *CommandM {
	return &CommandM{
		cmdm: make(map[uint32]Services),
	}
}

func (this *CommandM) Register(id uint32, service Services) {
	this.cmdm[id] = service
}

func (this *CommandM) Dispatcher(session Sessioner, ph *PackHead, data []byte) bool {
	if cmd, exist := this.cmdm[ph.Cmd]; exist {
		t1 := time.Now().UnixNano()
		ret := cmd(session, ph, data)
		if DEBUG {
			t2 := time.Now().UnixNano()
			dt := t2 - t1
			if _, ok := IGNORE_CODE[ph.Cmd]; !ok {
				l4g.Trace("cmd: %d, cost time: %d", ph.Cmd, dt/int64(time.Millisecond))
			}
		}
		if !ret {
			l4g.Error("Handler err, cmd: %d", ph.Cmd)
		}
		return ret
	}
	l4g.Error("[Command] no find cmd: %d %d", ph.Sid, ph.Cmd)
	return false
}
