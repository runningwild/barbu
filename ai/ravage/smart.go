package ravage

import (
  "bufio"
  "fmt"
  "github.com/runningwild/barbu/util"
  "os"
  "sort"
  "strings"
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

// func ravageRandomPlayer(input *bufio.Reader) {
//   // Read in hand
//   line, _, err := input.ReadLine()
//   if err != nil {
//     fmt.Fprintf(os.Stderr, "Error: %v\n")
//     return
//   }
//   cards := strings.Split(string(line), " ")

//   for len(cards) > 0 {
//     // Read in beginning of trick
//     line, _, err := input.ReadLine()
//     if err != nil {
//       fmt.Fprintf(os.Stderr, "Error: %v\n", err)
//       return
//     }
//     trick_start := strings.Split(string(line), " ")
//     play_index := -1
//     if len(line) > 0 {
//       suit := trick_start[0][1]
//       hits := 0
//       for i := range cards {
//         if cards[i][1] == suit {
//           hits++
//           if rand.Float64() <= 1/float64(hits) {
//             play_index = i
//           }
//         }
//       }
//     }
//     if play_index == -1 {
//       play_index = rand.Intn(len(cards))
//     }
//     fmt.Fprintf(os.Stdout, "%s\n", cards[play_index])
//     cards[play_index] = cards[len(cards)-1]
//     cards = cards[0 : len(cards)-1]

//     // Read in the rest of the trick
//     input.ReadLine()
//   }
// }

// // On a lead plays the highest rank in its largest suit
// // On a follow tries to duck, otherwise just plays its largest valid card
// func ravageStupidPlayer(input *bufio.Reader) {
//   // Read in hand
//   line, _, err := input.ReadLine()
//   if err != nil {
//     fmt.Fprintf(os.Stderr, "Error: %v\n")
//     return
//   }
//   cards := strings.Split(string(line), " ")
//   suits := make(map[byte][]string)
//   for i := range cards {
//     suits[cards[i][1]] = append(suits[cards[i][1]], cards[i])
//   }
//   for suit := range suits {
//     sort.Sort(handOfCards(suits[suit]))
//   }

//   for count := 0; count < 13; count++ {
//     // Read in beginning of trick
//     line, _, err := input.ReadLine()
//     if err != nil {
//       fmt.Fprintf(os.Stderr, "Error: %v\n", err)
//       return
//     }
//     trick_start := strings.Fields(string(line))
//     var play string
//     if len(trick_start) > 0 {
//       highest_rank := trick_start[0][0]
//       suit := trick_start[0][1]
//       for _, card := range trick_start {
//         if card[1] == suit && rank_map[card[0]] > rank_map[highest_rank] {
//           highest_rank = card[0]
//         }
//       }
//       cards := suits[suit]
//       if len(cards) > 0 {
//         if rank_map[cards[0][0]] > rank_map[highest_rank] {
//           play = cards[len(cards)-1]
//         } else {
//           var index int
//           for index = len(cards) - 1; index >= 0; index-- {
//             if rank_map[cards[index][0]] < rank_map[highest_rank] {
//               break
//             }
//           }
//           play = cards[index]
//           cards[index] = cards[len(cards)-1]
//         }
//         cards = cards[0 : len(cards)-1]
//         sort.Sort(handOfCards(cards))
//         suits[suit] = cards
//       }
//     }
//     if play == "" {
//       if len(trick_start) > 0 {
//         // We weren't able to follow suit, so play the highest rank card
//         // available
//         max_rank := 0
//         var suit, max_suit byte
//         var cards []string
//         for suit, cards = range suits {
//           if len(cards) == 0 {
//             continue
//           }
//           if rank_map[cards[len(cards)-1][0]] > max_rank {
//             max_rank = rank_map[cards[len(cards)-1][0]]
//             max_suit = suit
//           }
//         }
//         cards = suits[max_suit]
//         play = cards[len(cards)-1]
//         cards = cards[0 : len(cards)-1]
//         suits[max_suit] = cards
//       } else {
//         // We're leading, so play the lowest rank card from the suit that we
//         // have the least cards in.
//         min := 1000
//         var suit, min_suit byte
//         for suit = range suits {
//           if len(suits[suit]) < min && len(suits[suit]) > 0 {
//             min = len(suits[suit])
//             min_suit = suit
//           }
//         }
//         cards = suits[min_suit]
//         play = cards[0]
//         cards[0] = cards[len(cards)-1]
//         cards = cards[0 : len(cards)-1]
//         sort.Sort(handOfCards(cards))
//         suits[min_suit] = cards
//       }
//     }
//     fmt.Fprintf(os.Stdout, "%s\n", play)

//     // Read in the rest of the trick
//     input.ReadLine()
//   }
// }

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

func Smart(input *bufio.Reader) {
  // Read in hand
  line, _, err := input.ReadLine()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n")
    return
  }
  cards := strings.Split(string(line), " ")
  stats := util.MakeStats(cards)
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
    stats.TrickStart(trick_start)
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
        // First check if we can screw someone
        var suit, target_suit byte
        target := getLeader(trick_start)
        most := 0
        if len(trick_start) == 3 {
          for _, suit = range []byte{'c', 'd', 'h', 's'} {
            if len(suits[suit]) == 0 {
              continue
            }
            taken := stats.Taken(target, suit)
            if taken >= most {
              target_suit = suit
              most = taken
            }
          }
        }
        if target_suit == 0 {
          // We weren't able to follow suit, so play the highest rank card
          // from the suit that we have the least cards in.
          min := 1000
          for suit = range suits {
            if len(suits[suit]) < min && len(suits[suit]) > 0 {
              min = len(suits[suit])
              target_suit = suit
            }
          }
        }

        cards = suits[target_suit]
        play = cards[len(cards)-1]
        cards = cards[0 : len(cards)-1]
        sort.Sort(handOfCards(cards))
        suits[target_suit] = cards
      } else {
        // We're leading
        var best_suit byte
        var best_ratio float64 = 1.0
        for _, suit := range []byte{'c', 'd', 'h', 's'} {
          if len(suits[suit]) == 0 {
            continue
          }
          ratio := float64(len(suits[suit])) / (float64(stats.RemainingInSuit(suit) + len(suits[suit])))
          if ratio < best_ratio {
            best_ratio = ratio
            best_suit = suit
          }
        }
        cards = suits[best_suit]
        play = cards[0]
        cards[0] = cards[len(cards)-1]
        cards = cards[0 : len(cards)-1]
        sort.Sort(handOfCards(cards))
        suits[best_suit] = cards
      }
    }
    stats.TrickPlay(play)
    fmt.Fprintf(os.Stdout, "%s\n", play)

    // Read in the rest of the trick
    line, _, err = input.ReadLine()
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error: %v\n", err)
      return
    }
    trick_end := strings.Fields(string(line))
    stats.TrickEnd(trick_end)
  }
}

// func ravagemain() {
//   valid_modes := map[string]func(*bufio.Reader){
//     "random":  ravageRandomPlayer,
//     "stupid":  ravageStupidPlayer,
//     "smarter": ravageSmarterPlayer,
//   }
//   if f, ok := valid_modes[*mode]; ok {
//     buf := bufio.NewReader(os.Stdin)
//     for {
//       os.Stderr.Sync()
//       f(buf)
//       line, _, err := buf.ReadLine()
//       os.Stderr.Sync()
//       if string(line) != "RESET" || err != nil {
//         return
//       }
//     }
//   } else {
//     fmt.Fprintf(os.Stderr, "'%s' is not a valid mode, must be one of [", *mode)
//     var modes []string
//     for i := range valid_modes {
//       modes = append(modes, i)
//     }
//     sort.Strings(modes)
//     for i := range modes {
//       fmt.Fprintf(os.Stderr, " '%s'", modes[i])
//     }
//     fmt.Fprintf(os.Stderr, " ]\n")
//   }
// }
