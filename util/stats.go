package util

type Stats struct {
  // map from [player, suit] to whether that player is void in that suit
  voids map[int]map[byte]bool

  // List of cards remaining in other players' hands
  remaining_cards map[string]bool

  // List of cards remaining in each suit
  remaining_suits map[byte]int

  trick struct {
    // Suit lead this trick
    lead byte

    // Number of cards played so far
    played int
  }
}

// Creates a stats object with the knowledge of your hand.
// All arrays are three elements: {left, across, right}
func MakeStats(hand []string) Stats {
  return Stats{}
}

// player: [0, 1, 2, 3] == [left, across, right, self]
func (s *Stats) update(player int, card string) {
  if s.trick.played == 0 {
    s.trick.lead = card[1]
  }
  s.trick.played++

  // This player didn't follow suit - they must be void
  if player != 3 && card[1] != s.trick.lead {
    s.voids[player][s.trick.lead] = true
  }

  delete(s.remaining_cards, card)
  s.remaining_suits[card[1]]--
  if s.remaining_suits[card[1]] == 0 {
    s.voids[0][card[1]] = true
    s.voids[1][card[1]] = true
    s.voids[2][card[1]] = true
  }
}

// Updates Stats with cards played before you this trick, which may be empty.
func (s *Stats) TrickStart(cards []string) {
  s.trick.played = 0
  for i, card := range cards {
    player := i + (3 - len(cards))
    s.update(player, card)
  }
}

// Updates Stats with the card you played this trick
func (s *Stats) TrickPlay(card string) {
  s.update(3, card)
}

// Updates Stats with cards played after you this trick, which may be empty.
func (s *Stats) TrickEnd(cards []string) {
  for i, card := range cards {
    s.update(i, card)
  }
}
