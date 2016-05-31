package main

import (
    "log"
    "net/http"
    "encoding/json"

    "github.com/kamilmac/cubesdb/utils"
    "github.com/kamilmac/cubesdb/db"
    "goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
    "github.com/satori/go.uuid"
)

var authAddress = "http://userauth/api/v1/auth"

type Cube struct {
    ID          string
    Username    string
    Title       string
    Suffix      string
}

type App struct {
    db      *db.DB
    utils   *Utils
}

type Request map[string]string

type Response map[string]interface{}

func init() {
    
}

func main() {
    app := App{}
    app.db = db.Init("./cubes.db")
    defer app.db.Close()
    mux := goji.NewMux()
	mux.HandleFuncC(pat.Post("/api/v1/getall"), app.getAll)
	mux.HandleFuncC(pat.Post("/api/v1/set"), app.set)
    mux.UseC(app.parseJSON)
    mux.UseC(app.validate)
	http.ListenAndServe("localhost:5010", mux)
}

func (app *App) validate(inner goji.Handler) goji.Handler {
	mw := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		log.Print("A: before")
        token := // parse token from request
        resp := app.utils.postJSON(authAddress, fmt.Sprintf(`{"token": "%v}`, token)
		inner.ServeHTTPC(ctx, w, r)
		log.Print("A: after")
	}
	return goji.HandlerFunc(mw)
}

func (app *App) parseJSON(inner goji.Handler) goji.Handler {
	mw := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		res := Response{}
        req := Request{}
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            res["status"] = "error"
            res["message"] = "Json req decoding error"
        } else {
            res["status"] = "success"
            res["data"] = app.getAllCubes(req["username"])
        }
		inner.ServeHTTPC(ctx, w, r)
		log.Print("A: after")
	}
	return goji.HandlerFunc(mw)
}

func (app *App) getAll(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    res := Response{}
    req := Request{}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        res["status"] = "error"
        res["message"] = "Json req decoding error"
    } else {
        res["status"] = "success"
        res["data"] = app.getAllCubes(req["username"])
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
    cube := Cube{
        ID: uuid.NewV4().String(),
        Username: username,
        Title: title,
        Suffix: suffix,
    }
    cubeJSON, err := json.Marshal(cube)
    if err != nil {
		log.Println("createCube json marshall error:", err)
	}
    log.Println(string(cubeJSON))
    app.db.Put(cube.Username, cube.ID, cubeJSON)
}

func (app *App) getAllCubes(username string) []Cube {
    all := app.db.GetAll(username)
    cubes := []Cube{}
    for _, v := range(all) {
        cube := Cube{}
        _ = json.Unmarshal(v, &cube)
        cubes = append(cubes, cube)
    }
    return cubes
}