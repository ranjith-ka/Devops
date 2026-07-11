package main

import (
  "log"
  "net/http"
  "os"
  "time"

  "github.com/ranjith-ka/Devops/services/inference/internal/server"
)

func main() {
  addr := os.Getenv("PORT")
  if addr == "" {
    addr = ":8080"
  } else if addr[0] != ':' {
    addr = ":" + addr
  }

  httpServer := &http.Server{
    Addr:              addr,
    Handler:           server.NewRouter(),
    ReadHeaderTimeout: 5 * time.Second,
  }

  log.Printf("inference service listening on %s", addr)
  if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
    log.Fatal(err)
  }
}
