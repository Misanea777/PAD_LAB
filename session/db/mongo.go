package db

import (
	"context"
	"errors"
	"log"

	"session/myTypes"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client

func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongodb:27017"))

	if err != nil {
		log.Printf("could not connect with mongodb:  %v", err)
		panic(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for {
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Default().Println("db not up, wait 10 sec")
			time.Sleep(5 * time.Second)
			continue
		} else {
			break
		}
	}
	log.Default().Println("successfuly connected")
	mongoClient = client
	mongoClient.Database("session").Collection("games")
	mongoClient.Database("session").Collection("maps")
	log.Default().Println("db Created")
}

type FullGameState struct {
	Id      uint64           `bson:"_id"`
	Chunks  []myTypes.Chunk  `bson:"chunks"`
	Players []myTypes.Player `bson:"players"`
}

func UpdateMap(id uint64, chnks []myTypes.Chunk) error {
	current, _ := GetMap(id)

	current = append(current, chnks...)

	log.Default().Println("updating....")
	log.Default().Println(len(current))

	collection := mongoClient.Database("session").Collection("maps")
	toStore := FullGameState{Id: id, Chunks: current}

	update := bson.M{
		"$set": bson.D{{"chunks", toStore.Chunks}},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, e1 := collection.UpdateOne(ctx, bson.D{{"_id", id}}, update, options.Update().SetUpsert(true))
	log.Default().Println(e1)
	return nil
}

func GetMap(gameId uint64) ([]myTypes.Chunk, error) {
	collection := mongoClient.Database("session").Collection("maps")
	var m FullGameState
	filter := bson.D{{"_id", gameId}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&m)

	return m.Chunks, err
}

func GetChunk(gameId uint64, id myTypes.PosAsID) (myTypes.Chunk, error) {
	var chnk myTypes.Chunk
	m, err := GetMap(gameId)
	if err != nil {
		return chnk, err
	}

	for _, c := range m {
		if c.Id == id {
			chnk = c
			return chnk, nil
		}
	}

	return chnk, errors.New("Not in db")
}

func SaveState(id uint64, players []myTypes.Player) error {
	collection := mongoClient.Database("session").Collection("maps")
	update := bson.M{
		"$set": bson.D{{"players", players}},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, e1 := collection.UpdateOne(ctx, bson.D{{"_id", id}}, update, options.Update().SetUpsert(true))
	log.Default().Println(e1)
	return nil
}

func GetState(id uint64) ([]myTypes.Player, error) {
	collection := mongoClient.Database("session").Collection("maps")
	var m FullGameState
	filter := bson.D{{"_id", id}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&m)

	return m.Players, err
}
