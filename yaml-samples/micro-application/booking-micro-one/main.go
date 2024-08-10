// main.go

package main

import (
 "fmt"
 "log"
 "net/http"
)

func main() {
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, booking!")
 })

 fmt.Printf("Starting server at port 4000\n")
 if err := http.ListenAndServe(":4000", nil); err != nil {
  log.Fatal(err)
 }
}
