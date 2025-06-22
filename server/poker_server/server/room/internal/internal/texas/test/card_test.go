package test

import (
	"poker_server/server/room/internal/internal/texas"
	"poker_server/server/room/internal/module/card"
	"testing"
)

func TestCard(t *testing.T) {
	card1 := card.StrToCard("a1", "21", "32", "43", "52", "81", "73")
	card2 := card.StrToCard("21", "61", "32", "43", "52", "81", "73")
	//card2 := card.StrToCard("41", "51", "22", "33", "k2", "k1", "k3")

	t1, v1, b1 := texas.Compare(card1...)
	t2, v2, b2 := texas.Compare(card2...)
	t.Log(t1, v1, card.CardList(b1).String())
	t.Log(t2, v2, card.CardList(b2).String())
}
