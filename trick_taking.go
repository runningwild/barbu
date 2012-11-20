package main

import (
  "fmt"
  "github.com/runningwild/barbu/util"
  "strings"
)

type StandardDoubling struct {
  Players []Player
}

func stringSliceToIntSlice(ss []string) []int {
  var ns []int
  for _, s := range ss {
    var n int
    _, err := fmt.Sscanf(s, "%d", &n)
    if err == nil {
      // TODO: What to do if there is an error here?
      ns = append(ns, n)
    }
  }
  return ns
}

func (d *StandardDoubling) Double() [4][4]int {
  var ret [4][4]int
  // TODO: ERROR CHECKING!

  // Inform each player that doubling is starting
  for i := range d.Players {
    d.Players[i].Stdin().Write([]byte("DOUBLING\n"))
  }

  for i := range d.Players {
    cur := (i + 1) % 4
    // Tell player (i+1)%4 that it is their turn to double.
    d.Players[cur].Stdin().Write([]byte("DOUBLE\n"))
    line, _, _ := d.Players[cur].Stdout().ReadLine()
    doubles := stringSliceToIntSlice(strings.Fields(string(line)))
    for _, j := range doubles {
      ret[cur][j]++
      ret[j][cur]++
    }

    // TODO: Here is where we would check that the dealer only doubled players
    // that doubled him.

    // Tell each other player who that player doubled.
    for j := range d.Players {
      if j == cur {
        continue
      }
      d.Players[j].Stdin().Write([]byte(fmt.Sprintf("%d DOUBLE", cur)))
      for _, v := range doubles {
        d.Players[j].Stdin().Write([]byte(fmt.Sprintf(" %d", v)))
      }
      d.Players[j].Stdin().Write([]byte("\n"))
    }
  }

  return ret
}

type Trick struct {
  // Which players lead / took the trick.
  Lead, Took int

  // Mapping from player number to the card that player played this trick.
  Cards [4]string
}

type StandardTrickTaking struct {
  Players []Player
  Hands   []util.Hand
  Tricks  []Trick

  // Called after every trick if not nil.  If it returns true the game ends
  // immediately.
  End_early func([]Trick) bool
}

func (t *StandardTrickTaking) Run() {
  leader := 0
  for len(t.Hands[0]) > 0 && (t.End_early == nil || len(t.Tricks) == 0 || !t.End_early(t.Tricks)) {
    // Let each player know that we're starting a new trick
    for _, player := range t.Players {
      player.Stdin().Write([]byte("TRICK\n"))
    }

    var trick Trick
    trick.Lead = leader
    for i := range t.Players {
      cur := (i + leader) % 4
      // Tell the current player that it is their turn to play
      t.Players[cur].Stdin().Write([]byte("PLAY\n"))
      line, _, _ := t.Players[cur].Stdout().ReadLine()

      // TODO: check that line is properly formatted (i.e. exactly one card)
      card := strings.Fields(string(line))[0]
      t.Hands[cur].Remove(card)
      trick.Cards[cur] = card

      // Let all other players know what card this player played.
      for j := range t.Players {
        if j == cur {
          continue
        }
        t.Players[j].Stdin().Write([]byte(fmt.Sprintf("%d %s\n", cur, card)))
      }
    }

    // Note the winner
    trick.Took = trick.Lead
    high := trick.Cards[trick.Lead]
    for i := range t.Players {
      if trick.Cards[i][1] != high[1] {
        continue
      }
      if rank_map[trick.Cards[i][0]] > rank_map[high[0]] {
        high = trick.Cards[i]
        trick.Took = i
      }
    }

    leader = trick.Took
    t.Tricks = append(t.Tricks, trick)
  }
}
