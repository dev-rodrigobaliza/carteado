package card

const (
	CUSTOM_VALUE = 10
)

var (
	stringJoker   = "joker"
	stringUnknown = "unknown"

	graphicHearts   = "♥"
	graphicDiamonds = "♦"
	graphicClubs    = "♣"
	graphicSpades   = "♠"
	graphicJoker    = "!"
	graphicUnknown  = "?"

	stringHearts   = "hearts"
	stringDiamonds = "diamonds"
	stringClubs    = "clubs"
	stringSpades   = "spades"

	symbolHearts   = "h"
	symbolDiamonds = "d"
	symbolClubs    = "c"
	symbolSpades   = "s"
	symbolJoker    = "!"
	symbolUnknown  = "?"

	suitHearts   = &Suit{SymbolHearts, int(SymbolHearts)}
	suitDiamonds = &Suit{SymbolDiamonds, int(SymbolDiamonds)}
	suitClubs    = &Suit{SymbolClubs, int(SymbolClubs)}
	suitSpades   = &Suit{SymbolSpades, int(SymbolSpades)}
	suitJoker    = &Suit{SymbolJoker, int(SymbolJoker)}
	suitUnknown  = &Suit{SymbolUnknown, int(SymbolUnknown)}

	stringAce   = "ace"
	string2     = "2"
	string3     = "3"
	string4     = "4"
	string5     = "5"
	string6     = "6"
	string7     = "7"
	string8     = "8"
	string9     = "9"
	string10    = "10"
	stringJack  = "jack"
	stringQueen = "queen"
	stringKing  = "king"

	heightAce     = "1"
	height2       = "2"
	height3       = "3"
	height4       = "4"
	height5       = "5"
	height6       = "6"
	height7       = "7"
	height8       = "8"
	height9       = "9"
	height10      = "0"
	heightJack    = "j"
	heightQueen   = "q"
	heightKing    = "k"
	heightJoker   = "!"
	heightUnknown = "?"

	faceAce     = &Face{HeightAce, int(HeightAce)}
	face2       = &Face{Height2, int(Height2)}
	face3       = &Face{Height3, int(Height3)}
	face4       = &Face{Height4, int(Height4)}
	face5       = &Face{Height5, int(Height5)}
	face6       = &Face{Height6, int(Height6)}
	face7       = &Face{Height7, int(Height7)}
	face8       = &Face{Height8, int(Height8)}
	face9       = &Face{Height9, int(Height9)}
	face10      = &Face{Height10, int(Height10)}
	faceJack    = &Face{HeightJack, int(HeightJack)}
	faceQueen   = &Face{HeightQueen, int(HeightQueen)}
	faceKing    = &Face{HeightKing, int(HeightKing)}
	faceJoker   = &Face{HeightJoker, int(HeightJoker)}
	faceUnknown = &Face{HeightUnknown, int(HeightUnknown)}

	faceAceCustomValue    = &Face{HeightAce, CUSTOM_VALUE}
	faceJokerCustomValue  = &Face{HeightJoker, CUSTOM_VALUE}
	suitHeartsCustomValue = &Suit{SymbolHearts, CUSTOM_VALUE}
	suitJokerCustomValue  = &Suit{SymbolJoker, CUSTOM_VALUE}

	cardAceHearts, _                 = New(stringAce, stringHearts, 0, 0, false)
	cardAceHeartsCustomFaceValues, _ = New(stringAce, stringHearts, CUSTOM_VALUE, 0, false)
	cardAceHeartsCustomSuitValues, _ = New(stringAce, stringHearts, 0, CUSTOM_VALUE, false)
	cardAceHeartsCustomValues, _     = New(stringAce, stringHearts, CUSTOM_VALUE, CUSTOM_VALUE, false)
	cardJoker, _                     = New(stringJoker, stringJoker, 0, 0, false)
	cardJokerForcedFace, _           = New(stringJoker, stringJoker, CUSTOM_VALUE, 0, false)
	cardJokerForcedSuit, _           = New(stringJoker, stringJoker, 0, CUSTOM_VALUE, false)
)
