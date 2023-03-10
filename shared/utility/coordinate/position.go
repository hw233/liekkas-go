package coordinate

import "shared/protobuf/pb"

type Position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

func NewPosition(x, y int32) *Position {
	return &Position{
		X: x,
		Y: y,
	}
}

func NewPositionFromVo(vo *pb.VOPosition) *Position {
	return NewPosition(vo.PosX, vo.PosY)
}

func (p *Position) VOPosition() *pb.VOPosition {
	return &pb.VOPosition{
		PosX: p.X,
		PosY: p.Y,
	}
}

func (p *Position) Clone() *Position {
	return NewPosition(p.X, p.Y)
}
