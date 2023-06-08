package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
)

func mainy() {

  url := "http://localhost:8080/videos"
  method := "GET"

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return
  }
  req.Header.Add("Authorization", "Basic dXNlcm5hbWVFeGFtcGxlOnBhc3N3b3JkRXhhbXBsZQ==")

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body))
}