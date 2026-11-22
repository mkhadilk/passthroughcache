package mongocacheadaptor

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChangeEvent struct {
	OperationType string `bson:"operationType"`
	Namespace     struct {
		DB   string `bson:"db"`
		Coll string `bson:"coll"`
	} `bson:"ns"`
	DocumentKey struct {
		ID interface{} `bson:"_id"` // interface{} handles ObjectID, string, int, UUID, etc.
	} `bson:"documentKey"`
}

func AddListner(db *mongo.Database, collectionname string, handlechange func(changedid string)) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pipeline := mongo.Pipeline{
		bson.D{{
			Key: "$match",
			Value: bson.D{{
				Key: "operationType",
				Value: bson.D{{
					Key:   "$in",
					Value: bson.A{"update", "delete", "replace"},
				}},
			}},
		}},
	}

	cs, err := db.Collection(collectionname).Watch(ctx, pipeline)
	if err != nil {
		log.Printf("Can not listen on the %s collection; error is %+v", collectionname, err)
		return err
	}

	defer cs.Close(ctx)

	for cs.Next(ctx) {
		var event ChangeEvent
		if err := cs.Decode(&event); err != nil {
			log.Printf("Failed to decode event: %v", err)
			continue
		}
		handlechange(event.DocumentKey.ID.(string))
	}

	return nil

}
