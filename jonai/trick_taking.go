package main

import (
  "bufio"
  "fmt"
  "github.com/runningwild/barbu/util"
  "os"
  "strings"
)

type trickTakingBasics struct {
  seat  int
  hand  util.Hand
  stats *util.Stats
  trick []string
}

func (t *trickTakingBasics) Begin(seat int, hand []string) {
  for _, card := range hand {
    t.hand = append(t.hand, card)
  }
  t.seat = seat
  t.stats = util.MakeStats(seat, t.hand)
}
func (t *trickTakingBasics) Hand() util.Hand {
  return t.hand
}
func (t *trickTakingBasics) PlayCard(seat int, card string) {
  t.stats.Played(seat, card)
  if seat == t.seat {
    t.hand.Remove(card)
  }
  t.trick = append(t.trick, card)
  if len(t.trick) == 4 {
    t.trick = t.trick[0:0]
  }
}

type trickTakingPlayer interface {
  Begin(seat int, hand []string)
  Hand() util.Hand
  PlayCard(seat int, card string)

  // if double[i][j] == true then player i doubled player j.  Note that some
  // players may not have gotten to double yet and so some (or all) of these
  // values might be false even though players might still double.
  Double(doubles [4][4]bool) []bool

  Lead() string
  Follow(lead string) string
  Discard() string
}

func StandardTrickTakingAi(input *bufio.Reader, seating int, shand []string, player trickTakingPlayer) {
  player.Begin(seating, shand)

  // Do all of the doubling
  line, _, err := input.ReadLine()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    return
  }
  if string(line) != "DOUBLING" {
    return
  }

  var doubles [4][4]bool
  for i := 0; i < 4; i++ {
    line, _, err := input.ReadLine()
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error: %v\n", err)
      return
    }

    fields := strings.Fields(string(line))
    if fields[0] == "DOUBLE" {
      double := player.Double(doubles)
      for i := range double {
        if double[i] {
          fmt.Printf("%d ", i)
        }
      }
      fmt.Printf("\n")
    } else {
      var player int
      fmt.Sscanf(fields[0], "%d", &player)
      for i := 2; i < len(fields); i++ {
        var target int
        fmt.Sscanf(fields[i], "%d", &target)
        doubles[player][target] = true
      }
    }
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
          card = player.Lead()
        } else {
          shand := player.Hand().BySuit(trick[0][1])
          if shand.Len() > 0 {
            card = player.Follow(trick[0])
          } else {
            // We only get here if we were unable to follow suit
            card = player.Discard()
          }
        }
        player.PlayCard(seating, card)
        fmt.Printf("%s\n", card)
      } else {
        var player_pos int
        var card string
        _, err := fmt.Sscanf(string(line), "%d %s", &player_pos, &card)
        if err != nil {
          panic(err)
        }
        player.PlayCard(player_pos, card)
        trick = append(trick, card)
      }
    }
  }
}
