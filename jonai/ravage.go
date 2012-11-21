package main

import (
  "flag"
  "github.com/runningwild/barbu/util"
)

var P1 = flag.Int("p1", 0, "p1")
var P2 = flag.Int("p2", 0, "p2")
var P3 = flag.Int("p3", 0, "p3")

type ravageAi struct {
  trickTakingBasics
}

func NewRavageAi() trickTakingPlayer {
  return &ravageAi{}
}

func (r *ravageAi) Double(doubles [4][4]bool) []bool {
  return []bool{false, false, false}
}

func (r *ravageAi) Lead() string {
  var card string
  // Purposefully lead a trick with a high card if it is our only very high
  // card in that suit.
  diff := 0
  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    shand := r.hand.BySuit(suit)
    if shand.Len() <= 1 {
      continue
    }
    d := rank_map[shand[len(shand)-1][0]] - rank_map[shand[len(shand)-2][0]]
    if r.stats.Taken(r.seat, util.AnyRank, suit) == 0 && shand.Len() > 1 &&
      d > 8 && d > diff {
      card = shand[len(shand)-1]
      diff = d
    }
  }

  if card != "" {
    return card
  }

  lowest_ratio := 10000.0
  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    shand := r.hand.BySuit(suit)
    if shand.Len() == 0 {
      continue
    }
    ratio := float64(shand.Len()) / float64(r.stats.RemainingInSuit(suit))
    if ratio < lowest_ratio {
      lowest_ratio = ratio
      card = shand[0]
    }
  }
  return card
}

func (r *ravageAi) Follow(lead string) string {
  shand := r.hand.BySuit(lead[1])
  for _, c := range r.trick {
    if c[1] == r.trick[0][1] && rank_map[c[0]] > rank_map[lead[0]] {
      lead = c
    }
  }
  if len(shand) == 1 {
    return shand[0]
  }

  // Trying taking tricks on purpose if you have a lead card in that suit
  // if len(trick) == 1 && stats.Taken(seat, util.AnyRank, lead[1]) <= *P1 {
  //   if rank_map[shand[len(shand)-1][0]]-rank_map[shand[len(shand)-2][0]] > *P2 {
  //     return shand[len(shand)-1]
  //   }
  // }

  if rank_map[shand[0][0]] > rank_map[lead[0]] {
    if len(r.trick) == 3 {
      return shand[len(shand)-1]
    }
    return shand[len(shand)-1]
  }
  for i := len(shand) - 1; i >= 0; i-- {
    if rank_map[shand[i][0]] < rank_map[lead[0]] {
      return shand[i]
    }
  }
  return shand[len(shand)-1]
}

func (r *ravageAi) Discard() string {
  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    shand := r.hand.BySuit(suit)
    if shand.Len() == 0 {
      continue
    }
    if shand.Len() == 1 && r.stats.RemainingInSuit(suit) > 1 {
      return shand[0]
    }
    if shand.Len() > 1 && rank_map[shand[len(shand)-1][0]]-rank_map[shand[len(shand)-2][0]] > 5 {
      return shand[len(shand)-1]
    }
    if shand.Len() > 2 &&
      rank_map[shand[len(shand)-1][0]]-rank_map[shand[len(shand)-2][0]] > 1 &&
      rank_map[shand[len(shand)-2][0]]-rank_map[shand[len(shand)-3][0]] > 4 {
      return shand[len(shand)-1]
    }
  }
  best_ratio := 0.0
  var best_card string
  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    shand := r.hand.BySuit(suit)
    if shand.Len() == 0 {
      continue
    }
    ratio := float64(shand.Len()) / float64(r.stats.RemainingInSuit(suit))
    if ratio > best_ratio {
      best_ratio = ratio
      best_card = shand[len(shand)-1]
    }
  }
  return best_card
}
