package main

import (
  "bufio"
  "encoding/gob"
  "flag"
  "fmt"
  base "github.com/runningwild/barbu/net"
  "net"
  "strings"
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
    count int
  }
  host  net.Conn
  ready chan bool
}

func startGame(conn net.Conn) *Game {
  var g Game
  g.seat.take = make(chan connSeat)
  g.seat.got = make(chan bool)
  g.ready = make(chan bool)
  g.host = conn
  go g.startupRoutine()
  return &g
}
func (g *Game) TakeSeat(seat int, conn net.Conn) bool {
  g.seat.take <- connSeat{seat, conn}
  return <-g.seat.got
}

// Blocks until everyone has connected, then runs the game.
func (g *Game) Run() {
  <-g.ready
  g.activeRoutine()
}

func (g *Game) startupRoutine() {
  var ready chan bool
  for {
    select {
    case req := <-g.seat.take:
      if g.seat.taken[req.seat] != nil {
        g.seat.got <- false
      } else {
        g.seat.taken[req.seat] = req.conn
        g.seat.got <- true
        g.seat.count++
        if g.seat.count == 4 {
          ready = g.ready
        }
      }

    case ready <- true:
      return
    }
  }
}
func (g *Game) activeRoutine() {
  dead := make(chan bool, 10)
  from_host := make(chan base.RemoteData, 10)
  dec := gob.NewDecoder(g.host)
  go func() {
    for {
      var rd base.RemoteData
      err := dec.Decode(&rd)
      if err != nil {
        dead <- true
        return
      }
      from_host <- rd
    }
  }()
  from_client := make(chan base.RemoteData, 10)
  for i := 0; i < 4; i++ {
    go func(client int, conn net.Conn) {
      reader := bufio.NewReader(conn)
      for {
        // We don't trim this line of the trailing '\n' because it is expected
        // by the players.
        line, err := reader.ReadString('\n')
        if err != nil {
          // TODO: Should probably kill everything at this point
          dead <- true
          return
        }
        from_client <- base.RemoteData{client, line}
      }
    }(i, g.seat.taken[i])
  }

  enc := gob.NewEncoder(g.host)
  for {
    var err error
    select {
    case client_data := <-from_client:
      err = enc.Encode(client_data)
    case host_data := <-from_host:
      _, err = g.seat.taken[host_data.Client].Write([]byte(host_data.Line))
    case <-dead:
      return
    case <-g.seat.take:
      g.seat.got <- false
    }
    if err != nil {
      // TODO: Should kill everything here
      dead <- true
      return
    }
  }
}

var active_games_mutex sync.Mutex
var active_games map[string]*Game

func init() {
  active_games = make(map[string]*Game)
}

func registerGame(name string, conn net.Conn) bool {
  active_games_mutex.Lock()
  defer active_games_mutex.Unlock()
  if _, ok := active_games[name]; ok {
    return false
  }
  active_games[name] = startGame(conn)
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
  if game, ok := active_games[name]; ok {
    delete(active_games, name)
    for _, conn := range game.seat.taken {
      conn.Close()
    }
  }
}

func handleHost(conn *net.TCPConn) {
  defer conn.Close()
  defer fmt.Printf("Closed host conn.\n")

  r := bufio.NewReader(conn)
  w := bufio.NewWriter(conn)
  line, err := r.ReadString('\n')
  if err != nil {
    // TODO: Log this error
    fmt.Printf("Failed to read game from host: %v\n", err)
    return
  }
  name := strings.TrimSpace(line)
  fmt.Printf("Game name: %s\n", name)
  if !registerGame(name, conn) {
    w.Write([]byte(fmt.Sprintf("Unable to make game '%s', that name is already in use.\n", name)))
    return
  }
  defer unregisterGame(name)
  w.Flush()

  game := getGame(name)
  game.Run()
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
  success := false
  defer func() {
    if !success {
      conn.Close()
      fmt.Printf("Closed client conn.\n")
    }
  }()

  err := conn.SetDeadline(time.Now().Add(1 * time.Second))
  if err != nil {
    fmt.Printf("Error setting deadling on client conn.\n")
    return
  }
  w := bufio.NewReader(conn)
  line, err := w.ReadString('\n')
  if err != nil {
    fmt.Printf("Error reading from player: %v\n", err)
    return
  }
  name := strings.TrimSpace(line)
  game := getGame(name)
  active_games_mutex.Lock()
  fmt.Printf("%d active games:\n", len(active_games))
  for name := range active_games {
    fmt.Printf("%s\n", name)
  }
  active_games_mutex.Unlock()
  if game == nil {
    err_msg := fmt.Sprintf("Game '%s' does not exist.\n", name)
    conn.Write([]byte(err_msg))
    fmt.Printf(err_msg)
    return
  }

  line, err = w.ReadString('\n')
  if err != nil {
    fmt.Printf("Error reading from player: %v\n", err)
    return
  }
  var seat int = -1
  _, err = fmt.Sscanf(line, "%d", &seat)
  if err != nil || seat < 0 || seat > 3 {
    var err_msg string
    if err != nil {
      err_msg = fmt.Sprintf("Unable to parse seat: '%s': %v\n", strings.TrimSpace(line), err)
    } else {
      err_msg = fmt.Sprintf("Invalid seat: %d\n", seat)
    }
    conn.Write([]byte(err_msg))
    fmt.Printf(err_msg)
    return
  }
  if !game.TakeSeat(seat, conn) {
    err_msg := fmt.Sprintf("Seat %d is already taken.\n", seat)
    conn.Write([]byte(err_msg))
    fmt.Printf(err_msg)
    return
  }
  conn.SetDeadline(time.Time{})
  success = true
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
