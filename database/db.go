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
	Timestamp   int64             `bson:"timestamp" json:"-"`
}

type WebhookInfo struct {
	ID              objectid.ObjectID `bson:"_id" json:"-"`
	WebhookURL      string            `bson:"webhookURL" json:"webhookURL"`
	MinTriggerValue int64             `bson:"minTriggerValue" json:"minTriggerValue"`
	Counter         int64             `bson:"counter" json:"-"`
	LatestTimestamp int64             `bson:"latestTimestamp" json:"-"` // the latest timestamp that invoked this webhook
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
	if track == (TrackInfo{}) {
		return track, false
	}

	return track, true
}

// GetTrackCount returns the number of tracks in the database
func (db *DB) GetTrackCount() (int64, error) {
	count, err := db.db.Collection("tracks").Count(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return count, err
}

// DeleteAllTracks returns the number of tracks in the database
func (db *DB) DeleteAllTracks() (int64, error) {
	col := db.db.Collection("tracks")
	count, err := col.Count(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
		return count, err
	}
	col.DeleteMany(context.Background(), bson.NewDocument())
	return count, err
}

// GetAllTracks returns all the tracks in the database
func (db *DB) GetAllTracks() ([]TrackInfo, error) {
	var cursor mongo.Cursor
	var err error
	cursor, err = db.db.Collection("tracks").Find(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var tracks []TrackInfo
	track := TrackInfo{}
	for cursor.Next(context.Background()) { // cannot for the life of me figure out the new mongo driver version
		err := cursor.Decode(&track) // looping through to find last...
		if err != nil {
			log.Fatal(err)
		}
		tracks = append(tracks, track)
	}
	return tracks, err
}

// GetAllInvokeWebhooks returns an rray of every webhook that should be invoked
func (db *DB) GetAllInvokeWebhooks() ([]WebhookInfo, error) {
	// subtracts 1 from each webhook's counter
	coll := db.db.Collection("webhooks")
	_, err := coll.UpdateMany(context.Background(), nil, bson.EC.SubDocumentFromElements("$inc",
		bson.EC.Int64("counter", -1)))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// selects all webhooks that should be triggered (counter = 0)
	cursor, err2 := coll.Find(context.Background(), bson.NewDocument(bson.EC.SubDocumentFromElements("triggerCount",
		bson.EC.Int64("$lte", 0))))
	if err2 != nil {
		fmt.Println(err2)
		return nil, err2
	}
	defer cursor.Close(context.Background())

	var whs []WebhookInfo
	wh := WebhookInfo{}
	// adds the webhooks to be invoked into the array that will be returned
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&wh)
		if err != nil {
			log.Fatal(err)
		}
		whs = append(whs, wh)
	}
	return whs, nil
}

// ResetWebhookCounter resets the counter and updates LatestTimestamp for the passed webhook
func (db *DB) ResetWebhookCounter(webhook WebhookInfo) {
	_, err := db.db.Collection("webhooks").UpdateMany(context.Background(),
		bson.NewDocument(bson.EC.ObjectID("_id", webhook.ID)),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				bson.EC.Int64("counter", webhook.MinTriggerValue),
				bson.EC.Int64("latestTimestamp", webhook.LatestTimestamp))))
	if err != nil {
		log.Fatal(err)
	}
}
