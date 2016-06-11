package main

import (
    "log"
    "net/http"
    "encoding/json"

    "golang.org/x/net/context"
    "github.com/kamilmac/cubesdb/db"
    "github.com/kamilmac/cubesdb/middleware"
    "github.com/satori/go.uuid"
    "goji.io"
	"goji.io/pat"
)

type App struct {
    db           *db.DB
}

type Cube struct {
    ID          string
    Username    string
    Title       string
    Suffix      string
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
    req := ctx.Value("reqJSON").(middleware.Request)
    app.createCube(req["username"], req["title"], req["suffix"]) 
    res["status"] = "success"
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}

func (app *App) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    res := middleware.Response{}
    req := ctx.Value("reqJSON").(middleware.Request)
    app.delCube(req["username"], req["id"]) 
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
    app.db.Put(cube.Username, cube.ID, cubeJSON)
}

func (app *App) delCube(username, id string) {
    app.db.Delete(username, id)
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
    log.Println("Starting app")    
    app := App{}
    log.Println("initiating db")
    app.db = db.Init("/storage/cubes.db")
    log.Println("initiating db done")

    defer app.db.Close()
    
    setMux := goji.SubMux()
    setMux.UseC(middleware.Validate)
	setMux.HandleFuncC(pat.Post(""), app.set)
    
    getMux := goji.SubMux()
	getMux.HandleFuncC(pat.Post(""), app.getAll)
    
    delMux := goji.SubMux()
    delMux.UseC(middleware.Validate)
	delMux.HandleFuncC(pat.Post(""), app.delete)
    
    rootMux := goji.NewMux()
    rootMux.UseC(middleware.ParseJSON)
	rootMux.HandleC(pat.New("/api/v1/get"), getMux)
	rootMux.HandleC(pat.New("/api/v1/set"), setMux)
	rootMux.HandleC(pat.New("/api/v1/del"), delMux)
    
    log.Println("Running on port 5000")
	http.ListenAndServe(":5000", rootMux)
}
