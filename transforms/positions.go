package transforms

type Positions struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Rotations struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Transforms struct {
	Position Positions `json:"pozitions"`
	Rotation Rotations `json:"rotations"`
}

type Getter interface {
	Get()
}

type Setter interface {
	Set()
}

func (t *Transforms) Get() Transforms {
	return *t
}

func (t *Transforms) Set(newTransforms Transforms) {

	t.Position = newTransforms.Position
	t.Rotation = newTransforms.Rotation

}
