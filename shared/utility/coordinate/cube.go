package coordinate

import (
	"fmt"
	"math"
	"shared/utility/errors"
	utilMath "shared/utility/math"
)

/**
 * cube坐标系（六边形 2D坐标系）
 * 游戏内坐标系是偏移坐标:
 * 逢单数行会往右偏移半格
 *
 * 计算逻辑来自: https://www.redblobgames.com/grids/hexagons/
 */

const AllDirection = 6

var cubeDirections []Cube

func init() {
	cubeDirections = append(cubeDirections, *NewCube(+1, -1, 0))
	cubeDirections = append(cubeDirections, *NewCube(+1, 0, -1))
	cubeDirections = append(cubeDirections, *NewCube(0, +1, -1))
	cubeDirections = append(cubeDirections, *NewCube(-1, +1, 0))
	cubeDirections = append(cubeDirections, *NewCube(-1, 0, +1))
	cubeDirections = append(cubeDirections, *NewCube(0, -1, +1))
}

type Cube struct {
	X, Y, Z int32
}

func NewCube(X, Y, Z int32) *Cube {
	return &Cube{
		X: X,
		Y: Y,
		Z: Z,
	}
}

func NewCubeFromPosition(position Position) *Cube {
	cubeX := position.X - int32(math.Floor(float64(position.Y)/2))
	cubeY := position.Y
	cubeZ := -cubeX - cubeY
	return NewCube(cubeX, cubeY, cubeZ)
}
func (c *Cube) ToPosition() Position {
	return *NewPosition(c.X+int32(math.Floor(float64(c.Y)/2)), c.Y)
}

func (c *Cube) Neighbor(direction int32) {
	c.Add(cubeDirections[direction])
}

func (c *Cube) Add(addVal Cube) {
	c.X += addVal.X
	c.Y += addVal.Y
	c.Z += addVal.Z

}

func Direction(direction, radius int32) Cube {
	cubeDirection := cubeDirections[direction]
	return *NewCube(cubeDirection.X*radius, cubeDirection.Y*radius, cubeDirection.Z*radius)
}

// CubeRange 得到坐标 x,y半径r内的点（包括自身）
func CubeRange(position Position, n int32) []Position {
	var positions []Position
	cube := NewCubeFromPosition(position)
	minX := cube.X - n
	maxX := cube.X + n
	minY := cube.Y - n
	maxY := cube.Y + n
	minZ := cube.Z - n
	maxZ := cube.Z + n

	for i := minX; i <= maxX; i++ {
		for j := utilMath.MaxInt32Between(minY, -i-maxZ); j <= utilMath.MinInt32Between(maxY, -i-minZ); j++ {
			positions = append(positions, NewCube(i, j, -i-j).ToPosition())
		}
	}
	return positions
}

// CubeRing 获得第N圈的点
func CubeRing(position Position, n int32) []Position {
	var positions []Position
	cube := NewCubeFromPosition(position)

	if n == 0 {
		positions = append(positions, cube.ToPosition())
		return positions
	}
	//从4方向开始
	cube.Add(Direction(4, n))
	for i := 0; i < AllDirection; i++ {
		for j := 0; j < int(n); j++ {
			positions = append(positions, cube.ToPosition())
			cube.Neighbor(int32(i))

		}
	}
	return positions
}

// CubeDistance 距离
func CubeDistance(from, to Position) int32 {
	cubeFrom := NewCubeFromPosition(from)
	cubeTo := NewCubeFromPosition(to)
	return cubeDistance(cubeFrom, cubeTo)
}

func cubeDistance(cubeFrom, cubeTo *Cube) int32 {
	return (utilMath.AbsInt32(cubeFrom.X-cubeTo.X) + utilMath.AbsInt32(cubeFrom.Y-cubeTo.Y) + utilMath.AbsInt32(cubeFrom.Z-cubeTo.Z)) / 2

}

// CubeContain 以center为中的r半径内是否包含position
func CubeContain(center Position, r int32, position Position) bool {
	cubeCenter := NewCubeFromPosition(center)
	cubePosition := NewCubeFromPosition(position)

	CenterMinX := cubeCenter.X - r
	CenterMaxX := cubeCenter.X + r
	CenterMinY := cubeCenter.Y - r
	CenterMaxY := cubeCenter.Y + r
	CenterMinZ := cubeCenter.Z - r
	CenterMaxZ := cubeCenter.Z + r

	return !(CenterMinX > cubePosition.X || cubePosition.X > CenterMaxX ||
		CenterMinY > cubePosition.Y || cubePosition.Y > CenterMaxY ||
		CenterMinZ > cubePosition.Z || cubePosition.Z > CenterMaxZ)
}

// 获得to相对from的方向
func CubeDirection(from, to Position) (int32, error) {
	cubeFrom := NewCubeFromPosition(from)
	cubeTo := NewCubeFromPosition(to)
	if cubeDistance(cubeFrom, cubeTo) != 1 {
		return -1, errors.New("distance not 1")
	}
	direction := *NewCube(cubeTo.X-cubeFrom.X, cubeTo.Y-cubeFrom.Y, cubeTo.Z-cubeFrom.Z)

	for i, tmp := range cubeDirections {
		if tmp == direction {
			return int32(i), nil
		}
	}
	return -1, errors.New(fmt.Sprintf("unknow dirction ,from:%+v to:%+v", from, to))
}
