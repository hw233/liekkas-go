package coordinate

import (
	"testing"
)

func TestCubeRange(t *testing.T) {
	cubeRange := CubeRange(*NewPosition(-19, -2), 12)
	t.Log(cubeRange)
}

func TestCubeDistance(t *testing.T) {
	distance := CubeDistance(*NewPosition(0, -2), *NewPosition(-1, -3))
	t.Log(distance)

}

func TestCubeRing(t *testing.T) {
	//ring := CubeRing(NewPosition(0, 0), 1)
	//t.Log(ring)
	var x int32 = -19
	var y int32 = -2
	var n int32 = 12
	cubeRange := CubeRange(*NewPosition(x, y), n)

	posMap := map[Position]bool{}
	for _, position := range cubeRange {
		posMap[position] = true
	}
	for i := 1; i <= int(n); i++ {
		ring := CubeRing(*NewPosition(x, y), int32(i))
		for _, position := range ring {
			if posMap[position] == false {
				t.Errorf("err %v", position)
			}
			delete(posMap, position)

		}
	}
	t.Logf("rest %v", posMap)

}

func TestCubeRing2(t *testing.T) {
	cubeRing := CubeRing(*NewPosition(1, 1), 0)
	t.Logf("cubeRing %v", cubeRing)

}

func TestCubeContain(t *testing.T) {
	contain := CubeContain(*NewPosition(-19, -2), 12, *NewPosition(0, 0))
	t.Log(contain)
	contain = CubeContain(*NewPosition(0, 0), 0, *NewPosition(0, 0))
	t.Log(contain)

}
