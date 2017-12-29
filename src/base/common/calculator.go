package common

// 利用后缀表达式计算
/*
*	将中缀表达式转换为后缀表达式的算法思想：
*	·开始扫描；
*	·数字时，加入后缀表达式；
*	·运算符：
*		a. 若为 '('，入栈；
*		b. 若为 ')'，则依次把栈中的的运算符加入后缀表达式中，直到出现'('，从栈中删除'(' ；
*		c. 若为 除括号外的其他运算符， 当其优先级高于除'('以外的栈顶运算符时，直接入栈。否则从栈顶开始，依次弹出比当前处理的运算符优先级高和优先级相等的运算符，直到一个比它优先级低的或者遇到了一个左括号为止。
*	·当扫描的中缀表达式结束时，栈中的的所有运算符出栈；
 */

import (
	l4g "base/log4go"
	"math"
	"strconv"
)

type Calculator struct {
}

func NewCalculator() *Calculator {
	return &Calculator{}
}

type StackNode struct {
	Data interface{}
	next *StackNode
}

type LinkStack struct {
	top   *StackNode
	Count int
}

func (this *LinkStack) Init() {
	this.top = nil
	this.Count = 0
}

func (this *LinkStack) Push(data interface{}) {
	var node *StackNode = new(StackNode)
	node.Data = data
	node.next = this.top
	this.top = node
	this.Count++
}

func (this *LinkStack) Pop() interface{} {
	if this.top == nil {
		return nil
	}
	returnData := this.top.Data
	this.top = this.top.next
	this.Count--
	return returnData
}

//Look up the top element in the stack, but not pop.
func (this *LinkStack) LookTop() interface{} {
	if this.top == nil {
		return nil
	}
	return this.top.Data
}

//检查公式是否合法
func (this *Calculator) Check(data string) bool {
	return true
}

func (this *Calculator) Count(data string, param map[string]float64) float64 {
	//TODO 检查字符串输入
	var arr []string = this.GenerateRPN(data)
	return this.calculateRPN(arr, param)
}

func (this *Calculator) Execute(arr []string, param map[string]float64) float64 {
	return this.calculateRPN(arr, param)
}

func (this *Calculator) calculateRPN(datas []string, param map[string]float64) float64 {
	l4g.Debug("Calculator calculateRPN datas is :%+v %+v", datas, param)
	var stack LinkStack
	stack.Init()
	for i := 0; i < len(datas); i++ {
		//l4g.Debug("Calculator isNumberString datas is i %d data:%s", i, datas[i])
		if this.isNumberString(datas[i]) {
			//l4g.Debug("Calculator isNumberString datas is :%s", datas[i])
			if this.isParamString(datas[i]) {
				p := datas[i]
				if f, exists := param[p]; exists {
					stack.Push(f)
				} else {
					l4g.Error("Calculator no find  param :%s", p)
					panic("operatin process go wrong.")
				}
			} else {
				if f, err := strconv.ParseFloat(datas[i], 64); err != nil {
					panic("operatin process go wrong.")
				} else {
					stack.Push(f)
				}
			}
		} else {
			p1 := stack.Pop().(float64)
			p2 := stack.Pop().(float64)
			p3 := this.normalCalculate(p2, p1, datas[i])
			stack.Push(p3)
		}
	}
	res := stack.Pop().(float64)
	return res
}

func (this *Calculator) normalCalculate(a, b float64, operation string) float64 {
	switch operation {
	case "*":
		return a * b
	case "-":
		return a - b
	case "+":
		return a + b
	case "/":
		return a / b
	case "max":
		if a >= b {
			return a
		} else {
			return b
		}
	case "min":
		if a >= b {
			return b
		} else {
			return a
		}
	case "pow":
		return math.Pow(a, b)
	default:
		panic("invalid operator")
	}
}

func (this Calculator) GenerateRPN(exp string) []string {
	var stack LinkStack
	var stack_op LinkStack
	stack.Init()
	stack_op.Init()

	var spiltedStr []string = this.convertToStrings(exp)
	l4g.Debug("Calculator spiltedStr is :%+v", spiltedStr)
	var datas []string

	for i := 0; i < len(spiltedStr); i++ { // 遍历每一个字符
		tmp := spiltedStr[i]                //当前字符
		if this.IsBaseOperatorString(tmp) { //是否是操作符
			// 四种情况入栈
			// 1 左括号直接入栈
			// 2 栈内为空直接入栈
			// 3 栈顶为左括号，直接入栈
			// 4 当前元素不为右括号时，在比较栈顶元素与当前元素，如果当前元素大，直接入栈。
			if tmp == "," {
				tmp = stack_op.Pop().(string) //替换为正常的操作符
			}
			// l4g.Debug("tmp  is :%s", tmp)
			if tmp == "(" ||
				stack.LookTop() == nil ||
				stack.LookTop().(string) == "(" ||
				(this.compareOperator(tmp, stack.LookTop().(string)) == 1 && tmp != ")") {
				stack.Push(tmp)
				// l4g.Debug("%s push 1", tmp)
			} else { // ) priority
				if tmp == ")" { //当前元素为右括号时，提取操作符，直到碰见左括号
					for {
						if pop := stack.Pop().(string); pop == "(" {
							break
						} else {
							datas = append(datas, pop)
						}
					}
				} else { //当前元素为操作符时，不断地与栈顶元素比较直到遇到比自己小的（或者栈空了或者为左括号），然后入栈。
					for {
						pop := stack.LookTop()
						// l4g.Debug("look top is %s", pop)
						if pop != nil && pop != "(" && this.compareOperator(tmp, pop.(string)) != 1 {
							datas = append(datas, stack.Pop().(string))
						} else {
							stack.Push(tmp)
							l4g.Debug("%s push 2", tmp)
							break
						}
					}
				}
			}
		} else if this.IsOperatorString(tmp) {
			stack_op.Push(tmp)
		} else {
			datas = append(datas, tmp)
		}
	}

	//将栈内剩余的操作符全部弹出。
	for {
		if pop := stack.Pop(); pop != nil {
			datas = append(datas, pop.(string))
		} else {
			break
		}
	}
	return datas
}

// if return 1, o1 > o2.
// if return 0, o1 = 02
// if return -1, o1 < o2
func (this *Calculator) compareOperator(o1, o2 string) int {
	// + - * /
	var o1Priority int
	if o1 == "+" || o1 == "-" {
		o1Priority = 1
	} else {
		o1Priority = 2
	}
	var o2Priority int
	if o2 == "+" || o2 == "-" {
		o2Priority = 1
	} else {
		o2Priority = 2
	}
	if o1Priority > o2Priority {
		return 1
	} else if o1Priority == o2Priority {
		return 0
	} else {
		return -1
	}
}

func (this *Calculator) isNumberString(o1 string) bool {
	if o1 == "+" || o1 == "-" || o1 == "*" || o1 == "/" || o1 == "(" || o1 == ")" || o1 == "max" || o1 == "min" || o1 == "pow" {
		return false
	} else {
		return true
	}
}

func (this *Calculator) convertToStrings(s string) []string {
	var strs []string
	bys := []byte(s)
	var tmp string
	for i := 0; i < len(bys); i++ {
		if this.isNumber(bys[i]) {
			tmp = tmp + string(bys[i])
		} else if this.IsBaseOperator(bys[i]) {
			if tmp != "" {
				strs = append(strs, tmp)
				tmp = ""
			}
			strs = append(strs, string(bys[i]))
		} else if this.isParam(bys[i]) {
			tmp = tmp + string(bys[i])
		} else if this.IsOperator(bys[i]) {
			tmp = tmp + string(bys[i])
		} else {
			//空格之类的
		}
	}
	if len(tmp) > 0 {
		strs = append(strs, tmp)
	}
	return strs
}

func (this *Calculator) isNumber(o1 byte) bool {
	if (o1 >= '0' && o1 <= '9') || o1 == '.' {
		return true
	} else {
		return false
	}
}

//+ - * / ( )
func (this *Calculator) IsBaseOperator(o1 byte) bool {
	if o1 == '+' || o1 == '-' || o1 == '*' || o1 == '/' || o1 == '(' || o1 == ')' || o1 == ',' {
		return true
	} else {
		return false
	}
}

func (this *Calculator) IsOperator(o1 byte) bool {
	if o1 >= 'a' && o1 <= 'z' {
		return true
	} else {
		return false
	}
}

//大写的字母为参数
func (this *Calculator) isParam(o1 byte) bool {
	if o1 == '$' || (o1 >= 'A' && o1 <= 'Z') || o1 == '_' {
		return true
	} else {
		return false
	}
}

func (this *Calculator) isParamString(o1 string) bool {
	bys := []byte(o1)
	if bys[0] == '$' {
		return true
	}
	return false
}

func (this *Calculator) IsBaseOperatorString(o1 string) bool {
	if o1 == "+" || o1 == "-" || o1 == "*" || o1 == "/" || o1 == "(" || o1 == ")" || o1 == "," {
		return true
	} else {
		return false
	}
}

func (this *Calculator) IsOperatorString(o1 string) bool {
	if o1 == "max" || o1 == "min" || o1 == "pow" {
		return true
	} else {
		return false
	}
}
