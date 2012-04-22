package main

import (
  "fmt"
  "strings"
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

var player_names = flag.String("players", "", "comma delimited list of player binaries")

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

type Player struct {
  cmd *exec.Cmd
  Stdin  io.Writer
  Stdout *bufio.Reader
  Stderr *bufio.Reader
}
func MakePlayer(log_filename, name string) (*Player, error) {
  var p Player
  p.cmd = exec.Command(name)
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
  p.Stdin = io.MultiWriter(in, log)
  p.Stdout = bufio.NewReader(out)
  p.Stderr = bufio.NewReader(stderr)
  return &p, nil
}

func main() {
  flag.Parse()
  player_names_slice := strings.Split(*player_names, ",")
  if len(player_names_slice) != 4 {
    fmt.Printf("Must specify exactly 4 players\n")
    return
  }

  var players [4]*Player
  var err error
  for i := range player_names_slice {
    players[i], err = MakePlayer(fmt.Sprintf("%d.out", i), player_names_slice[i])
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
  fmt.Printf("Scores: %v\n", r.Score())
}









