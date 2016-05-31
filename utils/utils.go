package utils

import (
    
)

type Utils interface{}

func (u *Utils) postJSON(url, payload string) []byte {
    req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
    if err != nil {
        return fmt.Sprintln("post JSON request parse err: ", err)
    }
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{
        Timeout: time.Duration(5 * time.Second),
    }
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Sprintln("post JSON response err: ", err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}