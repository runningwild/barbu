package barbu

import (
  "github.com/runningwild/barbu/util"
)

type KillerQueens struct {
  StandardDoubling
  StandardTrickTaking
}

func makeKillerQueens(players []Player, hands [][]string) BarbuGame {
  var kq KillerQueens
  kq.StandardDoubling.Players = players[:]
  kq.StandardTrickTaking.Players = players[:]
  for _, hand := range hands {
    kq.StandardTrickTaking.Hands = append(kq.StandardTrickTaking.Hands, util.Hand(hand))
  }
  kq.End_early = func(tricks []Trick) bool {
    count := 0
    for _, trick := range tricks {
      for _, card := range trick.Cards {
        if card[0] == 'q' {
          count++
        }
      }
    }
    return count == 4
  }
  return &kq
}
func init() {
  RegisterBarbuGame("killerqueens", makeKillerQueens)
}

func (kq *KillerQueens) Scores() [4]int {
  var ret [4]int
  for _, trick := range kq.StandardTrickTaking.Tricks {
    for _, card := range trick.Cards {
      if card[0] == 'q' {
        ret[trick.Took] -= 6
      }
    }
  }
  return ret
}
