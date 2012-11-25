package main

import (
  "bufio"
  "flag"
  "fmt"
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

// params := strings.Fields(name)
// p.cmd = exec.Command(params[0], params[1:]...)
// log, err := os.Create(log_filename)
// if err != nil {
//   return nil, err
// }
// in, err := p.cmd.StdinPipe()
// if err != nil {
//   return nil, err
// }
// out, err := p.cmd.StdoutPipe()
// if err != nil {
//   return nil, err
// }
// stderr, err := p.cmd.StderrPipe()
// if err != nil {
//   return nil, err
// }
// err = p.cmd.Start()

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
    return
  }

  buf := bufio.NewReader(conn)
  go func() {
    buf := bufio.NewReader(os.Stdin)
    for {
      line, err := buf.ReadString('\n')
      if err != nil {
        panic(err)
      }
      conn.Write([]byte(line))
    }
  }()
  for {
    line, err := buf.ReadString('\n')
    if err != nil {
      fmt.Printf("ERROR: %v\n", err)
      return
    }
    fmt.Printf("%s\n", line)
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
