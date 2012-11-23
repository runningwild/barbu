package barbu

import (
  "github.com/runningwild/barbu/util"
)

type Barbu struct {
  StandardDoubling
  StandardTrickTaking
}

func MakeBarbu(players []Player, hands [][]string) BarbuGame {
  var lt Barbu
  lt.StandardDoubling.Players = players[:]
  lt.StandardTrickTaking.Players = players[:]
  for _, hand := range hands {
    lt.StandardTrickTaking.Hands = append(lt.StandardTrickTaking.Hands, util.Hand(hand))
  }
  lt.End_early = func(tricks []Trick) bool {
    for _, card := range tricks[len(tricks)-1].Cards {
      if card == "kh" {
        return true
      }
    }
    return false
  }
  return &lt
}

func (lt *Barbu) Scores() [4]int {
  var ret [4]int
  for _, trick := range lt.Tricks {
    for _, card := range trick.Cards {
      if card == "kh" {
        ret[trick.Took] = -18
        return ret
      }
    }
  }
  panic("This should be unreachable if the game ran properly.")
}
