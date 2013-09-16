package main

/*
  A Lease server which manages Lease allocations, backed by Mongo.
*/

import (
  "github.com/hayesgm/crates/server"
  "labix.org/v2/mgo"
  "log"
)

// Main runs a lease server
// First, we'll connect to a MongoDB
func main() {
  log.Println("Welcome to Crates\n")

  session, err := mgo.Dial("localhost")
  if err != nil {
    panic(err)
  }
  defer session.Close()

  log.Println("Connected to localhost")

  server.RunLeaseServer(session.DB("crates"))
}