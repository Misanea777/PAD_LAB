package db

import (
	"context"
	"log"
	state "session/game/gameState"
	mapstate "session/game/mapState"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"crypto"
	_ "crypto/md5"
	"encoding/binary"
	"fmt"
	"reflect"
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

type ChunkModel struct {
	Id   uint32         `bson:"_id"`
	Chnk mapstate.Chunk `bson:"chunk"`
}

func newChunkModel(id uint32, chnk mapstate.Chunk) ChunkModel {
	return ChunkModel{Id: Hash(id, chnk.Id.PosX, chnk.Id.PosY), Chnk: chnk}
}

func Hash(objs ...interface{}) uint32 {
	digester := crypto.MD5.New()
	for _, ob := range objs {
		fmt.Fprint(digester, reflect.TypeOf(ob))
		fmt.Fprint(digester, ob)
	}
	data := binary.BigEndian.Uint32(digester.Sum(nil))
	return data
}

func UpdateChnk(id uint32, chnk mapstate.Chunk) error {
	collection := mongoClient.Database("session").Collection("maps")
	model := newChunkModel(id, chnk)
	update := bson.M{
		"$set": bson.M{"chunk": model.Chnk},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(ctx, bson.M{"_id": model.Id}, update, options.Update().SetUpsert(true))
	log.Default().Println("Atemting to save.....")
	return err
}

func GetChnk(id uint32, chnkId mapstate.PosAsID) (mapstate.Chunk, error) {
	collection := mongoClient.Database("session").Collection("maps")
	filter := bson.M{"_id": Hash(id, chnkId.PosX, chnkId.PosY)}

	var model ChunkModel

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&model)
	log.Default().Println("Atemting to retrevice.....")
	log.Default().Println(err)
	return model.Chnk, err
}

func UpdateSt(gms *state.GameState) error {
	collection := mongoClient.Database("session").Collection("games")

	update := bson.M{
		"$set": bson.M{"players": gms.Players},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.UpdateOne(ctx, bson.M{"_id": gms.Id}, update, options.Update().SetUpsert(true))
	return err
}

func GetSt(id uint32) (state.GameState, error) {
	var res state.GameState
	collection := mongoClient.Database("session").Collection("games")
	filter := bson.M{"_id": id}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(&res)
	return res, err
}

func GetAll() (int64, error) {
	collection := mongoClient.Database("session").Collection("games")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": `*`}
	itemCount, err := collection.CountDocuments(ctx, filter)

	return itemCount, err
}
