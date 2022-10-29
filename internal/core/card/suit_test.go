package card

import (
	"reflect"
	"testing"
)

func TestSuit_Graphic(t *testing.T) {
	tests := []struct {
		name string
		s    *Suit
		want string
	}{
		{"hearts", suitHearts, graphicHearts},
		{"diamonds", suitDiamonds, graphicDiamonds},
		{"clubs", suitClubs, graphicClubs},
		{"spades", suitSpades, graphicSpades},
		{"joker", suitJoker, graphicJoker},
		{"unknown", suitUnknown, graphicUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Graphic(); got != tt.want {
				t.Errorf("Suit.Graphic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSuit_String(t *testing.T) {
	tests := []struct {
		name string
		s    *Suit
		want string
	}{
		{"hearts", suitHearts, stringHearts},
		{"diamonds", suitDiamonds, stringDiamonds},
		{"clubs", suitClubs, stringClubs},
		{"spades", suitSpades, stringSpades},
		{"joker", suitJoker, stringJoker},
		{"unknown", suitUnknown, stringUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("Suit.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSuit(t *testing.T) {
	type args struct {
		suit  string
		value int
		joker bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Suit
		wantErr bool
	}{
		{"hearts", args{"h", 0, false}, suitHearts, false},
		{"hearts with value", args{"h", 5, false}, &Suit{SymbolHearts, 5}, false},
		{"diamonds", args{"d", 0, false}, suitDiamonds, false},
		{"clubs", args{"c", 0, false}, suitClubs, false},
		{"spades", args{"s", 0, false}, suitSpades, false},
		{"joker", args{"!", 0, false}, suitJoker, false},
		{"unknown", args{"?", 0, false}, nil, true},
		{"hearts", args{"h", CUSTOM_VALUE, false}, suitHeartsCustomValue, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSuit(tt.args.suit, tt.args.value, tt.args.joker)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSuit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSuit() = %v, want %v", got, tt.want)
			}
		})
	}
	//random no joker
	for i := 0; i < 100; i++ {
		suit, err := NewSuit("", 0, false)
		if (err != nil) {
			t.Errorf("NewSuit() random with no joker error = %v", err)
			return
		}
		if suit.Symbol == SymbolJoker {
			t.Error("NewSuit() random with no joker got joker")
		}
	}
	//random joker
	gotJoker := 0
	for i := 0; i < 100; i++ {
		suit, err := NewSuit("", 0, true)
		if (err != nil) {
			t.Errorf("NewSuit() random with joker error = %v", err)
			return
		}
		if suit.Symbol == SymbolJoker {
			gotJoker++
		}
	}
	if gotJoker == 0 {
		t.Error("NewSuit() random with joker got no joker")
	}
}
