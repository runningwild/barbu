package main

// import (
//   "fmt"
//   "strings"
// )

// type LastTwo struct {
//   StandardTrickTakingGame
//   players        [4]Player
//   hands          [4]map[string]bool
//   tricks         [4][]string
//   deck           Deck
//   current        int
//   second_to_last int
//   last           int
// }

// func MakeLastTwo(players [4]Player, deck Deck) BarbuGame {
//   var bg LastTwo
//   bg.players = players
//   bg.deck = deck
//   return &bg
// }

// func (bg *LastTwo) Deal() {
//   hands := bg.deck.Deal()
//   for i := range bg.players {
//     fmt.Fprintf(bg.players[i].Stdin(), "LASTTWO\n")
//     bg.hands[i] = make(map[string]bool)
//     for j := range hands[i] {
//       bg.hands[i][hands[i][j]] = true
//     }
//     fmt.Fprintf(bg.players[i].Stdin(), "%s\n", hands[i])
//   }
//   // for i := range bg.players {
//   //   fmt.Printf("Player %d: %s\n", i, hands[i])
//   // }
// }

// func (bg *LastTwo) Score() [4]int {
//   var scores [4]int
//   scores[bg.second_to_last] -= 10
//   scores[bg.last] -= 20
//   return scores
// }

// func (bg *LastTwo) Round() bool {
//   var trick_so_far string

//   // Have everyone play their card for the trick
//   for i := 0; i < 4; i++ {
//     c := (bg.current + i) % 4
//     p := bg.players[c]
//     h := bg.hands[c]
//     // fmt.Printf("P%d: Send '%s'\n", c, trick_so_far)
//     fmt.Fprintf(p.Stdin(), "%s\n", trick_so_far)
//     bline, _, _ := p.Stdout().ReadLine()
//     line := card(strings.TrimSpace(string(bline)))
//     if !h[line] {
//       fmt.Printf("ERROR: Player %d tried to player '%s' which he doesn't have!\n", c, line)
//       // TODO: Handle error appropriately
//     } else {
//       delete(h, line)
//     }
//     // TODO: Check for valid play
//     trick_so_far = strings.TrimSpace(trick_so_far + " " + string(line))
//   }

//   // Let everyone know what cards were played after them in this trick
//   plays := strings.Split(trick_so_far, " ")
//   for i := 0; i < 4; i++ {
//     c := (bg.current + i) % 4
//     p := bg.players[c]
//     index := i + 1
//     if index > 4 {
//       index = 4
//     }
//     // fmt.Printf("P%d: Send '%s'\n", c, strings.Join(plays[index:], " "))
//     fmt.Fprintf(p.Stdin(), "%s\n", strings.Join(plays[index:], " "))
//   }
//   suit := plays[0][1]
//   index := 0
//   for i := 0; i < 4; i++ {
//     if plays[i][1] != suit {
//       continue
//     }
//     if less(card(plays[index]), card(plays[i])) {
//       index = i
//     }
//   }

//   bg.current = (bg.current + index) % 4
//   for i := range plays {
//     bg.tricks[bg.current] = append(bg.tricks[bg.current], card(plays[i]))
//   }
//   if len(bg.hands[0]) == 1 {
//     bg.second_to_last = bg.current
//   }
//   if len(bg.hands[0]) == 0 {
//     bg.last = bg.current
//   }

//   return len(bg.hands[0]) == 0
// }
