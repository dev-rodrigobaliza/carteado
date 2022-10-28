package card

import (
	"github.com/dev-rodrigobaliza/carteado/errors"
)

type Face struct {
	Height Height
	Value  int
}

func NewFace(face string, value int, joker bool) (*Face, error) {
	var height Height

	if face == "" {
		height = RandomHeight(joker)
	} else {
		height = NewHeight(face, joker)
		if height == HeightUnknown {
			return nil, errors.ErrInvalidCardFace
		}
	}

	if value == 0 {
		value = int(height)
	}

	f := &Face{
		Height: height,
		Value: value,
	}

	return f, nil
}

func (f *Face) String() string {
	return f.Height.String()
}
