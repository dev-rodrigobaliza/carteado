package card

import (
	"reflect"
	"testing"
)

func TestFace_String(t *testing.T) {
	tests := []struct {
		name string
		f    *Face
		want string
	}{
		{"ace", faceAce, stringAce},
		{"2", face2, string2},
		{"3", face3, string3},
		{"4", face4, string4},
		{"5", face5, string5},
		{"6", face6, string6},
		{"7", face7, string7},
		{"8", face8, string8},
		{"9", face9, string9},
		{"10", face10, string10},
		{"jack", faceJack, stringJack},
		{"queen", faceQueen, stringQueen},
		{"king2", faceKing, stringKing},
		{"joker", faceJoker, stringJoker},
		{"unknown", faceUnknown, stringUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Errorf("Face.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFace(t *testing.T) {
	type args struct {
		face  string
		value int
		joker bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Face
		wantErr bool
	}{
		{"ace", args{heightAce, 0, false}, faceAce, false},
		{"2", args{height2, 0, false}, face2, false},
		{"3", args{height3, 0, false}, face3, false},
		{"4", args{height4, 0, false}, face4, false},
		{"5", args{height5, 0, false}, face5, false},
		{"6", args{height6, 0, false}, face6, false},
		{"7", args{height7, 0, false}, face7, false},
		{"8", args{height8, 0, false}, face8, false},
		{"9", args{height9, 0, false}, face9, false},
		{"10", args{height10, 0, false}, face10, false},
		{"jack", args{heightJack, 0, false}, faceJack, false},
		{"queen", args{heightQueen, 0, false}, faceQueen, false},
		{"king", args{heightKing, 0, false}, faceKing, false},
		{"jocker", args{heightJoker, 0, false}, faceJoker, false},
		{"unknown", args{heightUnknown, 0, false}, nil, true},
		{"ace", args{heightAce, CUSTOM_VALUE, false}, faceAceCustomValue, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFace(tt.args.face, tt.args.value, tt.args.joker)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFace() = %v, want %v", got, tt.want)
			}
		})
	}
	//random no joker
	for i := 0; i < 100; i++ {
		face, err := NewFace("", 0, false)
		if (err != nil) {
			t.Errorf("NewFace() random with no joker error = %v", err)
			return
		}
		if face.Height == HeightJoker {
			t.Error("NewFace() random with no joker got joker")
		}
	}
	//random joker
	gotJoker := 0
	for i := 0; i < 100; i++ {
		face, err := NewFace("", 0, true)
		if (err != nil) {
			t.Errorf("NewFace() random with joker error = %v", err)
			return
		}
		if face.Height == HeightJoker {
			gotJoker++
		}
	}
	if gotJoker == 0 {
		t.Error("NewFace() random with joker got no joker")
	}
}
