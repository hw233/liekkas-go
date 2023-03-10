package common

import (
	"shared/utility/coordinate"
	"testing"
)

func TestArea_MergePiece(t *testing.T) {
	area := NewArea()

	area.AppendPiece(1, Piece{StartY: 1, EndY: 3})
	area.AppendPiece(1, Piece{StartY: 4, EndY: 8})
	area.AppendPiece(1, Piece{StartY: -1, EndY: 11})
	area.AppendPiece(1, Piece{StartY: 2, EndY: 123})

	for i, pieces := range area.Pieces {
		t.Logf("x :%+v", i)
		t.Logf("pieces :%+v", pieces)

	}

}

func TestArea_Cut(t *testing.T) {
	area := NewArea()
	area.AppendPiece(1, Piece{StartY: -1, EndY: 123})
	cut := area.CutRectangle(*coordinate.NewPosition(1, 1), 2, 2)
	for i, pieces := range cut.Pieces {
		t.Logf("x :%+v", i)
		t.Logf("pieces :%+v", pieces)

	}
}

func TestArea_Contains(t *testing.T) {
	area := NewArea()
	area.AppendPiece(1, Piece{StartY: -1, EndY: 123})
	contains := area.Contains(*coordinate.NewPosition(1, 124))
	t.Logf("contains :%+v", contains)

}
