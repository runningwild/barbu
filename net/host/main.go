package main

import (
  "bufio"
  "encoding/gob"
  "flag"
  "fmt"
  "github.com/runningwild/barbu/barbu"
  base "github.com/runningwild/barbu/net"
  "io"
  "net"
)

var addr = flag.String("addr", "", "Address of the server.")
var port = flag.Int("port", 9901, "Port to connect to server on.")
var name = flag.String("name", "", "Name of the game to join.")
var game = flag.String("game", "", "Name of the barbu sub-game to play.")

type subConnPlayer struct {
  seat   int
  rw     io.ReadWriter
  stdout *bufio.Reader
}

func makeSubConnPlayer(rw io.ReadWriter, seat int) barbu.Player {
  var p subConnPlayer
  p.seat = seat
  p.rw = rw
  p.stdout = bufio.NewReader(p.rw)
  return &p
}
func (p *subConnPlayer) Id() string {
  return fmt.Sprintf("%d", p.seat)
}
func (p *subConnPlayer) Stdin() io.Writer {
  return p.rw
}
func (p *subConnPlayer) Stdout() *bufio.Reader {
  return p.stdout
}
func (p *subConnPlayer) Stderr() *bufio.Reader {
  return bufio.NewReader(nil)
}
func (p *subConnPlayer) Close() {
}

func makePlayers(conn net.Conn) []barbu.Player {
  splits := splitConn(conn)
  var players []barbu.Player
  for seat, sub := range splits {
    players = append(players, makeSubConnPlayer(sub, seat))
  }
  return players
}

type subConn struct {
  in, out   chan []byte
  remaining []byte
}

func (c *subConn) transferRemaining(buf []byte) int {
  if len(buf) >= len(c.remaining) {
    copy(buf, c.remaining)
    r := len(c.remaining)
    c.remaining = nil
    return r
  }
  // else
  copy(buf, c.remaining)
  c.remaining = c.remaining[len(buf):]
  return len(buf)
}

func (c *subConn) Read(buf []byte) (int, error) {
  if len(c.remaining) == 0 {
    c.remaining = <-c.in
  }
  n := c.transferRemaining(buf)
  // fmt.Printf("READ: '%s'\n", buf)
  return n, nil
}

func (c *subConn) Write(buf []byte) (int, error) {
  // fmt.Printf("WRITE: '%s'\n", buf)
  b := make([]byte, len(buf))
  copy(b, buf)
  c.out <- b
  return len(b), nil
}

func splitConn(conn net.Conn) []io.ReadWriter {
  subs := make([]subConn, 4)
  for i := range subs {
    subs[i].in = make(chan []byte, 10)
    subs[i].out = make(chan []byte, 10)
  }

  go func() {
    dec := gob.NewDecoder(conn)
    for {
      var rd base.RemoteData
      err := dec.Decode(&rd)
      if err != nil {
        fmt.Printf("Error reading from server: %v\n", err)
        return
      }
      subs[rd.Client].in <- []byte(rd.Line)
    }
  }()

  collect := make(chan base.RemoteData, 10)
  for i := range subs {
    go func(n int, sc *subConn) {
      for data := range sc.out {
        collect <- base.RemoteData{n, string(data)}
      }
    }(i, &subs[i])
  }

  enc := gob.NewEncoder(conn)
  go func() {
    for data := range collect {
      err := enc.Encode(data)
      fmt.Printf("Sent(%d): '%s'\n", data.Client, data.Line)
      if err != nil {
        fmt.Printf("Error writing data from collector: %v\n", err)
        return
      }
    }
  }()

  ret := make([]io.ReadWriter, len(subs))
  for i := range ret {
    ret[i] = &subs[i]
  }
  return ret
}

func connectAsHost(addr string, port int, name, game string) {
  raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, port))
  if err != nil {
    fmt.Printf("Unable to resolve server address: %v\n", err)
    return
  }

  conn, err := net.DialTCP("tcp", nil, raddr)
  if err != nil {
    fmt.Printf("Unable to connect to server: %v\n", err)
    return
  }
  defer conn.Close()
  fmt.Printf("A\n")
  _, err = conn.Write([]byte(fmt.Sprintf("%s\n", name)))
  if err != nil {
    fmt.Printf("Failed to host game: %v\n", err)
    return
  }
  fmt.Printf("B\n")
  players := makePlayers(conn)
  err = barbu.RunGames(players, 0, game, 1, false)
  if err != nil {
    fmt.Printf("Error running games: %v\n", err)
    return
  }
  fmt.Printf("D\n")
  // players[0].Stdin().Write([]byte("RAWR!!!"))
  fmt.Printf("E\n")
  return
}

func main() {
  flag.Parse()
  if *addr == "" {
    fmt.Printf("Must specify an address with --addr.\n")
    return
  }
  if *name == "" {
    fmt.Printf("Must specify a game name with --name.\n")
    return
  }
  if *game == "" {
    fmt.Printf("Must specify a barbu sub-game with --game.\n")
    return
  }
  connectAsHost(*addr, *port, *name, *game)
}
