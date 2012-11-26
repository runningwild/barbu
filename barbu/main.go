package barbu

import (
  "bufio"
  "fmt"
  "github.com/runningwild/cmwc"
  "io"
  "os"
  "os/exec"
  "sort"
  "strings"
)

type BarbuGame interface {
  // Runs the doubling for this game.  Most games will probably have the same
  // doubling scheme, but this allows for special cases, like trumps.
  // Returns a mapping from player to a slice of each player that player
  // doubled.
  // Note: This matrix should be symmetric, a[i][j] == a[j][i]
  Double() [4][4]int

  // Runs the game to completion.  Will only be called once.
  Run()

  // Returns the PRE-doubling scores for each player this game.
  Scores() [4]int
}

type BarbuGameMaker func(players []Player, hands [][]string) BarbuGame

var allBarbuGames map[string]BarbuGameMaker
var allBarbuGameNames []string

// Registers a mapping from name to a factory function for barbu games.
func RegisterBarbuGame(name string, maker BarbuGameMaker) {
  if allBarbuGames == nil {
    allBarbuGames = make(map[string]BarbuGameMaker)
  }
  if _, ok := allBarbuGames[name]; ok {
    panic(fmt.Sprintf("Tried to register a BarbuGame '%s' twice.", name))
  }
  allBarbuGames[name] = maker
  allBarbuGameNames = append(allBarbuGameNames, name)
  sort.Strings(allBarbuGameNames)
}
func AllBarbuGameNames() []string {
  return allBarbuGameNames
}
func GetBarbuGame(name string, players []Player, hands [][]string) BarbuGame {
  if allBarbuGames == nil {
    return nil
  }
  return allBarbuGames[name](players, hands)
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

func less(a, b string) bool {
  return rank_map[a[0]] < rank_map[b[0]]
}

type Deck []string

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
func (d Deck) Deal() [][]string {
  return [][]string{[]string(d[0:13]), []string(d[13:26]), []string(d[26:39]), []string(d[39:52])}
}
func makeDeck(rng *cmwc.Cmwc) Deck {
  var d Deck
  for _, suit := range suits {
    for _, rank := range ranks {
      d = append(d, string([]byte{rank, suit}))
    }
  }
  for i := range d {
    k := int(rng.Int63()%int64(len(d)-i)) + i
    d[i], d[k] = d[k], d[i]
  }
  return d
}

type Player interface {
  // If two players return the same Id() they should be equivalent.  i.e. they
  // will always make the same moves as each other given the same exact
  // situation.
  Id() string
  Stdin() io.Writer
  Stdout() *bufio.Reader
  Stderr() *bufio.Reader
  Close()
}

type aiPlayer struct {
  id     string
  cmd    *exec.Cmd
  stdin  io.Writer
  stdout *bufio.Reader
  stderr *bufio.Reader
}

func (a *aiPlayer) Id() string {
  return a.id
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
  p.id = name
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

func equivalentPermutation(players []Player, pi, pj int) bool {
  for i := 0; i < 4; i++ {
    if players[perms[pi][i]].Id() != players[perms[pj][i]].Id() {
      return false
    }
  }
  return true
}

type runGamesErrors struct {
  msg string
}

func (rge runGamesErrors) Error() string {
  return rge.msg
}

func RunGames(players []Player, seed int64, game string, num_games int, all_perms bool) (err error) {
  defer func() {
    if r := recover(); r != nil {
      err = runGamesErrors{fmt.Sprintf("%v", r)}
    }
  }()

  rng := cmwc.MakeGoodCmwc()
  if seed != 0 {
    rng.Seed(int64(seed))
  } else {
    rng.SeedWithDevRand()
  }

  valid_game := false
  for _, valid_name := range AllBarbuGameNames() {
    if game == valid_name {
      valid_game = true
    }
  }
  if !valid_game {
    fmt.Printf("'%s' is not a valid game.  Valid games are %v.\n", AllBarbuGameNames())
  }

  var total [4]int
  N := num_games
  if !all_perms {
    perms = [][]int{{0, 1, 2, 3}}
  }
  total_games := N * len(perms)
  completed := 0

  for i := 0; i < N; i++ {
    // Map from N to the scores for the Nth permutation
    perm_score_map := make(map[int][]int)
    deck := makeDeck(rng)
    for perm_num, perm := range perms {
      equivalent := false
      equivalent_perm_num := -1
      for pi := 0; pi < perm_num; pi++ {
        if equivalentPermutation(players, pi, perm_num) {
          equivalent = true
          equivalent_perm_num = pi
          break
        }
      }
      scores := make([]int, 4)
      permed_players := make([]Player, 4)
      if equivalent {
        scores = perm_score_map[equivalent_perm_num]
      } else {
        for i := range permed_players {
          permed_players[i] = players[perm[i]]
        }

        // Tell the permed_players what position they are around the table, and what
        // their hand is.
        hands := deck.Copy().Deal()
        for i := range permed_players {
          permed_players[i].Stdin().Write([]byte(fmt.Sprintf("PLAYER %d\n", i)))
          for j := range hands[i] {
            var line string
            // Prevent having a trailing space
            if j == 0 {
              line = hands[i][j]
            } else {
              line = " " + hands[i][j]
            }
            permed_players[i].Stdin().Write([]byte(line))
          }
          permed_players[i].Stdin().Write([]byte("\n"))
        }

        // TODO: This is where the dealer should choose the game.
        for i := range permed_players {
          permed_players[i].Stdin().Write([]byte(strings.ToUpper(game) + "\n"))
        }
        the_game := GetBarbuGame(game, permed_players, hands)
        doubles := the_game.Double()
        the_game.Run()
        raw_scores := the_game.Scores()
        for i := range scores {
          for j := range scores {
            if i >= j {
              continue
            }
            if doubles[i][j] == 0 || raw_scores[i]-raw_scores[j] == 0 {
              continue
            }
            diff := raw_scores[i] - raw_scores[j]
            if diff < 0 {
              diff = -diff
            }
            diff *= doubles[i][j]
            if raw_scores[i] > raw_scores[j] {
              scores[i] += diff
              scores[j] -= diff
            } else {
              scores[j] += diff
              scores[i] -= diff
            }
          }
        }
        for i := range scores {
          scores[i] += raw_scores[i]
        }
      }
      perm_score_map[perm_num] = scores
      for i := range scores {
        if !equivalent {
          line := fmt.Sprintf("END %d %d %d %d\n", scores[0], scores[1], scores[2], scores[3])
          permed_players[i].Stdin().Write([]byte(line))
        }
        total[perm[i]] += scores[i]
      }
      completed++
      fmt.Printf("\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b")
      fmt.Printf("\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b")
      fmt.Printf("Finished %d/%d games.", completed, total_games)
    }
  }
  fmt.Printf("\nAverages:\n")
  for i := range total {
    fmt.Printf("Player %d: %.2f\n", i, float64(total[i])/float64(N*len(perms)))
  }
  return
}
