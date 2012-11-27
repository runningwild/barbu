package main

import (
  "bufio"
  "flag"
  "fmt"
  "github.com/runningwild/barbu/util"
  "io"
  "net"
  "os"
  "os/exec"
  "strings"
)

var addr = flag.String("addr", "", "Address of the server.")
var port = flag.Int("port", 9902, "Port to connect to server on.")
var seat = flag.Int("seat", -1, "Player's seat ([0-3]) when connecting as a player.")
var name = flag.String("name", "", "Name of the game to join.")
var cmd = flag.String("cmd", "", "Command to run ai, leave blank for interactive mode.")

func runCmd(cmd string, conn net.Conn) {
  params := strings.Fields(cmd)
  c := exec.Command(params[0], params[1:]...)
  stdin, err := c.StdinPipe()
  if err != nil {
    panic(err)
  }
  stdout, err := c.StdoutPipe()
  if err != nil {
    panic(err)
  }
  stderr, err := c.StderrPipe()
  if err != nil {
    panic(err)
  }
  err = c.Start()
  if err != nil {
    fmt.Printf("Error running cmd '%s': %v\n", cmd, err)
    return
  }
  dev_null, err := os.Open(os.DevNull)
  if err != nil {
    panic(err)
  }
  go io.Copy(dev_null, stderr)
  go io.Copy(conn, stdout)
  io.Copy(stdin, conn)
}

func interactiveMode(conn net.Conn) {
  buf := bufio.NewReader(conn)
  stdin := bufio.NewReader(os.Stdin)
  var hand util.Hand
  for {
    // Read in what seat you are
    line, err := buf.ReadString('\n')
    if err != nil {
      if err == io.EOF {
        return
      }
      panic(err)
    }
    line = strings.TrimSpace(line)
    fmt.Printf("You are: %s\n", line)

    // Read in your hand
    line, err = buf.ReadString('\n')
    hand = util.Hand(strings.Fields(line))
    hand.Sort()
    fmt.Printf("Your hand is: %s\n", hand)

    // Read in the game
    line, err = buf.ReadString('\n')
    line = strings.TrimSpace(line)
    fmt.Printf("The game is: %s\n", line)
    // TODO: Need to support non-trick-taking games

    // Read in 'DOUBLING', then all of the doubling info
    line, err = buf.ReadString('\n')
    line = strings.TrimSpace(line)
    fmt.Printf("READ DOUBLING: %s\n", line)
    for i := 0; i < 4; i++ {
      fmt.Printf("Reading ... ")
      line, err = buf.ReadString('\n')
      fmt.Printf(" Read: '%s'\n", line)
      line = strings.TrimSpace(line)
      if line == "DOUBLE" {
        fmt.Printf("Time to double!\n")
        line, err = stdin.ReadString('\n')
        _, err = conn.Write([]byte(line))
      } else {
        fmt.Printf("DOUBLING INFO: %s\n", line)
      }
    }

    for {
      line, err = buf.ReadString('\n')
      line = strings.TrimSpace(line)
      if line != "TRICK" {
        break
      }

      fmt.Printf("\nYour hand: %v\n", hand)
      for i := 0; i < 4; i++ {
        line, err = buf.ReadString('\n')
        line = strings.TrimSpace(line)
        if line == "PLAY" {
          line, err = stdin.ReadString('\n')
          _, err = conn.Write([]byte(line))
          hand.Remove(strings.TrimSpace(line))
        } else {
          fmt.Printf("TRICK INFO: %s\n", line)
        }
      }
    }
    fmt.Printf("END: %s\n", line)
  }
}

func connectAsPlayer(cmd, addr string, port int, name string, seat int) {
  raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", addr, port))
  if err != nil {
    fmt.Printf("Unable to resolve server address: %v\n", err)
    return
  }

  conn, err := net.DialTCP("tcp", nil, raddr)
  if err != nil {
    fmt.Printf("Unable to connect to server: %v\n", err)
    return
  }
  defer conn.Close()

  _, err = conn.Write([]byte(name + "\n"))
  if err != nil {
    fmt.Printf("Unable to join game: %v\n", err)
    return
  }

  _, err = conn.Write([]byte(fmt.Sprintf("%d\n", seat)))
  if err != nil {
    fmt.Printf("Failed to specify seat to server: %v\n", err)
    return
  }

  if cmd != "" {
    runCmd(cmd, conn)
  } else {
    interactiveMode(conn)
  }

}
func main() {
  flag.Parse()
  if *addr == "" {
    fmt.Printf("Must specify an address with --addr.\n")
    return
  }
  if *name == "" {
    fmt.Printf("Must specify a game name with --name.\n")
    return
  }
  connectAsPlayer(*cmd, *addr, *port, *name, *seat)
}
