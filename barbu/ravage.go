package barbu

import (
  "github.com/runningwild/barbu/util"
)

type Ravage struct {
  StandardDoubling
  StandardTrickTaking
}

func makeRavage(players []Player, hands [][]string) BarbuGame {
  var r Ravage
  r.StandardDoubling.Players = players[:]
  r.StandardTrickTaking.Players = players[:]
  for _, hand := range hands {
    r.StandardTrickTaking.Hands = append(r.StandardTrickTaking.Hands, util.Hand(hand))
  }
  return &r
}
func init() {
  RegisterBarbuGame("ravage", makeRavage)
}

func (r *Ravage) Scores() [4]int {
  player_suit_count := make([]map[byte]int, 4)
  for i := range player_suit_count {
    player_suit_count[i] = make(map[byte]int)
  }
  for _, trick := range r.StandardTrickTaking.Tricks {
    for _, card := range trick.Cards {
      player_suit_count[trick.Took][card[1]]++
    }
  }

  player_maxes := make([]int, 4)
  max := 0
  num_maxes := 0
  for i := range player_suit_count {
    for _, count := range player_suit_count[i] {
      if count > player_maxes[i] {
        player_maxes[i] = count
      }
    }
    if player_maxes[i] == max {
      num_maxes++
    }
    if player_maxes[i] > max {
      max = player_maxes[i]
      num_maxes = 1
    }
  }

  var ret [4]int
  for i := range player_maxes {
    if player_maxes[i] == max {
      ret[i] = -36 / num_maxes
    }
  }
  return ret
}
