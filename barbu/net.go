package barbu

import (
  "bufio"
  "bytes"
  "fmt"
  "github.com/runningwild/barbu/util"
  "io"
  "net"
  "strings"
  "time"
)

type netPlayer struct {
  stdin chanWriter

  conn *net.TCPConn

  stdout struct {
    base *bytes.Buffer
    read *bufio.Reader
  }

  stderr *bufio.Reader
}

// Net player will listen on port 5200+player_num until it receives a connection
// or times out.
func makeNetPlayer(player_num int) (*netPlayer, error) {
  addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", 5200+player_num))
  if err != nil {
    return nil, err
  }

  listener, err := net.ListenTCP("tcp", addr)
  if err != nil {
    return nil, err
  }

  err = listener.SetDeadline(time.Now().Add(time.Second * 15))
  if err != nil {
    return nil, err
  }

  var tp netPlayer
  tp.conn, err = listener.AcceptTCP()
  if err != nil {
    return nil, err
  }

  tp.stdout.base = bytes.NewBuffer(nil)
  tp.stdout.read = bufio.NewReader(tp.stdout.base)

  tp.stdin.c = make(chan string)
  go tp.routine()

  tp.stderr = bufio.NewReader(nil)

  return &tp, nil
}

func (tp *netPlayer) routine() {
  me := -1
  for {
    line := <-tp.stdin.c
    tp.stdin.finalizeWrite()
    _, err := tp.conn.Write([]byte(line))
    if err != nil {
      panic(err)
    }

    hand_line := <-tp.stdin.c
    tp.stdin.finalizeWrite()
    hand := util.Hand(strings.Fields(hand_line))
    hand.Sort()
    _, err = tp.conn.Write([]byte(hand_line))
    if err != nil {
      panic(err)
    }

    for i := 0; i < 13; i++ {
      fmt.Printf("\n")
      before_line := <-tp.stdin.c
      _, err = tp.conn.Write([]byte(before_line))
      if err != nil {
        panic(err)
      }
      before := strings.Fields(before_line)
      if me == -1 {
        me = len(before)
      }

      card_bytes := make([]byte, 2)
      _, err = tp.conn.Read(card_bytes)
      if err != nil {
        panic(err)
      }
      card := string(card_bytes)
      hand.Remove(card)
      tp.stdout.base.WriteString(card)

      tp.stdin.finalizeWrite()
      after_line := <-tp.stdin.c
      tp.stdin.finalizeWrite()

      _, err = tp.conn.Write([]byte(after_line))
      if err != nil {
        panic(err)
      }
    }
  }
}

func (tp *netPlayer) GiveDeck(deck Deck) {
}
func (tp *netPlayer) Stdin() io.Writer {
  return tp.stdin
}
func (tp *netPlayer) Stdout() *bufio.Reader {
  return tp.stdout.read
}
func (tp *netPlayer) Stderr() *bufio.Reader {
  return tp.stderr
}
func (tp *netPlayer) Close() {
}
