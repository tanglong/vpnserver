package common

import "fmt"

type KEY interface {
	IsLess(KEY) bool
	IsEqual(KEY) bool
}

type TSBNodeBase struct {
	Key    KEY
	Value  interface{}
	dwSize uint32
}

type TSBNode struct {
	TSBNodeBase
	ptLeft  *TSBNode
	ptRight *TSBNode
}

func newTSBNode() *TSBNode {
	return &TSBNode{ptLeft: nil, ptRight: nil}
}

type TSBTree struct {
	ptRoot *TSBNode
}

func (this *TSBTree) Size() uint32 {
	if this.ptRoot == nil {
		return 0
	}

	return this.ptRoot.dwSize
}

func (this *TSBTree) Root() *TSBNode {
	return this.ptRoot
}

func NewSBTree() *TSBTree {
	return &TSBTree{nil}
}

func (this *TSBTree) ClearSBTree() {
	clearNode(this.ptRoot)
	this.ptRoot = nil
}

func clearNode(pNode *TSBNode) {
	if pNode == nil {
		return
	}

	clearNode(pNode.ptLeft)
	clearNode(pNode.ptRight)
	pNode = nil
}

func getSize(pNode *TSBNode) uint32 {
	if pNode != nil {
		return pNode.dwSize
	}

	return 0
}

func leftRotate(pNode **TSBNode) {
	pRight := (*pNode).ptRight
	(*pNode).ptRight = pRight.ptLeft
	pRight.ptLeft = *pNode
	pRight.dwSize = (*pNode).dwSize
	(*pNode).dwSize = getSize((*pNode).ptLeft) + getSize((*pNode).ptRight) + 1
	*pNode = pRight
}

func rightRotate(pNode **TSBNode) {
	pLeft := (*pNode).ptLeft
	(*pNode).ptLeft = pLeft.ptRight
	pLeft.ptRight = (*pNode)
	pLeft.dwSize = (*pNode).dwSize
	(*pNode).dwSize = getSize((*pNode).ptLeft) + getSize((*pNode).ptRight) + 1
	*pNode = pLeft
}

func maintain(pNode **TSBNode, bIsRightDeeper bool) {
	if *pNode == nil {
		return
	}

	if !bIsRightDeeper {
		if (*pNode).ptLeft == nil {
			return
		}

		var dwRSize uint32 = getSize((*pNode).ptRight)
		if getSize((*pNode).ptLeft.ptLeft) > dwRSize {
			rightRotate(pNode)
		} else if getSize((*pNode).ptLeft.ptRight) > dwRSize {
			leftRotate(&((*pNode).ptLeft))
			rightRotate(pNode)
		} else {
			return
		}

		maintain(&((*pNode).ptLeft), false)
	} else {
		if (*pNode).ptRight == nil {
			return
		}

		var dwLSize uint32 = getSize((*pNode).ptLeft)
		if getSize((*pNode).ptRight.ptRight) > dwLSize {
			leftRotate(pNode)
		} else if getSize((*pNode).ptRight.ptLeft) > dwLSize {
			rightRotate(&((*pNode).ptRight))
			leftRotate(pNode)
		} else {
			return
		}

		maintain(&((*pNode).ptRight), true)
	}

	maintain(pNode, false)
	maintain(pNode, true)
}

func (this *TSBTree) Insert(key KEY, value interface{}) {
	pNode := newTSBNode()
	pNode.dwSize = 0
	pNode.Key = key
	pNode.Value = value
	sbtInsert(&this.ptRoot, pNode)
}

func sbtInsert(pTree **TSBNode, pNode *TSBNode) {
	if *pTree == nil {
		pNode.dwSize = 1
		*pTree = pNode
		return
	}

	(*pTree).dwSize++
	bIsLeft := pNode.Key.IsLess((*pTree).Key)
	if bIsLeft {
		sbtInsert(&((*pTree).ptLeft), pNode)
	} else {
		sbtInsert(&((*pTree).ptRight), pNode)
	}

	maintain(pTree, !bIsLeft)
}

func (this *TSBTree) Delete(key KEY) bool {
	if pNode := sbtDelete(&this.ptRoot, key, false); pNode != nil {
		pNode = nil
		return true
	}
	return false
}

func (this *TSBTree) DeleteNode(pNode *TSBNode) bool {
	if pNode == nil {
		return false
	}

	return this.Delete(pNode.Key)
}

func sbtDelete(pNode **TSBNode, key KEY, bIsFind bool) *TSBNode {
	// 查找到key所在的节点，然后用该节点左子树中最大值节点来替换掉需要删除的节点
	if *pNode == nil {
		return nil
	}

	var pRecord *TSBNode = nil
	// 只支持每次删除掉一个节点
	if bIsFind && (*pNode).ptRight == nil {
		// 查找最右的子树中最右的节点来替代删除的节点
		pRecord = *pNode
		*pNode = (*pNode).ptLeft
		return pRecord
	}

	if !bIsFind && key.IsEqual((*pNode).Key) {
		// 如果查询到是相等的节点
		if (*pNode).dwSize == 1 {
			// 叶子节点，需要删除掉此节点
			pRecord = *pNode
			*pNode = nil
			return pRecord
		}

		if (*pNode).dwSize == 2 {
			// 单枝节点，需要子节点继承被删除的节点
			if (*pNode).ptLeft != nil {
				pRecord = (*pNode).ptLeft
				(*pNode).ptLeft = nil
			} else {
				pRecord = (*pNode).ptRight
				(*pNode).ptRight = nil
			}
		} else {
			// 当前节点的左子树的最大值作为pNode的代替节点
			pRecord = sbtDelete(&((*pNode).ptLeft), key, true)
		}

		if pRecord != nil {
			(*pNode).Key = pRecord.Key
			(*pNode).Value = pRecord.Value
		}

		(*pNode).dwSize--
		maintain(pNode, true)
	} else if !bIsFind && key.IsLess((*pNode).Key) {
		pRecord = sbtDelete(&((*pNode).ptLeft), key, false)
		if pRecord != nil {
			(*pNode).dwSize--
		}
	} else {
		pRecord = sbtDelete(&((*pNode).ptRight), key, bIsFind)
		if pRecord != nil {
			(*pNode).dwSize--
		}
	}

	if pRecord != nil {
		maintain(pNode, !bIsFind && key.IsLess((*pNode).Key))
	}

	return pRecord
}

// 树中键值小于key的结点个数
func (this *TSBTree) LessCount(key KEY) uint32 {
	return sbtLessCount(this.ptRoot, key)
}

func sbtLessCount(pNode *TSBNode, key KEY) uint32 {
	if pNode == nil {
		return 0
	}

	if key.IsLess(pNode.Key) {
		return sbtLessCount(pNode.ptLeft, key)
	} else {
		if key.IsEqual(pNode.Key) {
			return sbtLessCount(pNode.ptLeft, key)
		} else {
			return getSize(pNode.ptLeft) + 1 + sbtLessCount(pNode.ptRight, key)
		}
	}
}

// 树中键值大于或等于key的结点个数
func (this *TSBTree) NonLessCount(key KEY) uint32 {
	if this.ptRoot == nil {
		return 0
	}

	return this.ptRoot.dwSize - sbtLessCount(this.ptRoot, key)
}

// 查找第一个值为key的结点
func (this *TSBTree) FindFirst(key KEY) *TSBNode {
	return sbtFind(this.ptRoot, key)
}

func sbtFind(pNode *TSBNode, key KEY) *TSBNode {
	if pNode == nil {
		return nil
	}

	if key.IsEqual(pNode.Key) {
		return pNode
	} else if key.IsLess(pNode.Key) {
		return sbtFind(pNode.ptLeft, key)
	} else {
		return sbtFind(pNode.ptRight, key)
	}
}

// 查找键值为key的结点个数
func (this *TSBTree) FindCount(key KEY) uint32 {
	return sbtFindCount(this.ptRoot, key)
}

func sbtFindCount(pNode *TSBNode, key KEY) uint32 {
	if pNode == nil {
		return 0
	}

	if key.IsEqual(pNode.Key) {
		return 1 + sbtFindCount(pNode.ptLeft, key) + sbtFindCount(pNode.ptRight, key)
	} else if key.IsLess(pNode.Key) {
		return sbtFindCount(pNode.ptLeft, key)
	} else {
		return sbtFindCount(pNode.ptRight, key)
	}
}

func (this *TSBTree) Head() *TSBNode {
	if this.ptRoot == nil {
		return nil
	}

	pNode := this.ptRoot
	for pNode.ptLeft != nil {
		pNode = pNode.ptLeft
	}

	return pNode
}

func (this *TSBTree) Tail() *TSBNode {
	if this.ptRoot == nil {
		return nil
	}

	pNode := this.ptRoot
	for pNode.ptRight != nil {
		pNode = pNode.ptRight
	}

	return pNode
}

// 比key小的最大节点
func (this *TSBTree) PreNode(key KEY) *TSBNode {
	if this.ptRoot == nil {
		return nil
	}

	pNode := this.ptRoot
	var pSave *TSBNode = nil
	for pNode != nil {
		bIsLess := pNode.Key.IsLess(key)
		if bIsLess {
			pSave = pNode
			pNode = pNode.ptRight
		} else {
			pNode = pNode.ptLeft
		}
	}

	return pSave
}

// 比key大的最小节点
func (this *TSBTree) SucNode(key KEY) *TSBNode {
	if this.ptRoot == nil {
		return nil
	}

	pNode := this.ptRoot
	var pSave *TSBNode = nil
	for pNode != nil {
		bIsLess := key.IsLess(pNode.Key)
		if bIsLess {
			pSave = pNode
			pNode = pNode.ptLeft
		} else {
			pNode = pNode.ptRight
		}
	}

	return pSave
}

// 选择排名为dwIndex的节点(从0开始到size-1之间，超过的返回空指针)
func (this *TSBTree) Rank(dwIndex uint32) *TSBNode {
	if this.ptRoot == nil {
		return nil
	}

	dwNowIndex := dwIndex
	var dwLeftSize uint32
	pNode := this.ptRoot
	for pNode != nil {
		dwLeftSize = getSize(pNode.ptLeft)
		if dwNowIndex == dwLeftSize {
			return pNode
		} else if dwNowIndex < dwLeftSize {
			pNode = pNode.ptLeft
		} else {
			dwNowIndex -= (dwLeftSize + 1)
			pNode = pNode.ptRight
		}
	}

	return nil
}

func (this *TSBTree) PopRoot() {
	if this.ptRoot == nil {
		return
	}

	this.Delete(this.ptRoot.Key)
}

func (this *TSBTree) PopHead() {
	if this.ptRoot == nil {
		return
	}

	pNode := sbtPopHead(&this.ptRoot)
	pNode.dwSize = 0
	pNode = nil
}

func sbtPopHead(pNode **TSBNode) *TSBNode {
	var pRecord *TSBNode = nil
	if (*pNode).ptLeft == nil {
		pRecord = *pNode
		*pNode = (*pNode).ptRight
		return pRecord
	}

	(*pNode).dwSize--
	pRecord = sbtPopHead(&((*pNode).ptLeft))
	maintain(pNode, true)
	return pRecord
}

func (this *TSBTree) PopTail() {
	if this.ptRoot == nil {
		return
	}

	pNode := sbtPopTail(&(this.ptRoot))
	pNode.dwSize = 0
	pNode = nil
}

func sbtPopTail(pNode **TSBNode) *TSBNode {
	var pRecord *TSBNode = nil
	if (*pNode).ptRight == nil {
		pRecord = *pNode
		*pNode = (*pNode).ptLeft
		return pRecord
	}

	(*pNode).dwSize--
	pRecord = sbtPopTail(&((*pNode).ptRight))
	maintain(pNode, false)
	return pRecord
}

func (this *TSBTree) Print() {
	print(this.ptRoot, 0, true)
}

func getAllKey(ptRoot *TSBNode, slcKeys *[]KEY) {
	if ptRoot == nil {
		return
	}

	getAllKey(ptRoot.ptLeft, slcKeys)
	*slcKeys = append(*slcKeys, ptRoot.Key)
	getAllKey(ptRoot.ptRight, slcKeys)
}

func (this *TSBTree) GetAllKey(slcKeys *[]KEY) {
	getAllKey(this.ptRoot, slcKeys)
}

func getAllValue(ptRoot *TSBNode, slcVals *[]interface{}) {
	if ptRoot == nil {
		return
	}

	getAllValue(ptRoot.ptLeft, slcVals)
	*slcVals = append(*slcVals, ptRoot.Value)
	getAllValue(ptRoot.ptRight, slcVals)
}

func (this *TSBTree) GetAllValue(slcVals *[]interface{}) {
	getAllValue(this.ptRoot, slcVals)
}

func print(pNode *TSBNode, dwHeight uint32, bIsLeft bool) {
	if pNode == nil {
		return
	}

	print(pNode.ptLeft, dwHeight+1, true)
	var dwI uint32 = 0
	for ; dwI < dwHeight; dwI++ {
		fmt.Print("        ")
	}

	if bIsLeft {
		fmt.Print("L")
	} else {
		fmt.Print("R")
	}

	fmt.Println(pNode.Key, ",", pNode.dwSize)
	fmt.Println()
	print(pNode.ptRight, dwHeight+1, false)
}

/*以下为排行榜涉及接口和结构及方法*/
type Element interface {
	KEY
	UniqueId() int64
}

type TSBTreeRank struct {
	ptSBTree    *TSBTree
	mapUidEle   map[int64]Element
	dwRankTotal uint32
}

func NewSBTreeRank(dwTotal uint32) *TSBTreeRank {
	return &TSBTreeRank{ptSBTree: NewSBTree(), mapUidEle: make(map[int64]Element), dwRankTotal: dwTotal}
}

//返回插入的Element在排行榜中的名次，不能进入排行榜返回-1
func (this *TSBTreeRank) Insert(ele Element) (int32, []Element) {
	slcRet := make([]Element, 0, 2)

	if tmp, ok := this.mapUidEle[ele.UniqueId()]; ok {
		this.ptSBTree.Delete(tmp)
		this.ptSBTree.Insert(ele, nil)
		this.mapUidEle[ele.UniqueId()] = ele
		slcRet = append(slcRet, ele)
		return int32(this.ptSBTree.LessCount(ele)) + 1, slcRet
	} else {
		this.ptSBTree.Insert(ele, nil)
		if this.ptSBTree.Size() > this.dwRankTotal {
			ptTail := this.ptSBTree.Tail()
			if ptTail.Key.IsEqual(ele) {
				this.ptSBTree.PopTail()
				return -1, slcRet
			} else {
				slcRet = append(slcRet, ele)
				slcRet = append(slcRet, ptTail.Key.(Element))
				delete(this.mapUidEle, ptTail.Key.(Element).UniqueId())
				this.ptSBTree.PopTail()
				this.mapUidEle[ele.UniqueId()] = ele
				return int32(this.ptSBTree.LessCount(ele)) + 1, slcRet
			}
		} else {
			slcRet = append(slcRet, ele)
			this.mapUidEle[ele.UniqueId()] = ele
			return int32(this.ptSBTree.LessCount(ele)) + 1, slcRet
		}
	}
}

func (this *TSBTreeRank) DelByUid(uid int64) Element {
	if tmp, ok := this.mapUidEle[uid]; ok {
		this.ptSBTree.Delete(tmp)
		delete(this.mapUidEle, uid)
		return tmp
	}

	return nil
}

//获取指定uid的排名，不在排行榜中返回-1，在则返回[1-dwRankTotal]
func (this *TSBTreeRank) GetRankByUid(uid int64) int32 {
	if val, ok := this.mapUidEle[uid]; ok {
		return int32(this.ptSBTree.LessCount(val)) + 1
	}

	return -1
}

//通过指定uid获得Element
func (this *TSBTreeRank) GetElementByUid(uid int64) Element {
	if val, ok := this.mapUidEle[uid]; ok {
		return val
	}

	return nil
}

//获取指定排名的Element，从1开始
func (this *TSBTreeRank) GetSpecRankElement(dwRank uint32) Element {
	if 0 == dwRank || this.ptSBTree.Size() < dwRank {
		return nil
	}
	ptRet := this.ptSBTree.Rank(dwRank - 1)
	if ptRet == nil {
		return nil
	}

	return ptRet.Key.(Element)
}

//获取排行榜中所有排好序的Element
func (this *TSBTreeRank) GetAllElement() []Element {
	slcKeys := make([]KEY, 0, this.ptSBTree.Size())
	this.ptSBTree.GetAllKey(&slcKeys)

	slcRet := make([]Element, 0, this.ptSBTree.Size())
	for _, val := range slcKeys {
		slcRet = append(slcRet, val.(Element))
	}

	return slcRet
}

//清空排行榜
func (this *TSBTreeRank) ClearRank() {
	this.mapUidEle = make(map[int64]Element)
	this.ptSBTree.ClearSBTree()
}
