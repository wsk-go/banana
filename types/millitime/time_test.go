package millitime

import (
	"context"
	"fmt"
	"github.com/JackWSK/banana/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

const MongoUrl = "mongodb://47.94.142.208:27017/admin?connect=direct&maxPoolSize=50&minPoolSize=10&slaveOk=true"

type Test struct {
	Hello *MilliTime `bson:"hello"`
}

func setupMongoClient(mongoUrl string) *mongo.Client {

	monitorOptions := options.Client().SetMonitor(&event.CommandMonitor{
		Started: func(i context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println("mongo command" + startedEvent.Command.String())
		},
	})

	//if conf.Profile != "dev" {
	//	monitorOptions.SetAuth(options.Credential{
	//		AuthSource: conf.MongoAuthSource,
	//		Username:   conf.MongoUserName,
	//		Password:   conf.MongoPassword,
	//	})
	//}

	//credentials := options.Client().SetAuth(options.Credential{
	//	AuthSource: "admin",
	//	Username:   "jcloudapp",
	//	Password:   "jcloudapp1231!",
	//})

	//client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl), monitorOptions, credentials)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl), monitorOptions)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	return client
}

func TestTime(t *testing.T) {
	c := setupMongoClient(MongoUrl)
	db := c.Database("test")
	col := db.Collection("hello")
	r, err := col.InsertOne(context.Background(), Test{
		Hello: utils.ToPtr(NewMilliTime()),
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(r.InsertedID)

	rr := col.FindOne(context.Background(), bson.M{})
	if rr.Err() != nil {
		t.Fatal(rr.Err())
	}
	var tt Test
	err = rr.Decode(&tt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tt)
}
