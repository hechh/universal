package test

import (
	"poker_server/server/room/internal/internal/texas"
	"poker_server/server/room/internal/module/card"
	"testing"
)

func TestCard(t *testing.T) {
	card1 := card.StrToCard("51", "k2", "q1", "j1", "a1", "23", "k2")
	card2 := card.StrToCard("61", "k2", "q1", "j1", "a1", "23", "k2")

	t1, v1, b1 := texas.Compare(card1...)
	t2, v2, b2 := texas.Compare(card2...)
	t.Log(t1, v1, card.CardList(b1).String())
	t.Log(t2, v2, card.CardList(b2).String())
}
