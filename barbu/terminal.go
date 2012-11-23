package barbu

import (
  "bufio"
  "bytes"
  "fmt"
  "github.com/runningwild/barbu/util"
  "io"
  "strings"
)

type termPlayer struct {
  stdin chanWriter

  stdout, stderr struct {
    base  *bytes.Buffer
    read  *bufio.Reader
    write *bufio.Writer
  }
}

type chanWriter struct {
  c chan string
}

func (c chanWriter) Write(data []byte) (int, error) {
  c.c <- string(data)
  <-c.c
  return len(data), nil
}
func (c chanWriter) finalizeWrite() {
  c.c <- ""
}

func makeTermPlayer() *termPlayer {
  var tp termPlayer

  tp.stdout.base = bytes.NewBuffer(nil)
  tp.stdout.read = bufio.NewReader(tp.stdout.base)
  tp.stdout.write = bufio.NewWriter(tp.stdout.base)

  tp.stderr.base = bytes.NewBuffer(nil)
  tp.stderr.read = bufio.NewReader(tp.stderr.base)

  tp.stdin.c = make(chan string)
  go tp.routine()

  return &tp
}

func (tp *termPlayer) routine() {
  me := -1
  for {
    line := <-tp.stdin.c
    tp.stdin.finalizeWrite()
    fmt.Printf("The name of the game is: %v\n", line)

    hand_line := <-tp.stdin.c
    tp.stdin.finalizeWrite()
    hand := util.Hand(strings.Fields(hand_line))
    hand.Sort()

    for i := 0; i < 13; i++ {
      fmt.Printf("\n")
      before_line := <-tp.stdin.c
      fmt.Printf("Played before you this trick: %v", before_line)
      before := strings.Fields(before_line)
      if me == -1 {
        me = len(before)
      }

      fmt.Printf("Play a card from your hand: %v\n", hand)
      var card string
      fmt.Scanf("%s", &card)
      hand.Remove(card)
      tp.stdout.base.WriteString(fmt.Sprintf("%s\n", card))

      tp.stdin.finalizeWrite()
      after_line := <-tp.stdin.c
      tp.stdin.finalizeWrite()
      fmt.Printf("Played after you this trick: %v", after_line)
      after := strings.Fields(after_line)

      var all []string
      for _, c := range before {
        all = append(all, c)
      }
      all = append(all, card)
      for _, c := range after {
        all = append(all, c)
      }
      lead := all[0][1]
      max := -1
      index := -1
      for i := range all {
        if all[i][1] != lead {
          continue
        }
        if rank_map[all[i][0]] > max {
          max = rank_map[all[i][0]]
          index = i
        }
      }
      player := (me - len(before) + index) % 4
      fmt.Printf("Player %d won the trick with %s\n", player, all[index])
    }
  }
}

func (tp *termPlayer) GiveDeck(deck Deck) {
}
func (tp *termPlayer) Stdin() io.Writer {
  return tp.stdin
}
func (tp *termPlayer) Stdout() *bufio.Reader {
  return tp.stdout.read
}
func (tp *termPlayer) Stderr() *bufio.Reader {
  return tp.stderr.read
}
func (tp *termPlayer) Close() {
}
