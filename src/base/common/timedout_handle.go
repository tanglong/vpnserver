package common

const INVALID_TIMEDOUT_HANDLE uint32 = 0

type ITimedOut interface {
	TimedOutHandle()
	ProcBeforeTimedOut(bIsLast bool, param interface{})
}

type TTimedoutKey struct {
	dwHandleId     uint32
	dwTimedOutTime uint32
	dwNeedCount    uint32
	iTimedOut      ITimedOut
}

func (this *TTimedoutKey) IsLess(b KEY) bool {
	c := b.(*TTimedoutKey)

	if this.dwTimedOutTime < c.dwTimedOutTime {
		return true
	} else if this.dwTimedOutTime == c.dwTimedOutTime {
		return this.dwHandleId < c.dwHandleId
	} else {
		return false
	}
}

func (this *TTimedoutKey) IsEqual(b KEY) bool {
	c := b.(*TTimedoutKey)

	return this.dwTimedOutTime == c.dwTimedOutTime && this.dwHandleId == c.dwHandleId
}

type TTimedOutHandler struct {
	ptSbTree       *TSBTree
	mapHandleIdKey map[uint32]*TTimedoutKey
	dwCurHandleId  uint32
}

func NewTimedOutHandler() *TTimedOutHandler {
	ptRet := &TTimedOutHandler{}
	ptRet.ptSbTree = NewSBTree()
	ptRet.mapHandleIdKey = make(map[uint32]*TTimedoutKey)
	ptRet.dwCurHandleId = 0
	return ptRet
}

func (this *TTimedOutHandler) CreateTimedOutHandle(dwTimedOutInterval, dwExpectCount uint32, iTimedOut ITimedOut) uint32 {
	ptKey := &TTimedoutKey{}

	this.dwCurHandleId++
	if this.dwCurHandleId == INVALID_TIMEDOUT_HANDLE {
		this.dwCurHandleId++
	}

	if dwExpectCount == 0 {
		dwExpectCount = 1
	}

	ptKey.dwHandleId = this.dwCurHandleId
	ptKey.dwNeedCount = dwExpectCount
	ptKey.dwTimedOutTime = CurDateTime().GetUnix() + dwTimedOutInterval
	ptKey.iTimedOut = iTimedOut

	this.ptSbTree.Insert(ptKey, nil)
	this.mapHandleIdKey[ptKey.dwHandleId] = ptKey
	return ptKey.dwHandleId
}

func (this *TTimedOutHandler) DeleteTimedOutHandle(dwHandleId uint32) {
	if val, ok := this.mapHandleIdKey[dwHandleId]; ok {
		this.ptSbTree.Delete(val)
		delete(this.mapHandleIdKey, dwHandleId)
	}
}

func (this *TTimedOutHandler) Update() {
	for this.ptSbTree.Size() != 0 {
		ptHead := this.ptSbTree.Head()
		ptKey := ptHead.Key
		ptRealKey := ptKey.(*TTimedoutKey)

		if ptRealKey.dwTimedOutTime <= CurDateTime().GetUnix() {
			this.ptSbTree.PopHead()
			delete(this.mapHandleIdKey, ptRealKey.dwHandleId)
			ptRealKey.iTimedOut.TimedOutHandle()
		} else {
			break
		}
	}
}

func (this *TTimedOutHandler) ProcBeforeTimedOut(dwHandleId uint32, param interface{}) {
	if val, ok := this.mapHandleIdKey[dwHandleId]; ok {
		val.dwNeedCount--
		if val.dwNeedCount == 0 {
			this.ptSbTree.Delete(val)
			delete(this.mapHandleIdKey, dwHandleId)
		}

		val.iTimedOut.ProcBeforeTimedOut(val.dwNeedCount == 0, param)
	}
}
