package db

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// DB TODO
type DB struct {
	Name string
	URI  string

	conn *mongo.Client
	db   *mongo.Database
}

type TrackInfo struct {
	ID          objectid.ObjectID `bson:"_id" json:"-"`
	HDate       string            `bson:"H_date" json:"H_Date"`
	Pilot       string            `bson:"pilot" json:"pilot"`
	Glider      string            `bson:"glider" json:"glider"`
	GliderID    string            `bson:"glider_id" json:"glider_id"`
	TrackLength string            `bson:"track_length" json:"track_length"`
	TrackURL    string            `bson:"track_url" json:"track_url"`
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

// GetAllTrackIDs returns an array of all the track ids in the database
func (db *DB) GetAllTrackIDs() ([]objectid.ObjectID, error) {
	var cursor mongo.Cursor
	var err error
	cursor, err = db.db.Collection("tracks").Find(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer cursor.Close(context.Background())
	var ids []objectid.ObjectID
	track := TrackInfo{}
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&track)
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, track.ID)
	}
	return ids, err
}

// GetTrackByID returns the track given an id and true/false wether it was found
func (db *DB) GetTrackByID(id string) (TrackInfo, bool) {
	var cursor mongo.Cursor
	var err error
	track := TrackInfo{}
	objectID, _ := objectid.FromHex(id)
	cursor, err = db.db.Collection("tracks").Find(context.Background(), bson.NewDocument(bson.EC.ObjectID("_id", objectID)))
	if err != nil {
		fmt.Println(err)
		return track, false
	}
	//defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&track)
		if err != nil {
			log.Fatal(err)
		}
	}
	/*if track == (TrackInfo{}) {
		return track, false
	}*/

	return track, true
}
