package db

import (
	"context"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// DB TODO
type DB struct {
	Name string
	URI  string

	conn *mongo.Client
	db   *mongo.Database
}

// Connect creates a connection to the database
func (db *DB) Connect() {
	conn, err := mongo.Connect(context.Background(), db.URI, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	db.conn = conn
	db.db = db.conn.Database(db.Name)
}

// Insert insert an object into specified collection. the id of the inserted object and and wether it was added
func (db *DB) Insert(collection string, obj interface{}) (string, bool) {
	res, err := db.db.Collection(collection).InsertOne(context.Background(), obj)
	if err != nil {
		log.Println(err)
		return "", false
	}
	return res.InsertedID.(*bson.Element).Value().ObjectID().Hex(), true
}
