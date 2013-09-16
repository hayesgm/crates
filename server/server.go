package server

import (
  "net/http"
  "github.com/hayesgm/crates/crates"
  "encoding/json"
  "log"
  "labix.org/v2/mgo"
  "time"
)

var registeredCrates = make(map[int]*crates.Crate)

type RegisterResponse struct {
  Crate *crates.Crate
}

type AcquireResponse struct {
  Lease *crates.Lease
}

type LeaseResponse struct {
  Crate *crates.Crate
}

// Register will set-up a new crate for a machine
// This is one of the few functions which will return right away
//
// Signature:
// <-
// info:: String describing the machine
// -> (json)
// { CrateId: '<new crate id>'}
func HandleRegister(db *mgo.Database, w http.ResponseWriter, req *http.Request) {
  machine := &crates.Machine{Info: req.FormValue("info")} // TODO: Verify

  // TODO: Verify parameters being passed in correctly

  crate, err := crates.NewCrate(db, machine)
  if err != nil {
    panic(err)
  }

  resp := RegisterResponse{Crate: crate}
  respBytes, err := json.Marshal(&resp)
  if err != nil {
    panic(err)
  }
  
  if _, err := w.Write(respBytes); err != nil {
    panic(err)
  }
}

// Acquires a lease for this server
// This is going to a long running HTTP request, which will return when a lease has been acquired
// If connection times out, client should reconnect
//
// Signature
// <-
// CrateId:: Crate id from registration
// -> (json)
// Lease:: Lease acquired
func HandleAcquire(db *mgo.Database, w http.ResponseWriter, req *http.Request) {
  var crate *crates.Crate

  crateId := req.FormValue("CrateId")
  
  for ;; time.Sleep(5*time.Second) {
    // TODO: Verify

    crate, err := crates.FindCrate(db, crateId)
    if err != nil {
      panic(err)
    }
    log.Printf("Acquiring Crate: %#v\n",crate)
    if crate.Lease != nil {
      break
    }
  }
  
  // log.Printf("Found crate: %#v, waiting for lease...", crate)


  // Move this crate into looking for a lease
  // crate.Lease = <- newLeaseCh

  resp := AcquireResponse{Lease: crate.Lease}
  respBytes, err := json.Marshal(&resp)
  if err != nil {
    panic(err)
  }
  
  if _, err := w.Write(respBytes); err != nil {
    panic(err)
  }
}

// This is the other half of crate acquisition
// Handle lease will look for a valid server and acquire a lease for it
//
// Signature
// <-
// CrateId:: Acquire a specific crate
// ->
// Crate:: Crate acquired
func HandleLease(db *mgo.Database, w http.ResponseWriter, req *http.Request) {
  
  crate, err := crates.MatchAndLease(db) // TODO: Verify
  if err != nil {
    panic(err)
  }

  log.Printf("Lease acquired for %#v\n", crate)

  resp := LeaseResponse{Crate: crate}
  respBytes, err := json.Marshal(&resp)
  if err != nil {
    panic(err)
  }
  
  if _, err := w.Write(respBytes); err != nil {
    panic(err)
  }
}

// This is run once a lease has been acquired
// This function will hold until the lease is no longer held
// If this returns forfeit, the lease has been forfeited.
// 
// Signature
// <-
// LeaseId:: Id of acquired lease
// ->
// Forfeit:: If true, this lease is forfeitted
func HandleHold(w http.ResponseWriter, req *http.Request) {

}

// This is for voluntarily dropping a lease
// When this occurs, lease should be immediately forfeitted
// 
// Signature
// <-
// LeaseId:: Id of lease to release
// ->
// Forfeit:: If true, the lease is forfeitted
func HandleForfeit(w http.ResponseWriter, req *http.Request) {

}

// Unregister a create
// Run this to unregister a crate
// Any lease associated with that crate will be forfeitted
//
// Signature
// <-
// CrateId:: Id of create to unregister
// ->
// Dropped:: True if crate has been unregistered
func HandleUnregister(w http.ResponseWriter, req *http.Request) {

}

func RunLeaseServer(db *mgo.Database) {

  handleWithDatabase := func(f func(db *mgo.Database, w http.ResponseWriter, req *http.Request)) (func(w http.ResponseWriter, req *http.Request)) {
    return func(w http.ResponseWriter, req *http.Request) {
      f(db, w, req)
    }
  }

  http.Handle("/register", http.HandlerFunc(handleWithDatabase(HandleRegister)))
  http.Handle("/acquire", http.HandlerFunc(handleWithDatabase(HandleAcquire)))
  http.Handle("/lease", http.HandlerFunc(handleWithDatabase(HandleLease)))
  http.Handle("/hold", http.HandlerFunc(HandleHold))
  http.Handle("/forfeit", http.HandlerFunc(HandleForfeit))
  http.Handle("/unregister", http.HandlerFunc(HandleUnregister))

  log.Println("Running Crates Server...")

  err := http.ListenAndServe("localhost:2353", nil)
  if err != nil {
    log.Fatal("ListenAndServer:",err)
  }
}