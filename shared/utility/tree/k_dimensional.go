package tree

import (
	"shared/utility/coordinate"
	"shared/utility/math"
	"sort"
)

const (
	XCUT int32 = 1
	YCUT int32 = 0
)

type Cube2DTree struct {
	root   *Node
	posMap map[coordinate.Position]bool
}

func NewCube2DTree() *Cube2DTree {
	return &Cube2DTree{
		posMap: map[coordinate.Position]bool{},
	}
}

type Node struct {
	value coordinate.Position
	left  *Node
	right *Node
	depth int32
}

func (c *Cube2DTree) Remove(posList ...coordinate.Position) {

	for _, position := range posList {
		delete(c.posMap, position)

	}
	//c.Rebuild()

}

func (c *Cube2DTree) Append(posList ...coordinate.Position) {

	for _, position := range posList {
		c.posMap[position] = true

	}
	//c.Rebuild()
}

func (c *Cube2DTree) Rebuild() {
	var list []coordinate.Position
	for position := range c.posMap {
		list = append(list, position)
	}

	c.root = rebuildRecursive(list, 1)
}
func (c *Cube2DTree) findClosestLow(target coordinate.Position) coordinate.Position {
	var list []coordinate.Position
	for position := range c.posMap {
		list = append(list, position)
	}

	distance := coordinate.CubeDistance(target, list[0])
	result := list[0]
	for _, position := range list {
		newDistance := coordinate.CubeDistance(target, position)
		if newDistance < distance {
			result = position
			distance = newDistance
		}

	}

	return result
}

func (c *Cube2DTree) FindClosest(target coordinate.Position) (*coordinate.Position, bool) {
	if c.root == nil {
		return nil, false
	}
	nodeList := NewNodeList()
	position, distance := c.root.search(target, nodeList, nil, 0)
	position, distance = nodeList.recall(target, position, distance)
	return position, true

}

func (n *Node) search(target coordinate.Position, nodeList *NodeList, position *coordinate.Position, distance int32) (*coordinate.Position, int32) {
	if n == nil {
		return position, distance
	}
	temp := n
	for temp.left != nil || temp.right != nil {
		*nodeList = append(*nodeList, temp)
		less := cutLess(temp.depth, target, temp.value)
		if less.X == target.X && less.Y == target.X {
			if temp.left == nil {
				break
			}
			temp = temp.left
		} else {
			if temp.right == nil {
				break
			}
			temp = temp.right
		}

	}
	if position == nil {
		return &temp.value, coordinate.CubeDistance(temp.value, target)
	}
	newDistance := coordinate.CubeDistance(temp.value, target)
	if newDistance < distance {
		return &temp.value, newDistance

	}
	return position, distance
}

type NodeList []*Node

func NewNodeList() *NodeList {
	return (*NodeList)(&[]*Node{})
}

func (n *NodeList) recall(target coordinate.Position, position *coordinate.Position, distance int32) (*coordinate.Position, int32) {
	for len(*n) != 0 {
		temp := (*n)[len(*n)-1]
		*n = append((*n)[:len(*n)-1])

		if temp.left == nil && temp.right == nil {
			newDistance := coordinate.CubeDistance(temp.value, target)
			if newDistance < distance {
				position, distance = &temp.value, newDistance

			}
		} else {
			if calCutDistance(target, temp.value, temp.depth) > distance {
				continue
			}
			newDistance := coordinate.CubeDistance(temp.value, target)
			if newDistance < distance {
				position, distance = &temp.value, newDistance
			}
			less := cutLess(temp.depth, target, temp.value)
			if less == target {
				if temp.right != nil {
					position, distance = temp.right.search(target, n, position, distance)
				}
			} else {
				if temp.left != nil {
					position, distance = temp.left.search(target, n, position, distance)
				}
			}

		}

	}
	return position, distance
}

func cutLess(depth int32, p1, p2 coordinate.Position) coordinate.Position {
	if calCut(depth) == XCUT {
		if p1.X < p2.X {
			return p1
		}
	} else {
		if p1.Y < p2.Y {
			return p1
		}
	}
	return p2
}

func calCutDistance(p1, p2 coordinate.Position, depth int32) int32 {
	if calCut(depth) == XCUT {
		return math.MaxInt32Between(p1.X, p2.X)
	} else {
		return math.MaxInt32Between(p1.Y, p2.Y)

	}
}
func rebuildRecursive(list []coordinate.Position, depth int32) *Node {
	node := &Node{}
	if len(list) == 0 {
		return nil
	}
	mid := len(list) / 2

	xless := func(i, j int) bool {
		return list[i].X < list[j].X
	}
	yless := func(i, j int) bool {
		return list[i].Y < list[j].Y
	}
	if calCut(depth) == XCUT {
		sort.Slice(list, xless)
	} else {
		sort.Slice(list, yless)
	}
	node.value = list[mid]
	node.depth = depth
	node.left = rebuildRecursive(list[0:mid], depth+1)
	node.right = rebuildRecursive(list[mid+1:], depth+1)
	return node
}

func calCut(depth int32) int32 {
	return depth / 2
}
