package card

import (
	"reflect"
	"testing"
)

func TestSymbol_Graphic(t *testing.T) {
	tests := []struct {
		name string
		s    Symbol
		want string
	}{
		{"hearts", SymbolHearts, graphicHearts},
		{"diamonds", SymbolDiamonds, graphicDiamonds},
		{"clubs", SymbolClubs, graphicClubs},
		{"spades", SymbolSpades, graphicSpades},
		{"joker", SymbolJoker, graphicJoker},
		{"unknown", SymbolUnknown, graphicUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Graphic(); got != tt.want {
				t.Errorf("Symbol.Graphic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSymbol(t *testing.T) {
	type args struct {
		symbol string
		joker  bool
	}
	tests := []struct {
		name string
		args args
		want Symbol
	}{
		{"hearts", args{symbolHearts, false}, SymbolHearts},
		{"diamonds", args{symbolDiamonds, false}, SymbolDiamonds},
		{"clubs", args{symbolClubs, false}, SymbolClubs},
		{"spades", args{symbolSpades, false}, SymbolSpades},
		{"joker", args{symbolJoker, false}, SymbolJoker},
		{"unknown", args{symbolUnknown, false}, SymbolUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSymbol(tt.args.symbol, tt.args.joker); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
	//random no joker
	for i := 0; i < 100; i++ {
		symbol := NewSymbol("", false)
		if symbol == SymbolJoker {
			t.Error("NewSymbol() random with no joker got joker")
		}
	}
	//random joker
	gotJoker := 0
	for i := 0; i < 100; i++ {
		symbol := NewSymbol("", true)
		if symbol == SymbolJoker {
			gotJoker++
		}
	}
	if gotJoker == 0 {
		t.Error("NewSymbol() random with joker got no joker")
	}
}

func TestSymbol_String(t *testing.T) {
	tests := []struct {
		name string
		s    Symbol
		want string
	}{
		{"hearts", SymbolHearts, stringHearts},
		{"diamonds", SymbolDiamonds, stringDiamonds},
		{"clubs", SymbolClubs, stringClubs},
		{"spades", SymbolSpades, stringSpades},
		{"joker", SymbolJoker, stringJoker},
		{"unknown", SymbolUnknown, stringUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("Symbol.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
