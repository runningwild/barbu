package main

import (
  "github.com/runningwild/barbu/util"
)

type lastTwoAi struct {
  trickTakingBasics
}

func NewLastTwoAi() trickTakingPlayer {
  return &lastTwoAi{}
}

func goodHand(hand util.Hand) bool {
  var low, med, high int
  lm := rank_map['5']
  mh := rank_map['9']
  for _, card := range hand {
    if rank_map[card[0]] <= lm {
      low++
    } else if rank_map[card[0]] <= mh {
      med++
    } else {
      high++
    }
  }
  return low > high && high > med
}

func (r *lastTwoAi) Double(doubles [4][4]bool) []bool {
  if goodHand(r.hand) {
    return []bool{true, true, true}
  }
  return []bool{false, false, false}
}

func (r *lastTwoAi) extreme(lowest bool) string {
  extreme := r.hand[0]
  for _, card := range r.hand {
    if lowest {
      if rank_map[card[0]] < rank_map[extreme[0]] {
        extreme = card
      }
    } else {
      if rank_map[card[0]] > rank_map[extreme[0]] {
        extreme = card
      }
    }
  }
  return extreme
}

func (r *lastTwoAi) Lead() string {
  if r.hand.Len() <= 2 {
    return r.extreme(true)
  }
  return r.extreme(false)
}

func (r *lastTwoAi) Follow(lead string) string {
  shand := r.hand.BySuit(lead[1])
  if r.hand.Len() <= 2 {
    return shand[0]
  }
  return shand[shand.Len()-1]
}

func (r *lastTwoAi) Discard() string {
  return r.extreme(false)
}
