package main

import (
  "bufio"
  "fmt"
  "os"
  "strings"
)
func startConsole() {
	reader := bufio.NewReader(os.Stdin)
  fmt.Println("Simple Shell")
  fmt.Println("---------------------")

  for {
    fmt.Print("-> ")
    text, _ := reader.ReadString('\n')
    // convert CRLF to LF
    text = strings.Replace(text, "\n", "", -1)

    if strings.Compare("q", text) == 0 {
	  fmt.Println("hello, Yourself")
	  return
    }

  }
}