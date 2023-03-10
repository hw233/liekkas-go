package common

type Vec2 struct {
	X float64
	Y float64
}

func NewVec2(x, y float64) *Vec2 {
	return &Vec2{
		X: x,
		Y: y,
	}
}

func NewVec2ByVec2(vec2 *Vec2) *Vec2 {
	v := &Vec2{}

	v.Set(vec2)

	return v
}

func (v *Vec2) Set(vec2 *Vec2) {
	v.X = vec2.X
	v.Y = vec2.Y
}
