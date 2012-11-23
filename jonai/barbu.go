package main

import (
// "github.com/runningwild/barbu/util"
)

type barbuAi struct {
  trickTakingBasics
}

func NewBarbuAi() trickTakingPlayer {
  return &barbuAi{}
}

func (r *barbuAi) Double(doubles [4][4]bool) []bool {
  return []bool{false, false, false}
}

func (r *barbuAi) holding() bool {
  for _, c := range r.hand {
    if c == "kh" {
      return true
    }
  }
  return false
}

func (r *barbuAi) pickLowestRatioInSuits(suits []byte) string {
  var card string
  lowest_ratio := 10000.0
  for _, suit := range suits {
    shand := r.hand.BySuit(suit)
    if shand.Len() == 0 {
      continue
    }
    ratio := float64(shand.Len()+1) / float64(r.stats.RemainingInSuit(suit)+1)
    if ratio < lowest_ratio {
      lowest_ratio = ratio
      card = shand[0]
    }
  }
  return card
}

func (r *barbuAi) Lead() string {
  var card string
  if r.holding() {
    card = r.pickLowestRatioInSuits([]byte{'c', 'd', 's'})
  }
  if card != "" {
    return card
  }
  card = r.pickLowestRatioInSuits([]byte{'c', 'd', 'h', 's'})
  return card
}

func (r *barbuAi) Follow(lead string) string {
  if lead == "ah" && r.holding() {
    return "kh"
  }

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

func (r *barbuAi) Discard() string {
  if r.holding() {
    return "kh"
  }
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
