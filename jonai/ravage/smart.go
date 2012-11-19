package ravage

import (
  "bufio"
  "fmt"
  "github.com/runningwild/barbu/util"
  "os"
  // "strings"
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

func Smart(input *bufio.Reader, seating int, shand []string) {
  hand := util.Hand(shand)

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
      switch seating {
      case 0:
        fmt.Printf("\n")
      case 1:
        fmt.Printf("0\n")
      case 2:
        fmt.Printf("0\n")
      case 3:
        fmt.Printf("0 1 2\n")
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
          card = hand[0]
        } else {
          for i := len(hand) - 1; i >= 0 && card == ""; i-- {
            if hand[i][1] == trick[0][1] {
              card = hand[i]
            }
          }
          if card == "" {
            // We only get here if we were unable to follow suit
            card = hand[len(hand)-1]
          }
        }
        hand.Remove(card)
        fmt.Printf("%s\n", card)
      } else {
        var player int
        var card string
        _, err := fmt.Sscanf(string(line), "%d %s", &player, &card)
        if err != nil {
          panic(err)
        }
        trick = append(trick, card)
      }
    }

  }
}
