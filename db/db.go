package db

import (
    "fmt"
    "log"
    "time"
    
    "github.com/satori/go.uuid"
	"github.com/boltdb/bolt"
)

type DB struct {
    core *bolt.DB
}

func getUID() (id string) {
    return uuid.NewV4().String()
}

func Init(path string) *DB {
    database := DB{}
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	database.core = db 
    if err != nil {
		log.Fatal(err)
	}
	return &database
}

func (db *DB) Put(bucket string, key string, value []byte) {
    err := db.core.Update(func(tx *bolt.Tx) error {
        b, err := tx.CreateBucketIfNotExists([]byte(bucket))
        if err != nil {
            return fmt.Errorf("create bucket: %s", err)
        }
        b.Put([]byte(key), []byte(value))
        return nil
    })
    if err != nil {
		log.Fatal(err)
	}
}

func (db *DB) Get(bucket, key string) (v []byte) {
    db.core.View(func(tx *bolt.Tx) error {
	    b := tx.Bucket([]byte(bucket))
        v = b.Get([]byte(key))
	    return nil
	})
	return
}

func (db *DB) Delete(bucket, key string) {
    db.core.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(bucket))
        b.Delete([]byte(key))
        return nil
    })
}

func (db *DB) GetAll(bucket string) []string {
    list := []string{}
    db.core.View(func(tx *bolt.Tx) error {
	    b := tx.Bucket([]byte(bucket))
        if(b != nil) {
            b.ForEach(func(k, v []byte) error {
                list = append(list, string(v))
                return nil
            })
        }
        fmt.Println("list: ", list)
	    return nil
	})
	return list
}