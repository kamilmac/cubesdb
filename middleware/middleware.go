package middleware

import (
    "net/http"
    "encoding/json"
    "fmt"

    "github.com/kamilmac/cubesdb/utils"
    "golang.org/x/net/context"
    "goji.io"
)

type Request    map[string]string
type Response   map[string]interface{}

var authAddress = "http://130.211.103.177/api/v1/auth"

func Validate(inner goji.Handler) goji.Handler {
	mw := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		res := Response{}
        req := ctx.Value("reqJSON").(Request)
        resp := utils.PostJSON(authAddress, fmt.Sprintf(`{"token": "%v"}`, req["token"]))
        if resp["status"] == "success" {
		    inner.ServeHTTPC(ctx, w, r)
        } else {
            res["status"] = "error"
            res["message"] = "Invalid token"
            json.NewEncoder(w).Encode(res)
            return
        }
	}
	return goji.HandlerFunc(mw)
}

func ParseJSON(inner goji.Handler) goji.Handler {
	mw := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		res := Response{}
        req := Request{}
        w.Header().Set("Access-Control-Allow-Origin", "*")     
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            res["status"] = "error"
            res["message"] = "Json req decoding error"
            json.NewEncoder(w).Encode(res)
            return
        }
        newCtx := context.WithValue(ctx, "reqJSON", req)
        inner.ServeHTTPC(newCtx, w, r)
	}
	return goji.HandlerFunc(mw)
}