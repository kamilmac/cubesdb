package utils

import (
    "net/http"
    "bytes"
    "fmt"
    "time"
    "encoding/json"
)

func PostJSON(url, payload string) map[string]interface{} {
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
    var bodyJSON map[string]interface{} 
    if err := json.NewDecoder(resp.Body).Decode(&bodyJSON); err != nil {
        fmt.Println("Couldn't parse bodyJSON response: ", err)
    }
    return bodyJSON
}