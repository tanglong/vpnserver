package zebra

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	l4g "base/log4go"

	"github.com/golang/protobuf/proto"
)

type Sessioner interface {
	Init(*Broker)              //初始化操作，比如心跳的设置...
	Process(*PackHead, []byte) //处理消息
	Close()                    //消除所有对Sessioner的引用,心跳...
}

const (
	StateInit = iota
	StateDisconnected
	StateConnected
)

type Broker struct {
	cn      *conn
	session Sessioner
	conf    *Config

	ReadMsgQueue  chan []byte
	writeMsgQueue chan *Message

	state     int32
	CloseChan chan struct{}

	wg sync.WaitGroup
}

func newBroker(se Sessioner, cf *Config) *Broker {
	return &Broker{
		session: se,
		conf:    cf,
	}
}

func (this *Broker) LocalAddr() string { return this.cn.localAddr }

func (this *Broker) RemoteAddr() string { return this.cn.remoteAddr }

func (this *Broker) State() int32 {
	return atomic.LoadInt32(&this.state)
}

func (this *Broker) Connect(timeout time.Duration) bool {
	rw, err := net.DialTimeout("tcp", this.conf.Address, timeout)
	if err != nil {
		l4g.Error("[Broker] Connect Error: %v", err)
		return false
	}
	if !this.serve(rw) {
		rw.Close()
		return false
	}
	return true
}

func (this *Broker) serve(rwc net.Conn) bool {
	if !atomic.CompareAndSwapInt32(&this.state, StateInit, StateConnected) {
		return false
	}

	this.cn = newconn(rwc, this)
	this.CloseChan = make(chan struct{})
	this.writeMsgQueue = make(chan *Message, this.conf.WriteMsgQueueSize)
	if this.conf.ReadMsgQueueSize > 0 {
		this.ReadMsgQueue = make(chan []byte, this.conf.ReadMsgQueueSize)
	}

	this.wg.Add(1)
	go this.cn.writeLoop()
	this.wg.Add(1)
	go this.cn.readLoop()

	this.session.Init(this)
	return true
}

func (this *Broker) AddWaitGroup() {
	this.wg.Add(1)
}

func (this *Broker) DecWaitGroup() {
	this.wg.Done()
}

func (this *Broker) transmitOrProcessMsg(buf []byte) {
	//l4g.Info("transmitOrProcessMsg %d %v", len(buf), buf)
	if this.conf.ReadMsgQueueSize > 0 {
		select {
		case this.ReadMsgQueue <- buf:
		case <-this.CloseChan:
		}
	} else {
		this.session.Process(GetInputMsgPackHead(buf), buf[16:])
	}
}

func GetInputMsgPackHead(buf []byte) *PackHead {
	return &PackHead{
		Length: uint32(len(buf) + 4),
		Cmd:    DecodeUint32(buf[0:]),
		Uid:    DecodeUint64(buf[4:]),
		Sid:    DecodeUint32(buf[12:]),
	}
}

func (this *Broker) Stop() {
	if !atomic.CompareAndSwapInt32(&this.state, StateConnected, StateDisconnected) {
		return
	}
	close(this.CloseChan)
	this.cn.rwc.Close()
	go func() {
		this.wg.Wait()
		this.cn.close()
		this.session.Close()
		l4g.Info("[Broker] closed, addr: (%s %s)", this.cn.localAddr, this.cn.remoteAddr)
		//this.cn = nil
		this.session = nil
		atomic.StoreInt32(&this.state, StateInit)
	}()
}

type Message struct {
	PH   *PackHead
	Info interface{}
}

func (this *Broker) Write(ph *PackHead, msg interface{}) bool {
	select {
	case this.writeMsgQueue <- &Message{ph, msg}:
		return true
	case <-this.CloseChan:
		l4g.Error("session close: %v %v", ph, msg)
		return false
	}
	/*
		if data, err := this.Marshal(ph, msg); err == nil {
			mq_len := len(this.writeMsgQueue)
			mq_cap := cap(this.writeMsgQueue)
			if mq_len > int(HIGH_WATER_MARK_SCALE*float64(mq_cap)) {
				l4g.Warn("[Broker] writeMsgQueue is HighWaterMark, len: %d cap: %d addr: (%s %s)",
					mq_len, mq_cap, this.cn.localAddr, this.cn.remoteAddr)
			}
			select {
			case this.writeMsgQueue <- data:
			case <-this.CloseChan:
			}
		}
	*/
}

/*
func (this *Broker) Marshal(ph *PackHead, msg interface{}) ([]byte, error) {
	var buf []byte
	switch v := msg.(type) {
	case []byte:
		buf = v
	case proto.Message:
		data, err := proto.Marshal(v)
		if err != nil {
			l4g.Error("[Broker] proto marshal cmd: %d sid: %d uid: %d error: %v",
				ph.Cmd, ph.Sid, ph.Uid, err)
			return nil, err
		}
		buf = data
	default:
		l4g.Error("[Broker] error msg type cmd: %d sid: %d uid: %d",
			ph.Cmd, ph.Sid, ph.Uid)
		return nil, ErrorMsgType
	}

	length := len(buf)
	length += PACK_HEAD_LEN
	if length > this.conf.MaxWriteMsgSize {
		l4g.Error("[Broker] write msg size overflow cmd: %d sid: %d uid: %d length: %d",
			ph.Cmd, ph.Sid, ph.Uid, length)
		return nil, WriteOverflow
	}

	data := make([]byte, length, length)
	ph.Length = uint32(length)
	l4g.Debug("[Broker] write head %v", ph)
	EncodePackHead(data, ph)
	copy(data[PACK_HEAD_LEN:], buf)
	return data, nil
}
*/

func (this *Broker) Marshal(ph *PackHead, msg interface{}) ([]byte, error) {
	var data []byte
	switch v := msg.(type) {
	case []byte:
		data = make([]byte, len(v)+PACK_HEAD_LEN)
		copy(data[PACK_HEAD_LEN:], v)
	case proto.Message:
		data = make([]byte, PACK_HEAD_LEN, 64)
		if mdata, err := proto.MarshalWithBytes(v, data); err == nil {
			data = mdata
		} else {
			l4g.Error("[Broker] proto marshal cmd: %d sid: %d uid: %d error: %v",
				ph.Cmd, ph.Sid, ph.Uid, err)
			return nil, err
		}
	default:
		l4g.Error("[Broker] error msg type cmd: %d sid: %d uid: %d",
			ph.Cmd, ph.Sid, ph.Uid)
		return nil, ErrorMsgType
	}

	length := len(data)
	if length > this.conf.MaxWriteMsgSize {
		l4g.Error("[Broker] write msg size overflow cmd: %d sid: %d uid: %d length: %d",
			ph.Cmd, ph.Sid, ph.Uid, length)
		return nil, WriteOverflow
	}

	ph.Length = uint32(length)
	//l4g.Debug("[Broker] write head %v", ph)
	EncodePackHead(data, ph)
	return data, nil
}
