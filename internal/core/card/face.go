package card

import (
	"github.com/dev-rodrigobaliza/carteado/errors"
)

type Face struct {
	Height Height
	Value  int
}

func NewFace(face string) (*Face, error) {
	if face == "" {
		return nil, errors.ErrInvalidCardFace
	}

	height := NewHeight(face, false)
	if height == HeightUnknown {
		return nil, errors.ErrInvalidCardFace
	}

	f := &Face{
		Height: height,
		Value:  int(height),
	}

	return f, nil
}

func NewFaceCustom(face string, value int) (*Face, error) {
	f, err := NewFace(face)
	if err != nil {
		return nil, err
	}

	f.Value = value

	return f, nil
}

func NewFaceRandom(joker bool) *Face {
	height := RandomHeight(joker)

	face := &Face{
		Height: height,
		Value:  int(height),
	}

	return face
}

func (f *Face) String() string {
	return f.Height.String()
}
