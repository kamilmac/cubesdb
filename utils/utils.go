package utils

import (
    "net/http"
    "bytes"
    "fmt"
    "time"
    "io/ioutil"
)

func PostJSON(url, payload string) []byte {
    fmt.Println(url, payload)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
    if err != nil {
        fmt.Println("post JSON request parse err: ", err)
    }
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{
        Timeout: time.Duration(5 * time.Second),
    }
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("post JSON response err: ", err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}