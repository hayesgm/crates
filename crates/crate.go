package crates

import (
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "time"
  "log"
)

type Crate struct {
  Id bson.ObjectId "_id,omitempty"
  Machine *Machine
  Lease *Lease
}

// Create a new Crate in Mongo
func NewCrate(db *mgo.Database, machine *Machine) (crate *Crate, err error) {
  crate = &Crate{Id: bson.NewObjectId(), Machine: machine}
  err = db.C("crates").Insert(crate)
  if err != nil {
    return nil, err // That's a failure
  }

  return
}

// Find a crate by Id from Mongo
func FindCrate(db *mgo.Database, crateId string) (crate *Crate, err error) {
  bid := bson.ObjectIdHex(crateId)

  crate = &Crate{}
  err = db.C("crates").FindId(bid).One(&crate)
  if err != nil {
    return nil, err
  }

  return
}

// This will find a matching crate [currently, no criteria]
// and add a lease for that crate with a test-and-set operation
func MatchAndLease(db *mgo.Database) (crate *Crate, err error) {
  crate = &Crate{}
  
  // To be fair, this logic is wrong
  // We only want to create a lease once we find our server
  // :-/
  
  lease, err := NewLease(db)
  if err != nil {
    return nil, err
  }
  
  change := mgo.Change{
    Update: bson.M{"Lease": lease},
    ReturnNew: true,
  }

  for ; len(crate.Id) == 0; time.Sleep(5 * time.Second) {
    _, err := db.C("crates").Find(bson.M{"Lease": nil}).Apply(change, &crate)
    if err != nil {
      log.Println("Error matching:", err)
    }
  }
  
  return crate, nil
}