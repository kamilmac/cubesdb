package main

import (
    // "fmt"
    "net/http"
    "encoding/json"

    "github.com/kamilmac/cubesdb/db"
    "goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
)

type Cube struct {
    ID          string
    title       string
    suffix      string
}

type App struct {
    db      *db.DB
}

type Request map[string]string

type Response map[string]interface{}

func init() {
    
}

func main() {
    app := App{}
    app.db = db.Init("./cubes.db")
    defer app.db.Close()
    // app.db.Put("prints","someid",[]byte("someValue1"))
    // app.db.Put("prints","someOtherid",[]byte("someValue2"))
    // app.db.Put("prints","moreid",[]byte("someValue3"))
    // app.db.Put("prints","moreid",[]byte("someValuef"))
    // app.db.Delete("prindts", "moreidd")
    // fmt.Println(app.db.Get("prindts", "moreidd"))
    // fmt.Println(app.db.GetAll("prints"))
    mux := goji.NewMux()
	mux.HandleFuncC(pat.Post("/api/v1/getall"), app.getAll)
	// mux.HandleFuncC(pat.Get("/api/v1/get"), get)
	mux.HandleFuncC(pat.Post("/api/v1/set"), app.set)
	// mux.HandleFuncC(pat.Get("/api/v1/delete"), auth(delete))

	http.ListenAndServe("localhost:5010", mux)
}

func (app *App) getAll(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    res := Response{}
    req := Request{}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        res["status"] = "error"
        res["message"] = "Json req decoding error"
    } else {
        res["status"] = "success"
        res["data"] = app.db.GetAll(req["username"])
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}

func (app *App) set(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    res := Response{}
    req := Request{}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        res["status"] = "error"
        res["message"] = "Json req decoding error"
    } else {
        app.createCube(req["username"], req["title"], req["suffix"])       
        res["status"] = "success"
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}


func (app *App) createCube(username, title, suffix string) {
    // generate ID
            
    app.db.Put(username, "id9878760987", []byte("hello.png"))
}