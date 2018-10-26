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
	URL  string

	conn *mongo.Client
	db   *mongo.Database
}

// Connect creates a connection to the database
func (db *DB) Connect() {
	conn, err := mongo.Connect(context.Background(), db.URL, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	db.conn = conn
	db.db = db.conn.Database(db.Name)
}

// Insert insert an object into specified collection. the id of the inserted object and an error if any
func (db *DB) Insert(collection string, obj interface{}) (string, error) {
	res, err := db.db.Collection(collection).InsertOne(context.Background(), obj)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return res.InsertedID.(*bson.Element).Value().ObjectID().Hex(), err
}
