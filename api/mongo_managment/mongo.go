package mongo_managment

import (
	p "api/util"
	"context"
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const props_name = "config/api.properties"

var (
	Ctx = context.TODO()
	TodoListCol *mongo.Collection
)

type TodoList struct {
	Message string `bson:"message"`
	Done bool `bson:"done"`
	Id int `bson:"_id"`
}

func Setup() {
	props := p.ReadProperties(props_name)
	host := props.GetString("mongo.host", "localhost")
	port := props.GetString("mongo.port", "27017")
	connectionURI := "mongodb://" + host + ":" + port
	cOptions := options.Client().ApplyURI(connectionURI)
	client, err := mongo.Connect(Ctx, cOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(Ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("todo-list")
	TodoListCol = db.Collection("list")
}

func getNextID() int {
	items, _ := GetItems()
	if items == nil {
		return 1
	}
	max := 1
	for _, item := range items {
		if item.Id > max {
			max = item.Id
		}
	}
	return max + 1
}

func CreateItem(t TodoList) (string, error) {
	t.Id = getNextID()
	result, err := TodoListCol.InsertOne(Ctx, t)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result.InsertedID), nil
}

func GetItem(id int) (TodoList, error) {
	var item TodoList
	filter := bson.D{{"_id", id}}
	err := TodoListCol.FindOne(Ctx, filter).Decode(&item)
	if err != nil {
		return item, err
	}
	return item, nil
}

func GetItems() ([]TodoList, error) {
	var items []TodoList
	var item TodoList

	cursor, err := TodoListCol.Find(Ctx, bson.D{})
	if err != nil {
		return items, err
	}
	defer cursor.Close(Ctx)

	for cursor.Next(Ctx) {
		err := cursor.Decode(&item)
		if err != nil {
			return items, err
		}
		items = append(items, item)
	}

	return items, nil
}

func UpdateItem(id int, status bool, message string) error {
	log.Println(id, status, message)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"done", status}}}}
	_, err := TodoListCol.UpdateOne(
		Ctx,
		filter,
		update,
	)
	update = bson.D{{"$set", bson.D{{"message", message}}}}
	_, err = TodoListCol.UpdateOne(
		Ctx,
		filter,
		update,
	)
	return err
}

func DeleteItem(id int) error {
	_, err := TodoListCol.DeleteOne(Ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}