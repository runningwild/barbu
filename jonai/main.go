package main

import (
  "bufio"
  "flag"
  "fmt"
  "github.com/runningwild/barbu/jonai/lasttwo"
  "github.com/runningwild/barbu/jonai/ravage"
  "os"
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
    switch string(line) {
    case "RAVAGE":
      ravage.Smart(stdin)

    case "LASTTWO":
      lasttwo.Smart(stdin)

    case "DONE":
      return

    default:
      fmt.Errorf("UNKNOWN GAME '%s'\n", line)
      return
    }
  }
}
