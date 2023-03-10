package common

import (
	"encoding/json"
	"fmt"
	"shared/utility/coordinate"
	"shared/utility/errors"
	utilMath "shared/utility/math"
	"sort"
	"strconv"
	"strings"
)

type Area struct {
	Pieces      map[int32]*Pieces
	CachePoints []coordinate.Position
}

func NewArea() *Area {
	return &Area{
		Pieces:      map[int32]*Pieces{},
		CachePoints: nil,
	}

}

func (a *Area) MarshalJSON() ([]byte, error) {
	str := ""
	for x, pieces := range a.Pieces {
		for _, piece := range pieces.Value {
			objStr := fmt.Sprintf("%d:%d:%d", x, piece.StartY, piece.EndY)
			if len(str) == 0 {
				str = objStr
			} else {
				str = fmt.Sprintf("%s,%s", str, objStr)

			}
		}

	}
	return json.Marshal(str)
}

func (a *Area) UnmarshalJSON(b []byte) error {
	org := string(b)

	split := strings.Split(org[1:len(org)-1], ",")

	return StringsToArea(split, a)
}

func StringsToArea(v []string, a *Area) error {
	for _, s := range v {
		if len(s) == 0 {
			continue
		}
		x, piece, err := StringToUnlockPiece(s)
		if err != nil {
			return err
		}
		a.AppendPiece(x, *piece)
	}
	a.Refresh()
	return nil
}

func stringToInt32(s string) (int32, error) {
	if s == "" {
		return 0, nil
	}
	i, err := strconv.ParseInt(s, 10, 32)
	return int32(i), err
}
func StringToUnlockPiece(v string) (int32, *Piece, error) {

	split := strings.Split(v, ":")
	if len(split) != 3 {
		return 0, nil, errors.WrapTrace(ErrCSVFormatInvalid)
	}
	X, err := stringToInt32(split[0])
	if err != nil {
		return 0, nil, errors.WrapTrace(err)
	}
	StartY, err := stringToInt32(split[1])
	if err != nil {
		return 0, nil, errors.WrapTrace(err)
	}
	EndY, err := stringToInt32(split[2])
	if err != nil {
		return 0, nil, errors.WrapTrace(err)
	}
	return X, &Piece{
		StartY: StartY,
		EndY:   EndY,
	}, nil
}

type Pieces struct {
	NeedMerge bool
	Value     []Piece
}

func NewPieces(value []Piece) *Pieces {
	return &Pieces{
		NeedMerge: false,
		Value:     value,
	}
}

type Piece struct {
	StartY int32
	EndY   int32
}

func NewPiece(StartY, EndY int32) Piece {
	return Piece{
		StartY: StartY,
		EndY:   EndY,
	}
}

// MergeArea 合并区域
func (a *Area) MergeArea(merge *Area) {
	for x, pieces := range merge.Pieces {
		for _, piece := range pieces.Value {
			a.AppendPiece(x, piece)
		}
	}
	a.Refresh()

}

func (a *Area) AppendPoint(positions ...coordinate.Position) {
	for _, position := range positions {

		newPiece := NewPiece(position.Y, position.Y)
		pieces, ok := a.Pieces[position.X]
		if !ok {
			pieces = NewPieces([]Piece{newPiece})
		} else {
			pieces.Value = append(pieces.Value, newPiece)
			pieces.NeedMerge = true
		}
		a.Pieces[position.X] = pieces
	}
	a.Refresh()

}

// AppendPiece 合并一条线段
func (a *Area) AppendPiece(x int32, merge Piece) {
	pieces, ok := a.Pieces[x]
	if !ok {
		pieces = NewPieces([]Piece{merge})
	} else {
		pieces.Value = append(pieces.Value, merge)
		pieces.NeedMerge = true
	}
	if a.Pieces == nil {
		a.Pieces = map[int32]*Pieces{}
	}
	a.Pieces[x] = pieces
	a.Refresh()
}

func (a *Area) Refresh() {
	newPiecesMap := map[int32]*Pieces{}

	for x, pieces := range a.Pieces {

		if len(pieces.Value) == 0 {
			continue
		}
		if !pieces.NeedMerge {
			newPiecesMap[x] = pieces
			continue
		}
		// 先按StartY从小到大排序
		less := func(i, j int) bool {
			return pieces.Value[i].StartY < pieces.Value[j].EndY
		}
		sort.Slice(pieces.Value, less)

		// 合并线段
		currPiece := pieces.Value[0]
		var newPieces []Piece
		for j := 1; j < len(pieces.Value); j++ {
			if currPiece.EndY >= pieces.Value[j].StartY-1 {
				currPiece = Piece{utilMath.MinInt32Between(currPiece.StartY, pieces.Value[j].StartY), utilMath.MaxInt32Between(currPiece.EndY, pieces.Value[j].EndY)}
			} else {
				newPieces = append(newPieces, currPiece)
				currPiece = pieces.Value[j]
			}
		}
		newPieces = append(newPieces, currPiece)
		newPiecesMap[x] = NewPieces(newPieces)
	}
	a.Pieces = newPiecesMap
	a.CachePoints = nil

}

func (p *Pieces) Minus(minus Piece) error {
	var newPieces []Piece
	fix := false

	for _, piece := range p.Value {
		if piece.StartY <= minus.StartY && piece.EndY >= minus.EndY {

			newPieces = append(newPieces, piece.Minus(minus)...)
			fix = true

		} else {
			newPieces = append(newPieces, piece)
		}
	}
	if !fix {
		return errors.WrapTrace(ErrGraveyardNoAreaCannotBuild)
	}
	p.Value = newPieces
	return nil

}

func (p Piece) Minus(minus Piece) []Piece {

	var pieces []Piece

	if p.StartY < minus.StartY {
		pieces = append(pieces, Piece{p.StartY, minus.StartY - 1})
	}
	if minus.EndY < p.EndY {
		pieces = append(pieces, Piece{minus.EndY + 1, p.EndY})

	}
	return pieces
}

func (a *Area) MinusBuildingArea(buildingArea BuildingArea, position *coordinate.Position) error {

	for i := position.X; i < position.X+buildingArea.HorizontalCount; i++ {

		pieces, ok := a.Pieces[i]
		if !ok {
			return errors.WrapTrace(ErrGraveyardNoAreaCannotBuild)
		}
		err := pieces.Minus(Piece{
			StartY: position.Y,
			EndY:   position.Y + buildingArea.VerticalCount - 1,
		})
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	return nil

}

func (a *Area) Contains(position coordinate.Position) bool {
	pieces, ok := a.Pieces[position.X]
	if !ok {
		return false
	}
	for _, piece := range pieces.Value {
		if position.Y >= piece.StartY && position.Y <= piece.EndY {
			return true
		}
	}
	return false
}

// CutRectangle 切下area一部分， position 表示左下角坐标， h,v分别表示水平地块数量和垂直地块数量
func (a *Area) CutRectangle(position coordinate.Position, h, v int32) *Area {
	result := NewArea()
	if h <= 0 || v <= 0 {
		return result
	}
	for i := position.X; i < position.X+h; i++ {
		pieces, ok := a.Pieces[i]
		if ok {
			for _, piece := range pieces.Value {
				if piece.EndY < position.Y {
					continue
				} else if piece.StartY > position.Y+h {
					continue
				} else {
					result.AppendPiece(i, Piece{utilMath.MaxInt32Between(position.Y, piece.StartY), utilMath.MinInt32Between(position.Y+v-1, piece.EndY)})
				}
			}
		}

	}
	return result
}

// Count 该区域包含的所有点数量
func (a *Area) Count() int32 {
	var ret int32
	for _, pieces := range a.Pieces {
		for _, piece := range pieces.Value {
			ret += piece.EndY - piece.StartY + 1
		}
	}
	return ret
}

func (a *Area) Points() []coordinate.Position {
	if a.CachePoints != nil {
		return a.CachePoints
	}

	var ret []coordinate.Position
	for x, pieces := range a.Pieces {
		for _, piece := range pieces.Value {
			for y := piece.StartY; y <= piece.EndY; y++ {
				ret = append(ret, *coordinate.NewPosition(x, y))
			}
		}
	}
	a.CachePoints = ret
	return ret
}

// CutArea 切下area一部分， 返回a和other重合的点
func (a *Area) CutArea(other *Area) *Area {
	ret := NewArea()

	var positions []coordinate.Position
	for _, position := range a.Points() {
		if other.Contains(position) {
			positions = append(positions, position)
		}
	}
	ret.AppendPoint(positions...)

	return ret
}

// MinusArea 返回a中和other不重合的点
func (a *Area) MinusArea(other *Area) *Area {
	ret := NewArea()
	var positions []coordinate.Position
	for _, position := range a.Points() {
		if !other.Contains(position) {
			positions = append(positions, position)
		}
	}
	ret.AppendPoint(positions...)
	return ret

}

func (a *Area) CoincidenceArea(other *Area) bool {
	for _, position := range a.Points() {
		if other.Contains(position) {
			return true
		}
	}
	return false

}

// BuildingArea 建筑占地块
type BuildingArea struct {
	HorizontalCount int32 // 水平地块数量
	VerticalCount   int32 // 垂直地块数量
}
