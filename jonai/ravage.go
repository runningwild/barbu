package main

import (
  "github.com/runningwild/barbu/util"
)

type ravageAi struct {
  trickTakingBasics
}

func NewRavageAi() trickTakingPlayer {
  return &ravageAi{}
}

func (r *ravageAi) Double(doubles [4][4]bool) []bool {
  return []bool{false, false, false}
}

func (r *ravageAi) winner(trick []string) int {
  suit := trick[0][1]
  var index int
  high := -1
  for i, card := range trick {
    if card[1] != suit {
      continue
    }
    if rank_map[card[0]] > high {
      high = rank_map[card[0]]
      index = i
    }
  }
  return (r.seat - len(trick) + index + 4) % 4
}

func (r *ravageAi) Lead() string {
  var card string
  // Purposefully lead a trick with a high card if it is our only very high
  // card in that suit.
  diff := 8
  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    shand := r.hand.BySuit(suit)
    if shand.Len() <= 1 || r.stats.Taken(r.seat, util.AnyRank, suit) > 0 {
      continue
    }
    A := rank_map[shand[len(shand)-1][0]]
    B := rank_map[shand[len(shand)-2][0]]
    cdiff := A - B
    if cdiff > diff {
      card = shand[len(shand)-1]
      diff = cdiff
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
  {
    var card string
    winner := r.winner(r.trick)
    biggest := 6

    for _, suit := range []byte{'c', 'd', 'h', 's'} {
      shand := r.hand.BySuit(suit)
      if shand.Len() == 0 {
        continue
      }

      if r.stats.Taken(winner, util.AnyRank, suit) > biggest {
        biggest = r.stats.Taken(winner, util.AnyRank, suit)
        card = shand[len(shand)-1]
      }
    }
    if card != "" {
      return card
    }
  }

  // If there is a big difference between the two highest cards in one suit,
  // then ditch the highest.  A card is more likely to be discarded in this
  // fashion the higher it is and the bigger the gap between it and the next
  // highest card.
  {
    diff := 16
    var card string
    for _, suit := range []byte{'c', 'd', 'h', 's'} {
      shand := r.hand.BySuit(suit)
      if shand.Len() <= 1 {
        continue
      }
      A := rank_map[shand[len(shand)-1][0]]
      B := rank_map[shand[len(shand)-2][0]]
      cdiff := A*A - B*B
      if cdiff > diff {
        card = shand[len(shand)-1]
        diff = cdiff
      }
    }
    if card != "" {
      return card
    }
  }

  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    shand := r.hand.BySuit(suit)
    if shand.Len() == 0 {
      continue
    }
    if shand.Len() == 1 && r.stats.RemainingInSuit(suit) > 1 {
      return shand[0]
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
