package client

import (
  "net/http"
  "github.com/hayesgm/crates/server"
  "fmt"
  "encoding/json"
  "io/ioutil"
)

func Register(endpoint string) (resp server.RegisterResponse, err error) {
  r, err := http.Get(fmt.Sprintf("http://%s/register", endpoint))
  if err != nil {
    return
  }

  rBytes, err := ioutil.ReadAll(r.Body)
  if err != nil {
    return
  }

  if err = json.Unmarshal([]byte(rBytes), &resp); err != nil {
    return
  }
  
  return
}

func Acquire(endpoint string, crateId string) (resp server.AcquireResponse, err error) {
  r, err := http.Get(fmt.Sprintf("http://%s/acquire?CrateId=%s", endpoint, crateId))
  if err != nil {
    return
  }

  rBytes, err := ioutil.ReadAll(r.Body)
  if err != nil {
    return
  }

  if err = json.Unmarshal([]byte(rBytes), &resp); err != nil {
    return
  }

  return
}