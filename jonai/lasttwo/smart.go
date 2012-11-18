package lasttwo

import (
  "bufio"
  "fmt"
  "os"
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

// func randomPlayer(input *bufio.Reader) {
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

// func stupidPlayer(input *bufio.Reader) {
//   // Read in hand
//   line, _, err := input.ReadLine()
//   if err != nil {
//     fmt.Fprintf(os.Stderr, "Error: %v\n")
//     return
//   }
//   cards := strings.Split(string(line), " ")

//   for count := 0; count < 13; count++ {
//     // Read in beginning of trick
//     line, _, err := input.ReadLine()
//     if err != nil {
//       fmt.Fprintf(os.Stderr, "Error: %v\n", err)
//       return
//     }
//     trick_start := strings.Fields(string(line))

//     play_index := -1
//     var suit byte = 255
//     if len(trick_start) > 0 {
//       suit = trick_start[0][1]
//     }
//     for play_index == -1 {
//       for i, card := range cards {
//         fmt.Errorf("card: '%s'\n", card)
//         if suit == 255 || card[1] == suit {
//           if play_index == -1 || rank_map[card[0]] < rank_map[cards[play_index][0]] {
//             play_index = i
//           }
//         }
//       }

//       // If we repeat this loop (because we didn't find a match) we will now
//       // relax the requirement on following suit because we were unable to do
//       // so.
//       suit = 255
//     }
//     fmt.Fprintf(os.Stdout, "%s\n", cards[play_index])
//     cards[play_index] = cards[len(cards)-1]
//     cards = cards[0 : len(cards)-1]

//     // Read in the rest of the trick
//     input.ReadLine()
//   }
// }

func Smart(stdin *bufio.Reader) {
  // Read in hand
  line, _, err := stdin.ReadLine()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n")
    return
  }
  cards := strings.Split(string(line), " ")

  for count := 0; count < 13; count++ {
    // Read in beginning of trick
    line, _, err := stdin.ReadLine()
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error: %v\n", err)
      return
    }
    trick_start := strings.Fields(string(line))

    play_index := -1
    var suit byte = 255
    if len(trick_start) > 0 {
      suit = trick_start[0][1]
    }
    for play_index == -1 {
      for i, card := range cards {
        fmt.Errorf("card: '%s'\n", card)
        if suit == 255 || card[1] == suit {
          if play_index == -1 || rank_map[card[0]] > rank_map[cards[play_index][0]] {
            play_index = i
          }
        }
      }

      // If we repeat this loop (because we didn't find a match) we will now
      // relax the requirement on following suit because we were unable to do
      // so.
      suit = 255
    }
    fmt.Fprintf(os.Stdout, "%s\n", cards[play_index])
    cards[play_index] = cards[len(cards)-1]
    cards = cards[0 : len(cards)-1]

    // Read in the rest of the trick
    stdin.ReadLine()
  }
}
