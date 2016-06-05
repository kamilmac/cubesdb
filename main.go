package main

import (
    "log"
    "net/http"
    "fmt"
    "encoding/json"

    "golang.org/x/net/context"
    "github.com/kamilmac/cubesdb/db"
    "github.com/kamilmac/cubesdb/middleware"
    "github.com/satori/go.uuid"
    "goji.io"
	"goji.io/pat"
)

type Cube struct {
    ID          string
    Username    string
    Title       string
    Suffix      string
}

type App struct {
    db           *db.DB
}

func (app *App) getAll(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    res := middleware.Response{}
    req := ctx.Value("reqJSON").(middleware.Request)
    res["status"] = "success"
    res["data"] = app.getAllCubes(req["username"])
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}

func (app *App) set(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    res := middleware.Response{}
    req, _ := ctx.Value("reqJSON").(middleware.Request)
    fmt.Println("CONTEXT: ", req)
    app.createCube(req["username"], req["title"], req["suffix"]) 
    res["status"] = "success"
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

func main() {
    app := App{}
    app.db = db.Init("./cubes.db")
    defer app.db.Close()
    
    setMux := goji.SubMux()
    setMux.UseC(middleware.Validate)
	setMux.HandleFuncC(pat.Post(""), app.set)
    
    getMux := goji.SubMux()
	getMux.HandleFuncC(pat.Post(""), app.getAll)
    
    rootMux := goji.NewMux()
    rootMux.UseC(middleware.ParseJSON)
	rootMux.HandleC(pat.New("/api/v1/get"), getMux)
	rootMux.HandleC(pat.New("/api/v1/set"), setMux)
    
	http.ListenAndServe("localhost:5010", rootMux)
}
