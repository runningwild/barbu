package main

import (
  "flag"
  "fmt"
  "net"
)

var addr = flag.String("addr", "", "Address of the server.")
var host_port = flag.Int("host_port", 9901, "Port to connect to server for hosting.")
var player_port = flag.Int("player_port", 9902, "Port to connect to server as a player.")
var seat = flag.Int("seat", -1, "Player's seat ([0-3]) when connecting as a player.")
var host = flag.Bool("host", false, "True to host, false to join.")
var name = flag.String("name", "", "Name of the game to join.")

func connectAsHost() {
  raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", *addr, *host_port))
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

  _, err = conn.Write([]byte(fmt.Sprintf("%s\n", *name)))
  if err != nil {
    fmt.Printf("Failed to host game: %v\n", err)
    return
  }
  buf := make([]byte, 1024)
  n, err := conn.Read(buf)
  if err != nil {
    fmt.Printf("Failed to host game: %v\n", err)
    return
  }
  fmt.Printf("%s\n", buf[0:n])
}

func connectAsPlayer() {
  raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", *addr, *player_port))
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

  _, err = conn.Write([]byte(*name + "\n"))
  if err != nil {
    fmt.Printf("Unable to join game: %v\n", err)
    return
  }

  _, err = conn.Write([]byte{byte(*seat), '\n'})
  if err != nil {
    fmt.Printf("Failed to specify seat to server: %v\n", err)
    return
  }
  buf := make([]byte, 1024)
  n, err := conn.Read(buf)
  if err != nil {
    fmt.Printf("Error reading from server: %v\n", err)
    return
  }
  fmt.Printf("%s\n", buf[0:n])
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
  if *host {
    connectAsHost()
  } else {
    connectAsPlayer()
  }
}
