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
    defer app.db.Close()
    app.db.Put("prints","someid",[]byte("someValue1"))
    app.db.Put("prints","someOtherid",[]byte("someValue2"))
    app.db.Put("prints","moreid",[]byte("someValue3"))
    app.db.Put("prints","moreid",[]byte("someValuef"))
    app.db.Delete("prindts", "moreidd")
    fmt.Println(app.db.Get("prindts", "moreidd"))
    fmt.Println(app.db.GetAll("prints"))
}