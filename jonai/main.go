package main

import (
  "bufio"
  "flag"
  "fmt"
  "github.com/runningwild/barbu/jonai/lasttwo"
  "github.com/runningwild/barbu/jonai/ravage"
  "os"
  "strings"
)

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
      ravage.Smart(stdin, seating, hand)

    case "LASTTWO":
      lasttwo.Smart(stdin, seating, hand)

    case "DONE":
      return

    default:
      fmt.Errorf("UNKNOWN GAME '%s'\n", line)
      return
    }
  }
}
