package main

import (
  "bufio"
  "flag"
  "fmt"
  "github.com/runningwild/cmwc"
  "io"
  "os"
  "os/exec"
  "runtime/pprof"
  "sort"
  "strings"
)

// profiling info
var cpu_profile = flag.String("cpuprof", "", "file to write cpu profile to")

var player_names = []*string{
  flag.String("player1", "", "command to run for player 1"),
  flag.String("player2", "", "command to run for player 2"),
  flag.String("player3", "", "command to run for player 3"),
  flag.String("player4", "", "command to run for player 4"),
}

var seed = flag.Int64("seed", 0, "Random seed - 0 uses current time.")
var game = flag.String("game", "", "The barbu game to run")
var num_games = flag.Int("n", 1, "Number of games")
var all_perms = flag.Bool("permute", false, "Run all permutations for each deck (24 runs per deck).")

var rng *cmwc.Cmwc

type BarbuGame interface {
  Deal()
  Round() bool   // returns true iff game is over
  Score() [4]int // only call this after the game is over

  // Given the string that a player would normally be given before choosing
  // what to play, returns an array containing all valid plays
  GetValidPlays(hand []string, lead string) []string
}

type StandardTrickTakingGame struct{}

func (StandardTrickTakingGame) GetValidPlays(hand []string, lead string) []string {
  if len(lead) == 0 {
    return hand
  }
  suit := strings.Split(lead, " ")[0][1]
  var valid []string
  for _, card := range hand {
    if card[1] == suit {
      valid = append(valid, card)
    }
  }
  if len(valid) == 0 {
    return hand
  }
  return valid
}

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

var suits = []byte{'s', 'h', 'c', 'd'}
var ranks = []byte{'2', '3', '4', '5', '6', '7', '8', '9', 't', 'j', 'q', 'k', 'a'}

type card string

func less(a, b card) bool {
  return rank_map[a[0]] < rank_map[b[0]]
}

type Deck []card

func (d Deck) Len() int {
  return len(d)
}
func (d Deck) Less(i, j int) bool {
  return d[i] < d[j]
}
func (d Deck) Swap(i, j int) {
  d[i], d[j] = d[j], d[i]
}

func (d Deck) Copy() Deck {
  d2 := make(Deck, len(d))
  copy(d2, d)
  return d2
}

func (d Deck) String() string {
  var s string
  for i := range d {
    s = s + string(d[i])
    if i < len(d)-1 {
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
    k := int(rng.Int63()%int64(len(d)-i)) + i
    d[i], d[k] = d[k], d[i]
  }
  return d
}

type Player interface {
  Stdin() io.Writer
  Stdout() *bufio.Reader
  Stderr() *bufio.Reader
  Close()
}

type aiPlayer struct {
  cmd    *exec.Cmd
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

  go func() {
    for {
      line, _, err := p.Stderr().ReadLine()
      if err != nil {
        return
      }
      fmt.Printf("Error(%s): %s\n", name, line)
    }
  }()

  return &p, nil
}

var perms = [][]int{
  {0, 1, 2, 3},
  {0, 1, 3, 2},
  {0, 2, 1, 3},
  {0, 2, 3, 1},
  {0, 3, 1, 2},
  {0, 3, 2, 1},
  {1, 0, 2, 3},
  {1, 0, 3, 2},
  {1, 2, 0, 3},
  {1, 2, 3, 0},
  {1, 3, 0, 2},
  {1, 3, 2, 0},
  {2, 0, 1, 3},
  {2, 0, 3, 1},
  {2, 1, 0, 3},
  {2, 1, 3, 0},
  {2, 3, 0, 1},
  {2, 3, 1, 0},
  {3, 0, 1, 2},
  {3, 0, 2, 1},
  {3, 1, 0, 2},
  {3, 1, 2, 0},
  {3, 2, 0, 1},
  {3, 2, 1, 0},
}

func main() {
  flag.Parse()
  rng = cmwc.MakeGoodCmwc()
  if *seed != 0 {
    rng.Seed(int64(*seed))
  } else {
    rng.SeedWithDevRand()
  }

  for i := range player_names {
    if player_names[i] == nil {
      fmt.Fprintf(os.Stderr, "Must specify all 4 players with --player1 - --player4\n")
      return
    }
  }

  if *game == "" {
    fmt.Fprintf(os.Stderr, "Must specify a game with --game.\n")
    return
  }

  if *cpu_profile != "" {
    profile_output, err := os.Create(*cpu_profile)
    if err != nil {
      fmt.Printf("Unable to start CPU profile: %v\n", err)
    } else {
      err = pprof.StartCPUProfile(profile_output)
      if err != nil {
        fmt.Printf("Unable to start CPU profile: %v\n", err)
        profile_output.Close()
      } else {
        defer profile_output.Close()
        defer pprof.StopCPUProfile()
      }
    }
  }

  makers := map[string]func([4]Player, Deck) BarbuGame{
    "ravage":  MakeRavage,
    "lasttwo": MakeLastTwo,
  }
  game_maker, ok := makers[*game]
  if !ok {
    var names []string
    for name := range makers {
      names = append(names, name)
    }
    sort.Strings(names)
    fmt.Printf("'%s' is not a valid game.  Valid games are %v.\n", names)
  }

  var total [4]int
  N := *num_games
  if !*all_perms {
    perms = [][]int{{0, 1, 2, 3}}
  }
  var orig_players [4]Player
  for i := range player_names {
    var err error
    if *player_names[i] == "terminal" {
      orig_players[i] = makeTermPlayer()
    } else if *player_names[i] == "net" {
      orig_players[i], err = makeNetPlayer(i)
      if err != nil {
        panic(err)
      }
    } else {
      orig_players[i], err = MakeAiPlayer(fmt.Sprintf("%d.out", i), *player_names[i])
      if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
      }
    }
  }
  for i := 0; i < N; i++ {
    deck := makeDeck()
    for _, perm := range perms {
      var players [4]Player
      for i := range players {
        players[i] = orig_players[perm[i]]
      }
      the_game := game_maker(players, deck.Copy())
      the_game.Deal()
      for !the_game.Round() {
      }
      scores := the_game.Score()
      for i := range scores {
        total[perm[i]] += scores[i]
      }
    }
  }
  fmt.Printf("Averages:\n")
  for i := range total {
    fmt.Printf("Player %d: %.2f\n", i, float64(total[i])/float64(N*len(perms)))
  }
}
