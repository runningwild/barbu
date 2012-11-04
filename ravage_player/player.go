package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "sort"
  "strings"
)

var mode = flag.String("mode", "random", "What ai to use")

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

type handOfCards []string

func (h handOfCards) Len() int {
  return len(h)
}
func (h handOfCards) Less(i, j int) bool {
  if h[i][1] != h[j][1] {
    return h[i][1] < h[j][1]
  }
  return rank_map[h[i][0]] < rank_map[h[j][0]]
}
func (h handOfCards) Swap(i, j int) {
  h[i], h[j] = h[j], h[i]
}

func randomPlayer() {
  input := bufio.NewReader(os.Stdin)
  // Read in hand
  line, _, err := input.ReadLine()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n")
    return
  }
  cards := strings.Split(string(line), " ")

  for len(cards) > 0 {
    // Read in beginning of trick
    line, _, err := input.ReadLine()
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error: %v\n", err)
      return
    }
    trick_start := strings.Split(string(line), " ")
    var play string
    if len(line) > 0 {
      suit := trick_start[0][1]
      for i := range cards {
        if cards[i][1] == suit {
          play = cards[i]
          cards[i] = cards[len(cards)-1]
          cards = cards[0 : len(cards)-1]
          break
        }
      }
    }
    if play == "" {
      play = cards[len(cards)-1]
      cards = cards[0 : len(cards)-1]
    }
    fmt.Fprintf(os.Stdout, "%s\n", play)

    // Read in the rest of the trick
    input.ReadLine()
  }
}

// On a lead plays the highest rank in its largest suit
// On a follow tries to duck, otherwise just plays its largest valid card
func stupidPlayer() {
  input := bufio.NewReader(os.Stdin)
  // Read in hand
  line, _, err := input.ReadLine()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n")
    return
  }
  cards := strings.Split(string(line), " ")
  suits := make(map[byte][]string)
  for i := range cards {
    suits[cards[i][1]] = append(suits[cards[i][1]], cards[i])
  }
  for suit := range suits {
    sort.Sort(handOfCards(suits[suit]))
  }

  for count := 0; count < 13; count++ {
    // Read in beginning of trick
    line, _, err := input.ReadLine()
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error: %v\n", err)
      return
    }
    trick_start := strings.Fields(string(line))
    var play string
    if len(trick_start) > 0 {
      highest_rank := trick_start[0][0]
      suit := trick_start[0][1]
      for _, card := range trick_start {
        if card[1] == suit && rank_map[card[0]] > rank_map[highest_rank] {
          highest_rank = card[0]
        }
      }
      cards := suits[suit]
      if len(cards) > 0 {
        if rank_map[cards[0][0]] > rank_map[highest_rank] {
          play = cards[len(cards)-1]
        } else {
          var index int
          for index = len(cards) - 1; index >= 0; index-- {
            if rank_map[cards[index][0]] < rank_map[highest_rank] {
              break
            }
          }
          play = cards[index]
          cards[index] = cards[len(cards)-1]
        }
        cards = cards[0 : len(cards)-1]
        sort.Sort(handOfCards(cards))
        suits[suit] = cards
      }
    }
    if play == "" {
      if len(trick_start) > 0 {
        // We weren't able to follow suit, so play the highest rank card
        // available
        max_rank := 0
        var suit, max_suit byte
        var cards []string
        for suit, cards = range suits {
          if len(cards) == 0 {
            continue
          }
          if rank_map[cards[len(cards)-1][0]] > max_rank {
            max_rank = rank_map[cards[len(cards)-1][0]]
            max_suit = suit
          }
        }
        cards = suits[max_suit]
        play = cards[len(cards)-1]
        cards = cards[0 : len(cards)-1]
        suits[max_suit] = cards
      } else {
        // We're leading, so play the lowest rank card from the suit that we
        // have the least cards in.
        min := 1000
        var suit, min_suit byte
        for suit = range suits {
          if len(suits[suit]) < min && len(suits[suit]) > 0 {
            min = len(suits[suit])
            min_suit = suit
          }
        }
        cards = suits[min_suit]
        play = cards[0]
        cards[0] = cards[len(cards)-1]
        cards = cards[0 : len(cards)-1]
        sort.Sort(handOfCards(cards))
        suits[min_suit] = cards
      }
    }
    fmt.Fprintf(os.Stdout, "%s\n", play)

    // Read in the rest of the trick
    input.ReadLine()
  }
}

type cardStats struct {
  // map from rank to set of cards that have been played
  by_rank map[byte]map[byte]bool

  // map from suit to set of cards that have been played
  by_suit map[byte]map[byte]bool

  // map from suit to number of cards played so far in that suit
  suit_count map[byte]int

  // map from suit to set of players void in that suit
  voids map[byte]map[int]bool
}

func makeCardStats() *cardStats {
  var cs cardStats

  cs.by_rank = make(map[byte]map[byte]bool)
  cs.by_suit = make(map[byte]map[byte]bool)
  cs.suit_count = make(map[byte]int)
  cs.voids = make(map[byte]map[int]bool)

  for rank := range rank_map {
    cs.by_rank[rank] = make(map[byte]bool)
  }

  for _, suit := range []byte{'h', 's', 'c', 'd'} {
    cs.by_suit[suit] = make(map[byte]bool)
    cs.voids[suit] = make(map[int]bool)
  }

  return &cs
}

func (cs *cardStats) Update(player int, card string) {
  rank := card[0]
  suit := card[1]
  cs.by_rank[rank][suit] = true
  cs.by_suit[suit][rank] = true
  cs.suit_count[suit]++
}

func (cs *cardStats) Trick(cards []string) {
  lead := cards[0][1]
  for i := range cards {
    if cards[i][1] != lead {
      cs.voids[lead][i] = true
    }
  }
}

var store *os.File

func init() {
  var err error
  store, err = os.Create("replay.txt")
  if err != nil {
    panic(err)
  }
}

func smarterPlayer() {
  defer store.Close()
  // cs := makeCardStats()
  input := bufio.NewReader(os.Stdin)
  // Read in hand
  line, _, err := input.ReadLine()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n")
    return
  }
  fmt.Fprintf(store, "Hand: %s\n", line)
  cards := strings.Split(string(line), " ")
  suits := make(map[byte][]string)
  for i := range cards {
    suits[cards[i][1]] = append(suits[cards[i][1]], cards[i])
  }
  for suit := range suits {
    sort.Sort(handOfCards(suits[suit]))
  }

  for count := 0; count < 13; count++ {
    // Read in beginning of trick
    line, _, err := input.ReadLine()
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error: %v\n", err)
      return
    }
    fmt.Fprintf(store, "Trick: %s\n", line)
    trick_start := strings.Fields(string(line))
    var play string
    if len(trick_start) > 0 {
      highest_rank := trick_start[0][0]
      suit := trick_start[0][1]
      for _, card := range trick_start {
        if card[1] == suit && rank_map[card[0]] > rank_map[highest_rank] {
          highest_rank = card[0]
        }
      }
      cards := suits[suit]
      if len(cards) > 0 {
        if rank_map[cards[0][0]] > rank_map[highest_rank] {
          play = cards[len(cards)-1]
        } else {
          var index int
          for index = len(cards) - 1; index >= 0; index-- {
            if rank_map[cards[index][0]] < rank_map[highest_rank] {
              break
            }
          }
          play = cards[index]
          cards[index] = cards[len(cards)-1]
        }
        cards = cards[0 : len(cards)-1]
        sort.Sort(handOfCards(cards))
        suits[suit] = cards
      }
    }
    if play == "" {
      if len(trick_start) > 0 {
        // We weren't able to follow suit, so play the highest rank card
        // available
        max_rank := 0
        var suit, max_suit byte
        var cards []string
        for suit, cards = range suits {
          if len(cards) == 0 || len(cards) >= 4 {
            continue
          }
          if rank_map[cards[len(cards)-1][0]] > max_rank {
            max_rank = rank_map[cards[len(cards)-1][0]]
            max_suit = suit
          }
        }
        cards = suits[max_suit]
        play = cards[len(cards)-1]
        cards = cards[0 : len(cards)-1]
        suits[max_suit] = cards
      } else {
        // We're leading, so play the lowest rank card from the suit that we
        // have the least cards in.
        min := 1000
        var suit, min_suit byte
        for suit = range suits {
          if len(suits[suit]) < min && len(suits[suit]) > 0 {
            min = len(suits[suit])
            min_suit = suit
          }
        }
        cards = suits[min_suit]
        play = cards[0]
        cards[0] = cards[len(cards)-1]
        cards = cards[0 : len(cards)-1]
        sort.Sort(handOfCards(cards))
        suits[min_suit] = cards
      }
    }
    fmt.Fprintf(os.Stdout, "%s\n", play)
    fmt.Fprintf(store, "Play: %s\n", play)
    // Read in the rest of the trick
    input.ReadLine()
  }
}

func main() {
  flag.Parse()
  valid_modes := map[string]func(){
    "random":  randomPlayer,
    "stupid":  stupidPlayer,
    "smarter": smarterPlayer,
  }
  if f, ok := valid_modes[*mode]; ok {
    f()
  } else {
    fmt.Fprintf(os.Stderr, "'%s' is not a valid mode, must be one of [", *mode)
    var modes []string
    for i := range valid_modes {
      modes = append(modes, i)
    }
    sort.Strings(modes)
    for i := range modes {
      fmt.Fprintf(os.Stderr, " '%s'", modes[i])
    }
    fmt.Fprintf(os.Stderr, " ]\n")
  }
}
