package main

import (
  "bufio"
  "flag"
  "fmt"
  "github.com/runningwild/barbu/util"
  "net"
  "strings"
)

var ip = flag.String("ip", "127.0.0.1", "Host's Ip address.")
var player = flag.Int("player", -1, "Which player [0-3] you expect to be.")

func main() {
  flag.Parse()

  if !(*player >= 0 && *player <= 3) {
    fmt.Printf("player must be in [0, 3].\n")
    return
  }

  raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", *ip, 5200+*player))
  if err != nil {
    fmt.Printf("Unable to resolve tcp addr: %v\n", err)
    return
  }

  conn, err := net.DialTCP("tcp", nil, raddr)
  if err != nil {
    fmt.Printf("Unable to dial host: %v\n", err)
    return
  }
  read := bufio.NewReader(conn)

  for {
    line, _, err := read.ReadLine()
    if err != nil {
      panic(err)
    }
    fmt.Printf("Game is %s\n", string(line))

    line, _, err = read.ReadLine()
    if err != nil {
      panic(err)
    }
    hand := util.Hand(strings.Fields(string(line)))

    for i := 0; i < 13; i++ {
      hand.Sort()
      fmt.Printf("Hand is %v\n", hand)
      fmt.Printf("\n")

      line, _, err = read.ReadLine()
      if err != nil {
        panic(err)
      }
      fmt.Printf("Played before you this round: %s\n", string(line))

      var card string
      fmt.Scanf("%s", &card)
      hand.Remove(card)
      _, err = conn.Write([]byte(card))
      if err != nil {
        panic(err)
      }

      line, _, err = read.ReadLine()
      if err != nil {
        panic(err)
      }
      fmt.Printf("Played after you this round: %s\n", string(line))

      fmt.Printf("\n\n")
    }
  }
}
