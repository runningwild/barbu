package main

import (
  "bufio"
  "fmt"
  "github.com/runningwild/barbu/ai/ravage"
  "os"
)

func main() {
  stdin := bufio.NewReader(os.Stdin)
  for {
    line, _, err := stdin.ReadLine()
    if err != nil {
      fmt.Errorf("ERROR: %v\n", err)
      return
    }
    switch string(line) {
    case "RAVAGE":
      ravage.Smart(stdin)

    case "LASTTWO":
      panic("NOT IMPLEMENTED")

    case "DONE":
      return

    default:
      fmt.Errorf("UNKNOWN GAME '%s'\n", line)
      return
    }
  }
}
