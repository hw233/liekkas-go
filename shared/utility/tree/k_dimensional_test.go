package tree

import (
	"shared/utility/coordinate"
	"shared/utility/rand"
	"testing"
	"time"
)

func TestKdTree(t *testing.T) {

	//tree.Append(coordinate.NewPosition(7,2))
	//tree.Append(coordinate.NewPosition(5,4))
	//tree.Append(coordinate.NewPosition(9,6))
	//tree.Append(coordinate.NewPosition(2,3))
	//tree.Append(coordinate.NewPosition(8,1))
	//tree.Append(coordinate.NewPosition(4,7))
	//closest, b := tree.findClosest(coordinate.NewPosition(3, 5))

	var looptime int64 = 1000
	var rebuildCost int64
	var findClosestCost int64
	var findClosestlowCost int64
	for j := 0; j < int(looptime); j++ {

		tree := NewCube2DTree()

		//tree.Append(coordinate.NewPosition(3 ,3))
		//tree.Append(coordinate.NewPosition(5 ,1))
		//tree.Append(coordinate.NewPosition(7 ,3))
		//tree.Append(coordinate.NewPosition(9 ,4))
		//tree.Append(coordinate.NewPosition(10, 4))
		//tree.Append(coordinate.NewPosition(0, 4))
		//tree.Append(coordinate.NewPosition(1 ,10))
		//tree.Append(coordinate.NewPosition(3 ,7))
		//tree.Append(coordinate.NewPosition(7 ,9))
		//tree.Append(coordinate.NewPosition(9,10))

		for i := 0; i < 10000; i++ {
			tree.Append(*coordinate.NewPosition(rand.RangeInt32(-10000, 10000), rand.RangeInt32(-10000, 10000)))
		}
		before := time.Now().UnixNano()

		tree.Rebuild()
		after := time.Now().UnixNano()

		//t.Logf("Rebuild cost,%d", after-before)
		rebuildCost += after - before
		target := *coordinate.NewPosition(0, 0)
		//
		before = time.Now().UnixNano()
		closest, _ := tree.FindClosest(target)
		after = time.Now().UnixNano()
		//t.Logf("findClosest cost ,%d", after-before)
		findClosestCost += after - before

		before = time.Now().UnixNano()

		closestLow := tree.findClosestLow(target)
		after = time.Now().UnixNano()

		//t.Logf("findClosestLow cost ,%d", after-before)
		findClosestlowCost += after - before

		//t.Log(closest)
		distance := coordinate.CubeDistance(*closest, target)
		//t.Log(distance)
		distancelow := coordinate.CubeDistance(closestLow, target)
		//t.Log(closestLow)

		//t.Log(distancelow)

		if distance != distancelow {
			for position := range tree.posMap {
				t.Error(position)

			}
			break
		}

	}

	t.Logf("rebuld average cost         ,%f", float64(rebuildCost)/float64(looptime)/1000000)
	t.Logf("findClosest average cost    ,%f", float64(findClosestCost)/float64(looptime)/1000000)
	t.Logf("findClosestLow average cost ,%f", float64(findClosestlowCost)/float64(looptime)/1000000)

}
