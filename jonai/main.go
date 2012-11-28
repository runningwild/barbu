package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "strings"
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

func main() {
  flag.Parse()
  stdin := bufio.NewReader(os.Stdin)
  for {
    line, _, err := stdin.ReadLine()
    if err != nil {
      fmt.Errorf("ERROR: %v\n", err)
      return
    }
    var seating int
    _, err = fmt.Sscanf(string(line), "PLAYER %d", &seating)
    if err != nil {
      fmt.Errorf("ERROR: %v\n", err)
      return
    }

    line, _, err = stdin.ReadLine()
    if err != nil {
      fmt.Errorf("ERROR: %v\n", err)
      return
    }
    hand := strings.Fields(string(line))

    line, _, err = stdin.ReadLine()
    if err != nil {
      fmt.Errorf("ERROR: %v\n", err)
      return
    }
    switch string(line) {
    case "RAVAGE":
      StandardTrickTakingAi(stdin, seating, hand, NewRavageAi())

    case "LASTTWO":
      StandardTrickTakingAi(stdin, seating, hand, NewLastTwoAi())

    case "BARBU":
      StandardTrickTakingAi(stdin, seating, hand, NewBarbuAi())

    case "KILLERQUEENS":
      StandardTrickTakingAi(stdin, seating, hand, NewKillerQueensAi())

    case "DONE":
      return

    default:
      fmt.Errorf("UNKNOWN GAME '%s'\n", line)
      return
    }
  }
}
