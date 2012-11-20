package ravage

import (
  "bufio"
  "flag"
  "fmt"
  "github.com/runningwild/barbu/util"
  "os"
  // "strings"
)

var P1 = flag.Int("p1", 0, "p1")
var P2 = flag.Int("p2", 0, "p2")
var P3 = flag.Int("p3", 0, "p3")

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

func getLeader(cards []string) int {
  suit := cards[0][1]
  best := -1
  for i, card := range cards {
    if card[1] != suit {
      continue
    }
    if best == -1 || rank_map[card[0]] > rank_map[cards[best][0]] {
      best = i
    }
  }
  return 3 - len(cards) + best
}

func goodEnoughToDouble(hand util.Hand) bool {
  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    subhand := hand.BySuit(suit)
    // Don't double if we have King and Ace in a single suit.
    if subhand.Len() >= 2 && rank_map[subhand[len(subhand)-2][0]] >= rank_map['k'] {
      return false
    }
  }
  return true
}

func lead(seat int, hand util.Hand, stats *util.Stats) string {
  var card string

  // Purposefully lead a trick with a high card if it is our only very high
  // card in that suit.
  diff := 0
  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    shand := hand.BySuit(suit)
    if shand.Len() <= 1 {
      continue
    }
    d := rank_map[shand[len(shand)-1][0]] - rank_map[shand[len(shand)-2][0]]
    if stats.Taken(seat, util.AnyRank, suit) == 0 && shand.Len() > 1 &&
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
    shand := hand.BySuit(suit)
    if shand.Len() == 0 {
      continue
    }
    ratio := float64(shand.Len()) / float64(stats.RemainingInSuit(suit))
    if ratio < lowest_ratio {
      lowest_ratio = ratio
      card = shand[0]
    }
  }
  return card
}

func follow(seat int, shand util.Hand, stats *util.Stats, trick []string) string {
  high := trick[0]
  for _, c := range trick {
    if c[1] == trick[0][1] && rank_map[c[0]] > rank_map[high[0]] {
      high = c
    }
  }
  if len(shand) == 1 {
    return shand[0]
  }

  // Trying taking tricks on purpose if you have a high card in that suit
  // if len(trick) == 1 && stats.Taken(seat, util.AnyRank, high[1]) <= *P1 {
  //   if rank_map[shand[len(shand)-1][0]]-rank_map[shand[len(shand)-2][0]] > *P2 {
  //     return shand[len(shand)-1]
  //   }
  // }

  if rank_map[shand[0][0]] > rank_map[high[0]] {
    if len(trick) == 3 {
      return shand[len(shand)-1]
    }
    return shand[len(shand)-1]
  }
  for i := len(shand) - 1; i >= 0; i-- {
    if rank_map[shand[i][0]] < rank_map[high[0]] {
      return shand[i]
    }
  }
  return shand[len(shand)-1]
}

func discard(hand util.Hand, stats *util.Stats) string {
  for _, suit := range []byte{'c', 'd', 'h', 's'} {
    shand := hand.BySuit(suit)
    if shand.Len() == 0 {
      continue
    }
    if shand.Len() == 1 && stats.RemainingInSuit(suit) > 1 {
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
    shand := hand.BySuit(suit)
    if shand.Len() == 0 {
      continue
    }
    ratio := float64(shand.Len()) / float64(stats.RemainingInSuit(suit))
    if ratio > best_ratio {
      best_ratio = ratio
      best_card = shand[len(shand)-1]
    }
  }
  return best_card
}

func Smart(input *bufio.Reader, seating int, shand []string) {
  hand := util.Hand(shand)
  stats := util.MakeStats(seating, hand)

  // Do all of the doubling
  line, _, err := input.ReadLine()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    return
  }
  if string(line) != "DOUBLING" {
    return
  }
  for i := 0; i < 4; i++ {
    line, _, err := input.ReadLine()
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error: %v\n", err)
      return
    }

    if string(line) == "DOUBLE" {
      if goodEnoughToDouble(hand) {
        switch seating {
        case 0:
          fmt.Printf("\n")
        case 1:
          fmt.Printf("0 2 3\n")
        case 2:
          fmt.Printf("0 1 3\n")
        case 3:
          fmt.Printf("0 1 2\n")
        }
      } else {
        fmt.Printf("\n")
      }
    }
    // TODO: Keep track of who doubled who so we can do something with it.
  }

  for {
    // Check that we're still going
    line, _, err := input.ReadLine()
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error: %v\n", err)
      return
    }
    if string(line) != "TRICK" {
      // This should be the end of the trick:
      score := make([]int, 4)
      _, err := fmt.Sscanf(
        string(line),
        "END %d %d %d %d",
        &score[0], &score[1], &score[2], &score[3])
      if err != nil {
        panic(err)
      }
      return
    }

    var trick []string
    for i := 0; i < 4; i++ {
      line, _, err := input.ReadLine()
      if err != nil {
        panic(err)
      }
      if string(line) == "PLAY" {
        // play something
        var card string
        if len(trick) == 0 {
          card = lead(seating, hand, stats)
        } else {
          shand := hand.BySuit(trick[0][1])
          if shand.Len() > 0 {
            card = follow(seating, shand, stats, trick)
          } else {
            // We only get here if we were unable to follow suit
            card = discard(hand, stats)
          }
        }
        stats.Played(seating, card)
        hand.Remove(card)
        fmt.Printf("%s\n", card)
      } else {
        var player int
        var card string
        _, err := fmt.Sscanf(string(line), "%d %s", &player, &card)
        if err != nil {
          panic(err)
        }
        stats.Played(player, card)
        trick = append(trick, card)
      }
    }

  }
}
