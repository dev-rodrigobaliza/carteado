package card

import (
	"reflect"
	"testing"
)

func TestNewHeight(t *testing.T) {
	type args struct {
		height string
		joker  bool
	}
	tests := []struct {
		name string
		args args
		want Height
	}{
		{"ace", args{heightAce, false}, HeightAce},
		{"2", args{height2, false}, Height2},
		{"3", args{height3, false}, Height3},
		{"4", args{height4, false}, Height4},
		{"5", args{height5, false}, Height5},
		{"6", args{height6, false}, Height6},
		{"7", args{height7, false}, Height7},
		{"8", args{height8, false}, Height8},
		{"9", args{height9, false}, Height9},
		{"10", args{height10, false}, Height10},
		{"jack", args{heightJack, false}, HeightJack},
		{"queen", args{heightQueen, false}, HeightQueen},
		{"king", args{heightKing, false}, HeightKing},
		{"joker", args{heightJoker, false}, HeightJoker},
		{"unknown", args{heightUnknown, false}, HeightUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHeight(tt.args.height, tt.args.joker); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHeight() = %v, want %v", got, tt.want)
			}
		})
	}
	//random no joker
	for i := 0; i < 100; i++ {
		height := NewHeight("", false)
		if height == HeightJoker {
			t.Error("NewHeight() 100 random with no joke got joke")
		}
	}
	//random joker
	gotJoker := 0
	for i := 0; i < 100; i++ {
		height := NewHeight("", true)
		if height == HeightJoker {
			gotJoker++
		}
	}
	if gotJoker == 0 {
		t.Error("NewHeight() 100 random with joke got no joke")
	}
}

func TestHeight_String(t *testing.T) {
	tests := []struct {
		name string
		h    Height
		want string
	}{
		{"ace", HeightAce, stringAce},
		{"2", Height2, string2},
		{"3", Height3, string3},
		{"4", Height4, string4},
		{"5", Height5, string5},
		{"6", Height6, string6},
		{"7", Height7, string7},
		{"8", Height8, string8},
		{"9", Height9, string9},
		{"10", Height10, string10},
		{"jack", HeightJack, stringJack},
		{"queen", HeightQueen, stringQueen},
		{"king2", HeightKing, stringKing},
		{"joker", HeightJoker, stringJoker},
		{"unknown", HeightUnknown, stringUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.String(); got != tt.want {
				t.Errorf("Height.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
