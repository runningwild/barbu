package main

import (
  "strings"
  "fmt"
  "flag"
  "io"
  "os"
  "bufio"
  "time"
  "math/rand"
  "os/exec"
)

var rank_map map[byte]int
func init() {
  rand.Seed(time.Now().UnixNano())
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

var player_names = []*string{
  flag.String("player1", "", "command to run for player 1"),
  flag.String("player2", "", "command to run for player 2"),
  flag.String("player3", "", "command to run for player 3"),
  flag.String("player4", "", "command to run for player 4"),
}

var suits = []byte{'s','h','c','d'}
var ranks = []byte{'2','3','4','5','6','7','8','9','t','j','q','k','a'}

type card string
func less(a, b card) bool {
  return rank_map[a[0]] < rank_map[b[0]]
}
type Deck []card

func (d Deck) String() string {
  var s string
  for i := range d {
    s = s + string(d[i])
    if i < len(d) - 1 {
      s = s + " "
    }
  }
  return s
}
func (d Deck) Deal() [4]Deck {
  return [4]Deck{d[0:13], d[13:26], d[26:39], d[39:52]}
}
func makeDeck() Deck {
  var d Deck
  for _, suit := range suits {
    for _, rank := range ranks {
      d = append(d, card([]byte{rank, suit}))
    }
  }
  for i := range d {
    k := rand.Intn(len(d))
    d[i], d[k] = d[k], d[i]
  }
  return d
}

type Player interface {
  Stdin()  io.Writer
  Stdout() *bufio.Reader
  Stderr() *bufio.Reader
  Close()
}

type termPlayer struct {
  // receive hand
  // receive start of trick and rest of trick
  // send play
}
func MakeTermPlayer() Player {
  var tp termPlayer
  go tp.routine()
  return nil
}
func (tp *termPlayer) routine() {
  // First get hand

  for {

  }
}
func (tp *termPlayer) Stderr() *bufio.Reader {
  return bufio.NewReader(os.Stderr)
}

type aiPlayer struct {
  cmd *exec.Cmd
  stdin  io.Writer
  stdout *bufio.Reader
  stderr *bufio.Reader
}
func (a *aiPlayer) Stdin() io.Writer {
  return a.stdin
}
func (a *aiPlayer) Stdout() *bufio.Reader {
  return a.stdout
}
func (a *aiPlayer) Stderr() *bufio.Reader {
  return a.stderr
}
func (a *aiPlayer) Close() {
  a.cmd.Wait()
}
func MakeAiPlayer(log_filename, name string) (Player, error) {
  var p aiPlayer
  params := strings.Fields(name)
  p.cmd = exec.Command(params[0], params[1:]...)
  log, err := os.Create(log_filename)
  if err != nil {
    return nil, err
  }
  in, err := p.cmd.StdinPipe()
  if err != nil {
    return nil, err
  }
  out, err := p.cmd.StdoutPipe()
  if err != nil {
    return nil, err
  }
  stderr, err := p.cmd.StderrPipe()
  if err != nil {
    return nil, err
  }
  err = p.cmd.Start()
  if err != nil {
    return nil, err
  }
  p.stdin = io.MultiWriter(in, log)
  p.stdout = bufio.NewReader(out)
  p.stderr = bufio.NewReader(stderr)
  return &p, nil
}

func main() {
  flag.Parse()
  for i := range player_names {
    if player_names[i] == nil {
      fmt.Fprintf(os.Stderr, "Must specify all 4 players\n")
      return
    }
  }

  var total [4]int
  N := 100
  for i := 0; i < N; i++ {
    var players [4]Player
    var err error
    for i := range player_names {
      players[i], err = MakeAiPlayer(fmt.Sprintf("%d.out", i), *player_names[i])
      if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
      }
    }

    r := MakeRavage(players, makeDeck())
    r.Deal()
    for i := 0; i < 13; i++ {
      r.Round()
    }
    scores := r.Score()
    for i := range scores {
      total[i] += scores[i]
      players[i].Close()
    }
  }
  fmt.Printf("Averages:\n")
  for i := range total {
    fmt.Printf("Player %d: %.2f\n", i, float64(total[i]) / float64(N))
  }
}









