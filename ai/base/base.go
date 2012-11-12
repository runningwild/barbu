package base

import (
  "flag"
  "fmt"
  "sort"
)

var mode = flag.String("mode", "random", "What ai to use")
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

type Hand []string

func (h Hand) Len() int {
  return len(h)
}
func (h Hand) Less(i, j int) bool {
  if h[i][1] != h[j][1] {
    return h[i][1] < h[j][1]
  }
  return rank_map[h[i][0]] < rank_map[h[j][0]]
}
func (h Hand) Swap(i, j int) {
  h[i], h[j] = h[j], h[i]
}
func (h *Hand) Remove(card string) {
  for i := range *h {
    if (*h)[i] == card {
      (*h)[i] = (*h)[len((*h))-1]
      (*h) = (*h)[0 : len((*h))-1]
      sort.Sort((*h))
      return
    }
  }
  panic(fmt.Sprintf("Didn't find '%s' in the hand '%v'.\n", card, h))
}
func (h Hand) Sort() {
  sort.Sort(h)
}
