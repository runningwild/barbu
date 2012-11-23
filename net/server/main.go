package main

import (
  "bufio"
  "flag"
  "fmt"
  "net"
  "sync"
  "time"
)

var host_port = flag.Int("host_port", 9901, "TCP port to listen for hosts on.")
var client_port = flag.Int("client_port", 9902, "TCP port to listen for clients on.")

type connSeat struct {
  seat int
  conn net.Conn
}

type Game struct {
  seat struct {
    take  chan connSeat
    got   chan bool
    taken [4]net.Conn
  }
}

func startGame() *Game {
  var g Game
  g.seat.take = make(chan connSeat)
  g.seat.got = make(chan bool)
  go g.routine()
  return &g
}
func (g *Game) TakeSeat(seat int, conn net.Conn) bool {
  g.seat.take <- connSeat{seat, conn}
  return <-g.seat.got
}
func (g *Game) routine() {
  for {
    select {
    case req := <-g.seat.take:
      if g.seat.taken[req.seat] != nil {
        g.seat.got <- false
      } else {
        g.seat.taken[req.seat] = req.conn
        g.seat.got <- true
      }
    }
  }
}

var active_games_mutex sync.Mutex
var active_games map[string]*Game

func init() {
  active_games = make(map[string]*Game)
}

func registerGame(name string) bool {
  active_games_mutex.Lock()
  defer active_games_mutex.Unlock()
  if _, ok := active_games[name]; ok {
    return false
  }
  active_games[name] = startGame()
  return true
}

func getGame(name string) *Game {
  active_games_mutex.Lock()
  defer active_games_mutex.Unlock()
  return active_games[name]
}

func unregisterGame(name string) {
  active_games_mutex.Lock()
  defer active_games_mutex.Unlock()
  delete(active_games, name)
}

func handleHost(conn *net.TCPConn) {
  defer conn.Close()
  r := bufio.NewReader(conn)
  w := bufio.NewWriter(conn)
  line, _, err := r.ReadLine()
  if err != nil {
    // TODO: Log this error
    fmt.Printf("Failed to read game from host: %v\n", err)
    return
  }
  name := string(line)
  fmt.Printf("Game name: %s\n", line)
  if !registerGame(name) {
    w.Write([]byte(fmt.Sprintf("Unable to make game '%s', that name is already in use.\n", name)))
    return
  }
  defer unregisterGame(name)
  w.Flush()

  getGame(name)
  for {
    _, _, err = r.ReadLine()
    if err != nil {
      return
    }
  }
}

func listenForHosts() error {
  flag.Parse()
  laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", *host_port))
  if err != nil {
    return err
  }

  listener, err := net.ListenTCP("tcp", laddr)
  if err != nil {
    return err
  }

  fmt.Printf("Listening for hosts on %v\n", laddr)
  go func() {
    for {
      conn, err := listener.AcceptTCP()
      if err != nil {
        panic(err)
      }
      fmt.Printf("Host connected: %v\n", conn.RemoteAddr())
      go handleHost(conn)
    }
  }()

  return nil
}

func handleClient(conn *net.TCPConn) {
  defer conn.Close()
  defer fmt.Printf("Closed client conn.\n")

  err := conn.SetDeadline(time.Now().Add(1 * time.Second))
  if err != nil {
    fmt.Printf("Error setting deadling on client conn.\n")
    return
  }
  w := bufio.NewReader(conn)
  line, _, err := w.ReadLine()
  if err != nil {
    fmt.Printf("Error reading from player: %v\n", err)
    return
  }
  name := string(line)
  game := getGame(name)
  if game == nil {
    err_msg := fmt.Sprintf("Game '%s' does not exist.\n", name)
    conn.Write([]byte(err_msg))
    fmt.Printf(err_msg)
    return
  }

  line, _, err = w.ReadLine()
  if err != nil {
    fmt.Printf("Error reading from player: %v\n", err)
    return
  }
  if len(line) != 1 || line[0] < 0 || line[0] > 3 {
    err_msg := fmt.Sprintf("Specified an invalid seat.\n")
    conn.Write([]byte(err_msg))
    fmt.Printf(err_msg)
    return
  }
  seat := int(line[0])
  if !game.TakeSeat(seat, conn) {
    err_msg := fmt.Sprintf("Seat %d is already taken.\n", seat)
    conn.Write([]byte(err_msg))
    fmt.Printf(err_msg)
    return
  }
  conn.Write([]byte("WOOO!!!\n"))
}

func listenForClients() error {
  flag.Parse()
  laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", *client_port))
  if err != nil {
    return err
  }

  listener, err := net.ListenTCP("tcp", laddr)
  if err != nil {
    return err
  }

  fmt.Printf("Listening for clients on %v\n", laddr)
  go func() {
    for {
      conn, err := listener.AcceptTCP()
      fmt.Printf("Client connected: %v\n", conn.RemoteAddr())
      if err != nil {
        panic(err)
      }
      go handleClient(conn)
    }
  }()

  return nil
}

func main() {
  err := listenForHosts()
  if err != nil {
    fmt.Printf("Unable to listen for hosts: %v\n", err)
    return
  }

  err = listenForClients()
  if err != nil {
    fmt.Printf("Unable to listen for clients: %v\n", err)
    return
  }

  select {}
}
