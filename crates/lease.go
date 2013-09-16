package crates

import (
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "time"
)

type Lease struct {
  Id bson.ObjectId "_id,omitempty"
  Expiration time.Time
}

func NewLease(db *mgo.Database) (lease *Lease, err error) {
  lease = &Lease{Id: bson.NewObjectId(), Expiration: time.Now()}
  err = db.C("leases").Insert(lease)
  if err != nil {
    return nil, err // That's a failure
  }

  return
}