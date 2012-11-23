package main

import (
  "flag"
  "fmt"
  "github.com/runningwild/barbu/barbu"
)

var player_names = []*string{
  flag.String("player0", "", "command to run for player 0"),
  flag.String("player1", "", "command to run for player 1"),
  flag.String("player2", "", "command to run for player 2"),
  flag.String("player3", "", "command to run for player 3"),
}

var seed = flag.Int64("seed", 0, "Random seed - 0 uses current time.")
var game = flag.String("game", "", "The barbu game to run")
var num_games = flag.Int("n", 1, "Number of games")
var all_perms = flag.Bool("permute", false, "Run all permutations for each deck (24 runs per deck).")

func main() {
  flag.Parse()
  players := make([]barbu.Player, 4)
  for i := range players {
    var err error
    players[i], err = barbu.MakeAiPlayer(fmt.Sprintf("%d.out", i), *player_names[i])
    if err != nil {
      fmt.Printf("Unable to make player '%s': %v\n", *player_names[i], err)
      return
    }
  }
  barbu.RunGames(players, *seed, *game, *num_games, *all_perms)
}
