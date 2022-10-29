package card

import (
	"reflect"
	"testing"
)

func TestCard_Graphic(t *testing.T) {
	type args struct {
		face bool
	}
	tests := []struct {
		name string
		c    *Card
		args args
		want string
	}{
		{"ace hearts without face", cardAceHearts, args{false}, graphicHearts},
		{"ace hearts with face", cardAceHearts, args{true}, (stringAce + graphicHearts)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Graphic(tt.args.face); got != tt.want {
				t.Errorf("Card.Graphic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCard_String(t *testing.T) {
	type args struct {
		face bool
		suit bool
	}
	tests := []struct {
		name string
		c    *Card
		args args
		want string
	}{
		{"ace hearts without face and suit", cardAceHearts, args{false, false}, ""},
		{"ace hearts with face without suit", cardAceHearts, args{true, false}, stringAce},
		{"ace hearts without face with suit", cardAceHearts, args{false, true}, stringHearts},
		{"ace hearts with face and suit", cardAceHearts, args{true, true}, (stringAce + stringHearts)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(tt.args.face, tt.args.suit); got != tt.want {
				t.Errorf("Card.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCard_Value(t *testing.T) {
	type args struct {
		face bool
		suit bool
	}
	tests := []struct {
		name string
		c    *Card
		args args
		want int
	}{
		{"ace hearts without face and suit", cardAceHearts, args{false, false}, 0},
		{"ace hearts with face without suit", cardAceHearts, args{true, false}, int(HeightAce)},
		{"ace hearts without face with suit", cardAceHearts, args{false, true}, int(SymbolHearts)},
		{"ace hearts with face and suit", cardAceHearts, args{true, true}, (int(HeightAce) + int(SymbolHearts))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Value(tt.args.face, tt.args.suit); got != tt.want {
				t.Errorf("Card.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		face      string
		suit      string
		faceValue int
		suitValue int
		joker     bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Card
		wantErr bool
	}{
		{"ace hearts default face and suit and default values", args{"1", "h", 0, 0, false}, cardAceHearts, false},
		{"ace hearts invalid face and default suit and default values", args{"a", "h", 0, 0, false}, nil, true},
		{"ace hearts default face and invalid suit and default values", args{"1", "a", 0, 0, false}, nil, true},
		{"ace hearts invalid face and suit and default values", args{"a", "a", 0, 0, false}, nil, true},
		{"ace hearts default face and suit and custom face and deafult suit values", args{"1", "h", CUSTOM_VALUE, 0, false}, cardAceHeartsCustomFaceValues, false},
		{"ace hearts default face and suit and default face and custom suit values", args{"1", "h", 0, CUSTOM_VALUE, false}, cardAceHeartsCustomSuitValues, false},
		{"ace hearts default face and suit and custom values", args{"1", "h", CUSTOM_VALUE, CUSTOM_VALUE, false}, cardAceHeartsCustomValues, false},
		{"joker default face and suit and default values", args{"!", "!", 0, 0, false}, cardJoker, false},
		{"joker forced face and default suit and default values", args{"1", "!", CUSTOM_VALUE, 0, false}, cardJokerForcedFace, false},
		{"joker default face and forced suit and default values", args{"!", "h", 0, CUSTOM_VALUE, false}, cardJokerForcedSuit, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.face, tt.args.suit, tt.args.faceValue, tt.args.suitValue, tt.args.joker)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
