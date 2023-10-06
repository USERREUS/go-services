package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name string `bson:"name"`
}

func main() {
	// Устанавливаем параметры подключения к MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Проверяем подключение к MongoDB
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Подключение к MongoDB установлено")
	// Получаем коллекцию, в которую будем вставлять записи
	collection := client.Database("mydatabase").Collection("mycollection")
	// Создаем несколько документов (записей) для вставки
	documents := []interface{}{
		Person{Name: "Alice"},
		Person{Name: "Bob"}, Person{Name: "Charlie"},
	}
	// Вставляем документы в коллекцию
	insertResult, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Вставлено %v документов\n", len(insertResult.InsertedIDs))
	// Запрашиваем список всех записей в коллекции
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	var results []Person
	// Итерируем по результатам и добавляем их в слайс results
	for cursor.Next(context.TODO()) {
		var elem Person
		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, elem)
	}
	// Проверяем наличие ошибок при итерации
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	// Выводим список записей
	fmt.Println("Список записей:")
	for _, result := range results {
		fmt.Println(result.Name)
	}
	// Закрываем подключение к MongoDB
	defer client.Disconnect(context.TODO())
}
