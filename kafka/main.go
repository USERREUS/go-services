package main

import (
	"log"

	"github.com/IBM/sarama"
)

func main() {
	go prod()
	go cons()
	select {}
}

func cons() {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer consumer.Close()
	partitionConsumer, err := consumer.ConsumePartition("order-topic", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Received message: %s\n", msg.Value)
			// Здесь вы можете сохранить сообщение в базу данных или выполнить другие действия
		}
	}
}

func prod() {
	config := sarama.NewConfig()
	//config.Producer.Return.Successes = true // Установите этот параметр в true
	producer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer producer.Close()
	// Генерируйте сообщение и отправляйте его в тему Kafka
	message := "orderid: 124, orderDate: 2023-10-06"
	producerMessage := &sarama.ProducerMessage{
		Topic: "order-topic",
		Value: sarama.StringEncoder(message),
	}
	producer.Input() <- producerMessage
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}
	log.Println("Message sent successfully")
}
