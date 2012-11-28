package util

import (
  "fmt"
  "sort"
)

var rank_map map[byte]int

func init() {
  rank_map = map[byte]int{
    '2': 1,
    '3': 2,
    '4': 3,
    '5': 4,
    '6': 5,
    '7': 6,
    '8': 7,
    '9': 8,
    't': 9,
    'j': 10,
    'q': 11,
    'k': 12,
    'a': 13,
  }
}

type Hand []string

func (h Hand) Len() int {
  return len(h)
}
func (h Hand) Less(i, j int) bool {
  if h[i][1] != h[j][1] {
    return h[i][1] < h[j][1]
  }
  return rank_map[h[i][0]] < rank_map[h[j][0]]
}
func (h Hand) Swap(i, j int) {
  h[i], h[j] = h[j], h[i]
}
func (h *Hand) Remove(card string) {
  for i := range *h {
    if (*h)[i] == card {
      (*h)[i] = (*h)[len((*h))-1]
      (*h) = (*h)[0 : len((*h))-1]
      sort.Sort((*h))
      return
    }
  }
  panic(fmt.Sprintf("Didn't find '%s' in the hand '%v'.\n", card, h))
}
func (h Hand) Sort() {
  sort.Sort(h)
}
func (h Hand) BySuit(suit byte) Hand {
  var h2 Hand
  for _, card := range h {
    if card[1] == suit {
      h2 = append(h2, card)
    }
  }
  h2.Sort()
  return h2
}

type Stats struct {
  // Seating of the player these stats are for
  seating int

  // map from [player, suit] to whether that player is void in that suit
  voids map[int]map[byte]bool

  // List of cards remaining in other players' hands
  remaining_cards map[string]bool

  // map from player to all of the cards that player has taken
  taken map[int][]string

  // List of cards remaining in each suit
  remaining_suits map[byte]int

  trick struct {
    // Suit lead this trick
    lead byte

    // Number of cards played so far
    played int

    cards [4]string
  }
}

// Creates a stats object with the knowledge of your hand.
func MakeStats(seating int, hand []string) *Stats {
  var s Stats
  s.seating = seating
  s.voids = make(map[int]map[byte]bool)
  for i := 0; i <= 3; i++ {
    s.voids[i] = make(map[byte]bool)
  }

  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    count := 0
    for i := range hand {
      if hand[i][1] == suit {
        count++
      }
    }
    if count == 0 {
      s.voids[s.seating][suit] = true
    }
  }

  s.remaining_suits = map[byte]int{
    'c': 13,
    'd': 13,
    'h': 13,
    's': 13,
  }

  s.taken = make(map[int][]string)

  s.remaining_cards = make(map[string]bool)
  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    for _, rank := range []byte{'2', '3', '4', '5', '6', '7', '8', '9', 'j', 'q', 'k', 'a'} {
      s.remaining_cards[string([]byte{rank, suit})] = true
    }
  }
  for _, card := range hand {
    s.remaining_cards[card] = false
  }

  return &s
}

func (s *Stats) update(player int, card string) {
  if s.trick.played%4 == 0 {
    s.trick.played = 0
    s.trick.lead = card[1]
  }
  s.trick.cards[player] = card
  s.trick.played++

  // This player didn't follow suit - they must be void
  if card[1] != s.trick.lead {
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

func (s *Stats) Played(player int, card string) {
  s.update(player, card)
}

// Returns the number of cards of the specified suit that a player has taken.
const AnyPlayer = 200
const AnySuit = 200
const AnyRank = 200

func (s *Stats) Taken(player int, rank, suit byte) int {
  count := 0
  for the_player, taken := range s.taken {
    if player != the_player && player != AnyPlayer {
      continue
    }
    for _, card := range taken {
      if rank != card[0] && rank != AnyRank {
        continue
      }
      if suit != card[1] && suit != AnySuit {
        continue
      }
      count++
    }
  }
  return count
}

func (s *Stats) Remaining(card string) bool {
  return s.remaining_cards[card]
}

func (s *Stats) RemainingInSuit(suit byte) int {
  return s.remaining_suits[suit]
}

func (s *Stats) RemainingAboveInSuit(rank, suit byte) int {
  count := 0
  for card := range s.remaining_cards {
    if card[1] == suit && rank_map[card[0]] > rank_map[rank] {
      count++
    }
  }
  return count
}

func (s *Stats) IsDefinitelyVoid(player int, suit byte) bool {
  return s.voids[player][suit]
}
