package astar

import (
	"container/heap"
	"errors"
	"math"
)

type Node struct {
	x, y, f, g, h int
	p             *Node
}

type PathQueue []*Node

func (this PathQueue) Len() int {
	return len(this)
}

func (this PathQueue) Less(i, j int) bool {
	return this[i].f < this[j].f
}

func (this PathQueue) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this *PathQueue) Push(x interface{}) {
	tmp := *this
	n := len(tmp)
	tmp = tmp[0 : n+1]
	node := x.(*Node)
	tmp[n] = node
	*this = tmp
}

func (this *PathQueue) Pop() interface{} {
	tmp := *this
	n := len(tmp)
	node := tmp[n-1]
	tmp[n-1] = nil
	*this = tmp[0 : n-1]
	return node
}

func (this *PathQueue) Clear() {
	tmp := *this
	for i := range tmp {
		tmp[i] = nil
	}
	*this = tmp[0:0]
}

var Dirs [8][2]int = [8][2]int{{1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}}

type Maper interface {
	In(int, int) bool
	IsWall(int, int, int, int) bool
	IsBlock(int, int) bool
}

type AStar struct {
	pq PathQueue
}

func New(a PathQueue) *AStar {
	return &AStar{pq: a}
}

func (this *AStar) heuristic(sx, sy, ex, ey int) int {
	return int(10 * (math.Abs(float64(ex-sx)) + math.Abs(float64(ey-sy))))
}

func (this *AStar) Clear() {
	this.pq.Clear()
}

func (this *AStar) FindMinPath(m Maper, sx, sy, ex, ey int) (int, int, error) {
	t_path := &Node{
		x: sx,
		y: sy,
		h: this.heuristic(sx, sy, ex, ey),
	}
	t_path.f = t_path.h
	heap.Push(&this.pq, t_path)
	var cx, cy, px, py int = 0, 0, 0, 0
	var visited [256][256]bool
	for this.pq.Len() > 0 {
		t_path = heap.Pop(&this.pq).(*Node)
		px, py = t_path.x, t_path.y
		if px == ex && py == ey {
			for tmp := t_path; tmp.p != nil; tmp = tmp.p {
				if tmp.p.x == sx && tmp.p.y == sy {
					return tmp.x, tmp.y, nil
				}
			}
		} else if !visited[px][py] {
			visited[px][py] = true
			for i := 0; i < 8; i++ {
				cx, cy = px+Dirs[i][0], py+Dirs[i][1]
				if !m.In(cx, cy) {
					continue
				}
				if i&1 == 1 {
					if m.IsWall(px, py, cx, cy) {
						continue
					}
				}
				if m.IsBlock(cx, cy) && !visited[cx][cy] {
					tmp := &Node{
						x: cx,
						y: cy,
						h: this.heuristic(cx, cy, ex, ey),
					}
					if i&1 == 1 {
						tmp.g = t_path.g + 14
					} else {
						tmp.g = t_path.g + 10
					}
					tmp.f = tmp.g + tmp.h
					tmp.p = t_path
					heap.Push(&this.pq, tmp)
				}
			}
		}
	}
	return 0, 0, errors.New("no find path")
}
