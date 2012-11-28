package main

type killerQueensAi struct {
  trickTakingBasics
}

func NewKillerQueensAi() trickTakingPlayer {
  return &killerQueensAi{}
}

func (r *killerQueensAi) Double(doubles [4][4]bool) []bool {
  return []bool{false, false, false}
}

func (r *killerQueensAi) Lead() string {
  var card string

  lowest_ratio := 10000.0
  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    shand := r.hand.BySuit(suit)
    if shand.Len() == 0 {
      continue
    }

    // If we have the queen and someone has a card higher than the queen, and
    // no one has a card less than the queen except for us, then we can play
    // the queen with impunity.
    queen := string([]byte{'q', suit})
    have := make(map[byte]bool)
    for _, card := range shand {
      have[card[0]] = true
    }
    if have['q'] {
      targets := 0
      if !have['k'] && r.stats.Remaining(string([]byte{'k', suit})) {
        targets++
      }
      if !have['a'] && r.stats.Remaining(string([]byte{'a', suit})) {
        targets++
      }
      if targets > 0 && r.stats.RemainingInSuit(suit) == shand.Len()+targets {
        return queen
      }
    }

    if (shand[0][0] == 'k' || shand[0][0] == 'a') && r.stats.Remaining(queen) {
      continue
    }
    if shand[0][0] == 'q' {
      continue
    }
    ratio := float64(shand.Len()) / float64(r.stats.RemainingInSuit(suit))
    if ratio < lowest_ratio {
      lowest_ratio = ratio
      card = shand[0]
    }
  }
  if card != "" {
    return card
  }

  lowest := 1000
  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    shand := r.hand.BySuit(suit)
    if shand.Len() == 0 {
      continue
    }
    val := shand.Len()
    if shand[0][0] == 'q' && shand.Len() == 1 {
      val += 100
    }
    if val < lowest {
      lowest = val
      card = shand[0]
      if card[0] == 'q' && shand.Len() > 1 {
        card = shand[1]
      }
    }
  }
  return card
}

func (r *killerQueensAi) Follow(lead string) string {
  shand := r.hand.BySuit(lead[1])
  if shand.Len() == 1 {
    return shand[0]
  }
  have_queen := false
  for _, card := range shand {
    if card[0] == 'q' {
      have_queen = true
    }
  }
  if have_queen && (lead[0] == 'k' || lead[0] == 'a') {
    return string([]byte{'q', lead[1]})
  }

  highest := shand[0]
  if highest[0] == 'q' {
    highest = shand[1]
  }
  for _, card := range shand {
    if card[0] == 'q' {
      continue
    }
    if rank_map[card[0]] < rank_map[lead[0]] {
      highest = card
    }
  }

  return highest
}

func (r *killerQueensAi) Discard() string {
  var card string
  lowest_ratio := 10000.0
  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    shand := r.hand.BySuit(suit)
    for _, c := range shand {
      if c[0] != 'q' {
        continue
      }
      if shand.Len() == 1 && r.stats.RemainingInSuit(suit) == 1 {
        return c
      }
      ratio := float64(shand.Len()) / float64(r.stats.RemainingInSuit(suit))
      if ratio > lowest_ratio {
        continue
      }
      ratio = lowest_ratio
      card = c
    }
  }
  if card != "" {
    return card
  }

  lowest_ratio = 10000.0
  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    shand := r.hand.BySuit(suit)
    for _, c := range shand {
      ratio := float64(shand.Len()) / float64(r.stats.RemainingInSuit(suit))
      if ratio > lowest_ratio {
        continue
      }
      ratio = lowest_ratio
      card = c
    }
  }
  return card
}
