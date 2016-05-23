package main

import (
    "fmt"
    
    "github.com/kamilmac/cubesdb/db"
)

type App struct {
    db *db.DB
}

func init() {
    
}

func main() {
    app := App{}
    app.db = db.Init("./cubes.db")
    fmt.Println(app)
}