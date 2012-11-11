package main

import (
  "fmt"
  "strings"
)

type Ravage struct {
  StandardTrickTakingGame
  players [4]Player
  hands   [4]map[card]bool
  tricks  [4][]card
  deck    Deck
  current int
}

func MakeRavage(players [4]Player, deck Deck) BarbuGame {
  var r Ravage
  r.players = players
  r.deck = deck
  return &r
}

func (r *Ravage) Deal() {
  hands := r.deck.Deal()
  for i := range r.players {
    r.hands[i] = make(map[card]bool)
    for j := range hands[i] {
      r.hands[i][hands[i][j]] = true
    }
    fmt.Fprintf(r.players[i].Stdin(), "%s\n", hands[i])
  }
  // for i := range r.players {
  //   fmt.Printf("Player %d: %s\n", i, hands[i])
  // }
}

func (r *Ravage) Score() [4]int {
  if len(r.hands[0]) > 0 {
    return [4]int{0, 0, 0, 0}
  }
  var maxes [4]map[byte]int
  for i := range maxes {
    maxes[i] = make(map[byte]int)
  }
  for i := range r.tricks {
    for _, card := range r.tricks[i] {
      maxes[i][card[1]]++
    }
  }
  max_max := 0
  for i := range maxes {
    if len(maxes[i]) > max_max {
      max_max = len(maxes[i])
    }
  }
  num_maxes := 0
  for i := range maxes {
    if len(maxes[i]) == max_max {
      num_maxes++
    }
  }
  screwedness := -36 / num_maxes
  scores := [4]int{0, 0, 0, 0}
  for i := range maxes {
    if len(maxes[i]) == max_max {
      scores[i] = screwedness
    }
  }
  return scores
}

func (r *Ravage) Round() bool {
  var trick_so_far string

  // Have everyone play their card for the trick
  for i := 0; i < 4; i++ {
    c := (r.current + i) % 4
    p := r.players[c]
    h := r.hands[c]
    // fmt.Printf("P%d: Send '%s'\n", c, trick_so_far)
    fmt.Fprintf(p.Stdin(), "%s\n", trick_so_far)
    bline, _, _ := p.Stdout().ReadLine()
    line := card(strings.TrimSpace(string(bline)))
    if !h[line] {
      fmt.Printf("ERROR: Player %d tried to player '%s' which he doesn't have!\n", c, line)
      // TODO: Handle error appropriately
    } else {
      delete(h, line)
    }
    // TODO: Check for valid play
    trick_so_far = strings.TrimSpace(trick_so_far + " " + string(line))
  }

  // Let everyone know what cards were played after them in this trick
  plays := strings.Split(trick_so_far, " ")
  for i := 0; i < 4; i++ {
    c := (r.current + i) % 4
    p := r.players[c]
    index := i + 1
    if index > 4 {
      index = 4
    }
    // fmt.Printf("P%d: Send '%s'\n", c, strings.Join(plays[index:], " "))
    fmt.Fprintf(p.Stdin(), "%s\n", strings.Join(plays[index:], " "))
  }
  suit := plays[0][1]
  index := 0
  for i := 0; i < 4; i++ {
    if plays[i][1] != suit {
      continue
    }
    if less(card(plays[index]), card(plays[i])) {
      index = i
    }
  }

  r.current = (r.current + index) % 4
  for i := range plays {
    r.tricks[r.current] = append(r.tricks[r.current], card(plays[i]))
  }

  return len(r.hands[0]) == 0
}
